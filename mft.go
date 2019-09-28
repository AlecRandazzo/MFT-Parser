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
	"fmt"
	ts "github.com/AlecRandazzo/Timestamp-Parser"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"sync"
)

type MasterFileTableRecord struct {
	BytesPerCluster               int64
	RecordHeader                  RecordHeader
	StandardInformationAttributes StandardInformationAttributes
	FileNameAttributes            []FileNameAttribute
	DataAttributes                DataAttribute
	Attributes                    Attributes
	UseFullMftFields              useFullMftFields
}

type useFullMftFields struct {
	RecordNumber     uint32       `json:"RecordNumber"`
	FilePath         string       `json:"FilePath"`
	FullPath         string       `json:"FullPath"`
	FileName         string       `json:"FileName"`
	SystemFlag       bool         `json:"SystemFlag"`
	HiddenFlag       bool         `json:"HiddenFlag"`
	ReadOnlyFlag     bool         `json:"ReadOnlyFlag"`
	DirectoryFlag    bool         `json:"DirectoryFlag"`
	DeletedFlag      bool         `json:"DeletedFlag"`
	FnCreated        ts.TimeStamp `json:"FnCreated"`
	FnModified       ts.TimeStamp `json:"FnModified"`
	FnAccessed       ts.TimeStamp `json:"FnAccessed"`
	FnChanged        ts.TimeStamp `json:"FnChanged"`
	SiCreated        ts.TimeStamp `json:"SiCreated"`
	SiModified       ts.TimeStamp `json:"SiModified"`
	SiAccessed       ts.TimeStamp `json:"SiAccessed"`
	SiChanged        ts.TimeStamp `json:"SiChanged"`
	PhysicalFileSize uint64       `json:"PhysicalFileSize"`
}

type RawMasterFileTableRecord []byte

// Parse an already extracted MFT and write the results to a file.
func ParseMFT(reader io.Reader, writer OutputWriters, numberOfWorkers int) (err error) {

	var directoryTree DirectoryTree
	err = directoryTree.Build(reader, numberOfWorkers)
	if err != nil {
		return
	}

	outputChannel := make(chan useFullMftFields, 100)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go writer.Write(&outputChannel, &waitGroup)

	for {
		buffer := make([]byte, 1024)
		offset, err := reader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}
		rawMftRecord := RawMasterFileTableRecord(buffer)
		mftRecord, err := rawMftRecord.Parse()
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

		mftRecord.UseFullMftFields.get(&mftRecord, &directoryTree)
		outputChannel <- mftRecord.UseFullMftFields

	}
	close(outputChannel)
	waitGroup.Wait()
	return
}

func (useFullMftFields *useFullMftFields) get(mftRecord *MasterFileTableRecord, directoryTree *DirectoryTree) {
	for _, record := range mftRecord.FileNameAttributes {
		if strings.Contains(record.FileNamespace, "WIN32") || strings.Contains(record.FileNamespace, "POSIX") {
			if directory, ok := (*directoryTree)[record.ParentDirRecordNumber]; ok {
				mftRecord.UseFullMftFields.FileName = record.FileName
				mftRecord.UseFullMftFields.FilePath = directory
				mftRecord.UseFullMftFields.FullPath = mftRecord.UseFullMftFields.FilePath + mftRecord.UseFullMftFields.FileName
			} else {
				mftRecord.UseFullMftFields.FileName = record.FileName
				mftRecord.UseFullMftFields.FilePath = "$ORPHANFILE"
				mftRecord.UseFullMftFields.FullPath = mftRecord.UseFullMftFields.FilePath + mftRecord.UseFullMftFields.FileName
			}
			mftRecord.UseFullMftFields.RecordNumber = mftRecord.RecordHeader.RecordNumber
			mftRecord.UseFullMftFields.SystemFlag = record.FileNameFlags.System
			mftRecord.UseFullMftFields.HiddenFlag = record.FileNameFlags.Hidden
			mftRecord.UseFullMftFields.ReadOnlyFlag = record.FileNameFlags.ReadOnly
			mftRecord.UseFullMftFields.DirectoryFlag = mftRecord.RecordHeader.Flags.FlagDirectory
			mftRecord.UseFullMftFields.DeletedFlag = mftRecord.RecordHeader.Flags.FlagDeleted
			mftRecord.UseFullMftFields.FnCreated = record.FnCreated
			mftRecord.UseFullMftFields.FnModified = record.FnModified
			mftRecord.UseFullMftFields.FnAccessed = record.FnAccessed
			mftRecord.UseFullMftFields.FnChanged = record.FnChanged
			mftRecord.UseFullMftFields.SiCreated = mftRecord.StandardInformationAttributes.SiCreated
			mftRecord.UseFullMftFields.SiModified = mftRecord.StandardInformationAttributes.SiModified
			mftRecord.UseFullMftFields.SiAccessed = mftRecord.StandardInformationAttributes.SiAccessed
			mftRecord.UseFullMftFields.SiChanged = mftRecord.StandardInformationAttributes.SiChanged
			mftRecord.UseFullMftFields.PhysicalFileSize = record.PhysicalFileSize
			break
		}
	}
}

// Parse the bytes of an MFT record
func (rawMftRecord *RawMasterFileTableRecord) Parse() (mftRecord MasterFileTableRecord, err error) {
	err = mftRecord.RecordHeader.Parse(*rawMftRecord)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	rawMftRecord.trimSlackSpace()
	mftRecord.Attributes.Parse(*rawMftRecord, mftRecord.RecordHeader.AttributesOffset)
	if err != nil {
		err = fmt.Errorf("failed to get Attribute list: %w", err)
		return
	}

	const codeFileName = 0x30
	const codeStandardInformation = 0x10
	const codeData = 0x80

	for _, attribute := range mftRecord.Attributes {
		switch attribute.AttributeType {
		case codeFileName:
			fileNameAttribute := FileNameAttribute{}
			err = fileNameAttribute.Parse(attribute)
			if err != nil {
				continue
			}
			mftRecord.FileNameAttributes = append(mftRecord.FileNameAttributes, fileNameAttribute)
		case codeStandardInformation:
			err = mftRecord.StandardInformationAttributes.Parse(attribute)
			if err != nil {
				err = fmt.Errorf("failed to get standard info Attribute %w", err)
				return
			}
		case codeData:
			err = mftRecord.DataAttributes.Parse(attribute, mftRecord.BytesPerCluster)
			if err != nil {
				err = fmt.Errorf("failed to get data Attribute %w", err)
				return
			}
		}
	}
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
