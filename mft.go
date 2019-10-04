/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package GoFor_MFT_Parser

import (
	"bytes"
	"errors"
	"fmt"
	ts "github.com/AlecRandazzo/Timestamp-Parser"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"sync"
)

type MasterFileTableRecord struct {
	RecordHeader                  RecordHeader
	StandardInformationAttributes StandardInformationAttribute
	FileNameAttributes            []FileNameAttribute
	DataAttribute                 DataAttribute
}

//TODO fill out these tags for json, csv, bson, and protobuf
type UseFulMftFields struct {
	RecordNumber     uint32       `json:"RecordNumber,number"`
	FilePath         string       `json:"FilePath,string"`
	FullPath         string       `json:"FullPath,string"`
	FileName         string       `json:"FileName,string"`
	SystemFlag       bool         `json:"SystemFlag,bool"`
	HiddenFlag       bool         `json:"HiddenFlag,bool"`
	ReadOnlyFlag     bool         `json:"ReadOnlyFlag,bool"`
	DirectoryFlag    bool         `json:"DirectoryFlag,bool"`
	DeletedFlag      bool         `json:"DeletedFlag,bool"`
	FnCreated        ts.TimeStamp `json:"FnCreated"`
	FnModified       ts.TimeStamp `json:"FnModified"`
	FnAccessed       ts.TimeStamp `json:"FnAccessed"`
	FnChanged        ts.TimeStamp `json:"FnChanged"`
	SiCreated        ts.TimeStamp `json:"SiCreated"`
	SiModified       ts.TimeStamp `json:"SiModified"`
	SiAccessed       ts.TimeStamp `json:"SiAccessed"`
	SiChanged        ts.TimeStamp `json:"SiChanged"`
	PhysicalFileSize uint64       `json:"PhysicalFileSize,number"`
}

type RawMasterFileTableRecord []byte

// Parse an already extracted MFT and write the results to a file.
func ParseMFT(reader io.Reader, writer OutputWriters, bytesPerCluster int64) (err error) {
	var buffer bytes.Buffer
	tee := io.TeeReader(reader, &buffer)
	directoryTree, _ := BuildDirectoryTree(tee)

	outputChannel := make(chan UseFulMftFields, 100)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go writer.Write(&outputChannel, &waitGroup)

	timeToBreak := false
	for timeToBreak == false {
		buffer := make([]byte, 1024)
		offset, err := reader.Read(buffer)
		if err == io.EOF {
			err = nil
			timeToBreak = true
		}
		rawMftRecord := RawMasterFileTableRecord(buffer)
		mftRecord, err := rawMftRecord.Parse(bytesPerCluster)
		if err != nil {
			log.WithFields(log.Fields{
				"mft_offset":   offset,
				"deleted_flag": mftRecord.RecordHeader.Flags.FlagDeleted,
			}).Debug(err)
			continue
		}
		if len(mftRecord.FileNameAttributes) == 0 {
			continue
		}

		usefulMftFields := GetUsefulMftFields(mftRecord, directoryTree)
		outputChannel <- usefulMftFields

	}
	close(outputChannel)
	waitGroup.Wait()
	return
}

func GetUsefulMftFields(mftRecord MasterFileTableRecord, directoryTree DirectoryTree) (useFulMftFields UseFulMftFields) {
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

// Parse the bytes of an MFT record
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
		err = fmt.Errorf("failed to parse the raw mft record: %v", err)
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
