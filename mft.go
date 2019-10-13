/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package mft

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// MasterFileTableRecord contains information on a parsed MFT record
type MasterFileTableRecord struct {
	RecordHeader                  RecordHeader
	StandardInformationAttributes StandardInformationAttribute
	FileNameAttributes            []FileNameAttribute
	DataAttribute                 DataAttribute
}

//TODO fill out these tags for json, csv, bson, and protobuf

// UsefulMftFields contains a downselected list of fields that are actually valuable to an analyst. This is the data that is written after parsing an MFT.
type UsefulMftFields struct {
	RecordNumber     uint32    `json:"RecordNumber,number"`
	FilePath         string    `json:"FilePath,string"`
	FullPath         string    `json:"FullPath,string"`
	FileName         string    `json:"FileName,string"`
	SystemFlag       bool      `json:"SystemFlag,bool"`
	HiddenFlag       bool      `json:"HiddenFlag,bool"`
	ReadOnlyFlag     bool      `json:"ReadOnlyFlag,bool"`
	DirectoryFlag    bool      `json:"DirectoryFlag,bool"`
	DeletedFlag      bool      `json:"DeletedFlag,bool"`
	FnCreated        time.Time `json:"FnCreated"`
	FnModified       time.Time `json:"FnModified"`
	FnAccessed       time.Time `json:"FnAccessed"`
	FnChanged        time.Time `json:"FnChanged"`
	SiCreated        time.Time `json:"SiCreated"`
	SiModified       time.Time `json:"SiModified"`
	SiAccessed       time.Time `json:"SiAccessed"`
	SiChanged        time.Time `json:"SiChanged"`
	PhysicalFileSize uint64    `json:"PhysicalFileSize,number"`
}

// RawMasterFileTableRecord is a []byte alias for raw mft record. Used with the Parse() method.
type RawMasterFileTableRecord []byte

// ParseMFT takes an input file os.File and writes the results to the io.Writer. The format of the data sent to the io.Writer is dependent on what ResultWriter is used. The bytes per cluster input is typically 4096
func ParseMFT(inputFile *os.File, writer ResultWriter, streamer io.Writer, bytesPerCluster int64) {
	directoryTree, _ := BuildDirectoryTree(inputFile)
	outputChannel := make(chan UsefulMftFields, 100)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go writer.ResultWriter(streamer, &outputChannel, &waitGroup)
	// Seek back to the beginning of the file
	_, _ = inputFile.Seek(0, 0)
	ParseMftRecords(inputFile, bytesPerCluster, directoryTree, &outputChannel)
	waitGroup.Wait()
	return
}

// ParseMftRecords parses a stream of mft record bytes and sends the results to an output channel. Records from this output channel are popped off by the ResultWriter used in the ParseMFT() method.
func ParseMftRecords(reader io.Reader, bytesPerCluster int64, directoryTree DirectoryTree, outputChannel *chan UsefulMftFields) {
	for {
		buffer := make([]byte, 1024)
		_, err := reader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}
		rawMftRecord := RawMasterFileTableRecord(buffer)
		mftRecord, err := rawMftRecord.Parse(bytesPerCluster)
		if err != nil {
			continue
		}

		usefulMftFields := GetUsefulMftFields(mftRecord, directoryTree)
		*outputChannel <- usefulMftFields

	}
	close(*outputChannel)
	return
}

// GetUsefulMftFields will pull out and return just the MFT record fields that are useful to an analyst.
func GetUsefulMftFields(mftRecord MasterFileTableRecord, directoryTree DirectoryTree) (useFulMftFields UsefulMftFields) {
	for _, record := range mftRecord.FileNameAttributes {
		if strings.Contains(record.FileNamespace, "WIN32") || strings.Contains(record.FileNamespace, "POSIX") {
			if directory, ok := directoryTree[record.ParentDirRecordNumber]; ok {
				useFulMftFields.FileName = record.FileName
				useFulMftFields.FilePath = directory
				useFulMftFields.FullPath = useFulMftFields.FilePath + useFulMftFields.FileName
			} else {
				useFulMftFields.FileName = record.FileName
				useFulMftFields.FilePath = "$ORPHANFILE\\"
				useFulMftFields.FullPath = useFulMftFields.FilePath + useFulMftFields.FileName
			}
			useFulMftFields.RecordNumber = mftRecord.RecordHeader.RecordNumber
			useFulMftFields.SystemFlag = record.FileNameFlags.System
			useFulMftFields.HiddenFlag = record.FileNameFlags.Hidden
			useFulMftFields.ReadOnlyFlag = record.FileNameFlags.ReadOnly
			useFulMftFields.DirectoryFlag = mftRecord.RecordHeader.Flags.FlagDirectory
			useFulMftFields.DeletedFlag = mftRecord.RecordHeader.Flags.FlagDeleted
			useFulMftFields.FnCreated = record.FnCreated
			useFulMftFields.FnModified = record.FnModified
			useFulMftFields.FnAccessed = record.FnAccessed
			useFulMftFields.FnChanged = record.FnChanged
			useFulMftFields.SiCreated = mftRecord.StandardInformationAttributes.SiCreated
			useFulMftFields.SiModified = mftRecord.StandardInformationAttributes.SiModified
			useFulMftFields.SiAccessed = mftRecord.StandardInformationAttributes.SiAccessed
			useFulMftFields.SiChanged = mftRecord.StandardInformationAttributes.SiChanged
			useFulMftFields.PhysicalFileSize = record.PhysicalFileSize
			break
		}
	}

	return
}

// Parse parses the raw MFT record receiver and returns a parsed mft record.
func (rawMftRecord RawMasterFileTableRecord) Parse(bytesPerCluster int64) (mftRecord MasterFileTableRecord, err error) {
	// Sanity checks
	sizeOfRawMftRecord := len(rawMftRecord)
	if sizeOfRawMftRecord == 0 {
		err = errors.New("received nil bytes")
		return
	}
	if bytesPerCluster == 0 {
		err = errors.New("bytes per cluster of 0, typically this value is 4096")
		return
	}
	result, err := rawMftRecord.IsThisAnMftRecord()
	if err != nil {
		err = fmt.Errorf("failed to parse the raw mft record: %v", err)
		return
	}
	if result == false {
		err = fmt.Errorf("this is not an mft record: %v", err)
		return
	}

	rawMftRecord.trimSlackSpace()

	rawRecordHeader, err := rawMftRecord.GetRawRecordHeader()
	if err != nil {
		err = fmt.Errorf("failed to parse MFT record header: %v", err)
		return
	}

	mftRecord.RecordHeader, _ = rawRecordHeader.Parse()

	var rawAttributes RawAttributes
	rawAttributes, err = rawMftRecord.GetRawAttributes(mftRecord.RecordHeader)
	if err != nil {
		err = fmt.Errorf("failed to get raw data attributes: %v", err)
		return
	}

	mftRecord.FileNameAttributes, mftRecord.StandardInformationAttributes, mftRecord.DataAttribute, _ = rawAttributes.Parse(bytesPerCluster)
	return
}

// Trims off slack space after end sequence 0xffffffff
func (rawMftRecord *RawMasterFileTableRecord) trimSlackSpace() {
	lenMftRecordBytes := len(*rawMftRecord)
	mftRecordEndByteSequence := []byte{0xff, 0xff, 0xff, 0xff}
	for i := 0; i < (lenMftRecordBytes - 4); i++ {
		if bytes.Equal([]byte(*rawMftRecord)[i:i+0x04], mftRecordEndByteSequence) {
			*rawMftRecord = []byte(*rawMftRecord)[:i]
			break
		}
	}
}
