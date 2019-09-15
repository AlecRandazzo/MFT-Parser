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
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type MasterFileTableRecord struct {
	BytesPerCluster               int64
	RecordHeader                  RecordHeader
	StandardInformationAttributes StandardInformationAttributes
	FileNameAttributes            []FileNameAttribute
	DataAttributes                DataAttribute
	MftRecordBytes                []byte
	Attributes                    Attributes
}

type MftFile struct {
	FileHandle        *os.File
	MappedDirectories map[uint64]string
	OutputChannel     chan MasterFileTableRecord
}

// Parse an already extracted MFT and write the results to a file.
func ParseMFT(mftFilePath, outFileName string) (err error) {
	file := MftFile{}
	file.FileHandle, err = os.Open(mftFilePath)
	if err != nil {
		err = fmt.Errorf("failed to open MFT file %s: %w", mftFilePath, err)
		return
	}
	defer file.FileHandle.Close()

	err = file.BuildDirectoryTree()
	if err != nil {
		return
	}

	file.OutputChannel = make(chan MasterFileTableRecord, 100)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go file.MftToCSV(outFileName, &waitGroup)

	var offset int64 = 0
	for {
		buffer := make([]byte, 1024)
		_, err = file.FileHandle.ReadAt(buffer, offset)
		if err == io.EOF {
			err = nil
			break
		}
		mftRecord := MasterFileTableRecord{}
		mftRecord.MftRecordBytes = buffer
		err = mftRecord.Parse()
		if err != nil {
			log.WithFields(log.Fields{
				"mft_offset":   offset,
				"deleted_flag": mftRecord.RecordHeader.Flags.FlagDeleted,
			}).Debug(err)
			offset += 1024
			continue
		}
		file.OutputChannel <- mftRecord
		offset += 1024
		if len(mftRecord.FileNameAttributes) == 0 {
			continue
		}

	}
	close(file.OutputChannel)
	waitGroup.Wait()
	return
}

// Parse the bytes of an MFT record
func (mftRecord *MasterFileTableRecord) Parse() (err error) {
	err = mftRecord.RecordHeader.Parse(mftRecord.MftRecordBytes)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}
	mftRecord.TrimSlackSpace()
	err = mftRecord.Attributes.Parse(mftRecord.MftRecordBytes, mftRecord.RecordHeader.AttributesOffset)
	if err != nil {
		err = fmt.Errorf("failed to get attribute list: %w", err)
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
				err = fmt.Errorf("failed to get standard info attribute %w", err)
				return
			}
		case codeData:
			err = mftRecord.DataAttributes.Parse(attribute, mftRecord.BytesPerCluster)
			if err != nil {
				err = fmt.Errorf("failed to get data attribute %w", err)
				return
			}
		}
	}
	return
}

// Trims off slack space after end sequence 0xffffffff
func (mftRecord *MasterFileTableRecord) TrimSlackSpace() {
	lenMftRecordBytes := len(mftRecord.MftRecordBytes)
	mftRecordEndByteSequence := []byte{0xff, 0xff, 0xff, 0xff}
	for i := 0; i < (lenMftRecordBytes - 4); i++ {
		if bytes.Equal(mftRecord.MftRecordBytes[i:i+0x04], mftRecordEndByteSequence) {
			mftRecord.MftRecordBytes = mftRecord.MftRecordBytes[:i]
			break
		}
	}
}
