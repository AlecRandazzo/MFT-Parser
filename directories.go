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
	"io"
	"strings"
	"sync"
)

type Directory struct {
	DirectoryName      string
	ParentRecordNumber uint64
}

type DirectoryList map[uint64]Directory

type MappedDirectories map[uint64]string

// Quickly checks the bytes of an MFT record to determine if it is a Directory or not.
func (mftRecord *MasterFileTableRecord) isThisADirectory() {

	const offsetRecordFlag = 0x16
	const codeDirectory = 0x03
	if len(mftRecord.MftRecordBytes) <= offsetRecordFlag {
		mftRecord.RecordHeader.Flags.FlagDirectory = false
		return
	}
	recordFlag := mftRecord.MftRecordBytes[offsetRecordFlag]
	if recordFlag == codeDirectory {
		mftRecord.RecordHeader.Flags.FlagDirectory = true
	} else {
		mftRecord.RecordHeader.Flags.FlagDirectory = false
	}
	return
}

// Creates a list of directories from a channel of MFR record bytes.
func (directoryList *DirectoryList) Create(inboundBuffer *chan []byte, directoryListChannel *chan map[uint64]Directory, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	var openChannel = true
	var err error
	*directoryList = make(map[uint64]Directory)
	for openChannel == true {
		var mftRecord MasterFileTableRecord
		mftRecord.MftRecordBytes, openChannel = <-*inboundBuffer
		mftRecord.isThisADirectory()
		if mftRecord.RecordHeader.Flags.FlagDirectory == false {
			continue
		}
		err = mftRecord.RecordHeader.Parse(mftRecord.MftRecordBytes)
		if err != nil {
			continue
		}

		err = mftRecord.Attributes.Parse(mftRecord.MftRecordBytes, mftRecord.RecordHeader.AttributesOffset)
		if err != nil {
			continue
		}

		const codeFileName = 0x30
		for _, attribute := range mftRecord.Attributes {
			switch attribute.AttributeType {
			case codeFileName:
				fileNameAttribute := FileNameAttribute{}
				err = fileNameAttribute.Parse(attribute)
				if err != nil {
					continue
				}
				mftRecord.FileNameAttributes = append(mftRecord.FileNameAttributes, fileNameAttribute)
			default:
				continue
			}
		}

		for _, attribute := range mftRecord.FileNameAttributes {
			if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") {

				(*directoryList)[uint64(mftRecord.RecordHeader.RecordNumber)] = Directory{
					DirectoryName:      attribute.FileName,
					ParentRecordNumber: attribute.ParentDirRecordNumber,
				}
				break
			}
		}
	}
	*directoryListChannel <- *directoryList
	return
}

// Combines a running list of directories from a channel in order to create the systems Directory trees.
func (file *MftFile) CombineDirectoryInformation(directoryListChannel *chan map[uint64]Directory, waitForDirectoryCombination *sync.WaitGroup) {
	defer waitForDirectoryCombination.Done()

	file.MappedDirectories = make(map[uint64]string)

	// Merge lists
	var masterDirectoryList map[uint64]Directory
	masterDirectoryList = make(map[uint64]Directory)
	openChannel := true

	for openChannel == true {
		var directoryList map[uint64]Directory
		directoryList = make(map[uint64]Directory)
		directoryList, openChannel = <-*directoryListChannel
		for key, value := range directoryList {
			masterDirectoryList[key] = value
		}
	}

	for recordNumber, directoryMetadata := range masterDirectoryList {
		mappingDirectory := directoryMetadata.DirectoryName
		parentRecordNumberPointer := directoryMetadata.ParentRecordNumber
		for {
			if _, ok := masterDirectoryList[parentRecordNumberPointer]; ok {
				if recordNumber == 5 {
					mappingDirectory = ":\\"
					file.MappedDirectories[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = ":\\" + mappingDirectory
					file.MappedDirectories[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = masterDirectoryList[parentRecordNumberPointer].DirectoryName + "\\" + mappingDirectory
				parentRecordNumberPointer = masterDirectoryList[parentRecordNumberPointer].ParentRecordNumber
				continue
			}
			file.MappedDirectories[recordNumber] = "$ORPHANFILE\\" + mappingDirectory
			break
		}
	}
	return
}

// Builds a list of directories for the purpose of of mapping MFT records to their parent directories.
func (file *MftFile) BuildDirectoryTree() (err error) {
	var waitGroup sync.WaitGroup
	numberOfWorkers := 4
	bufferChannel := make(chan []byte, 100)
	directoryListChannel := make(chan map[uint64]Directory, numberOfWorkers)
	for i := 0; i <= numberOfWorkers; i++ {
		waitGroup.Add(1)
		directoryList := DirectoryList{}
		go directoryList.Create(&bufferChannel, &directoryListChannel, &waitGroup)
	}

	var waitForDirectoryCombination sync.WaitGroup
	waitForDirectoryCombination.Add(1)
	go file.CombineDirectoryInformation(&directoryListChannel, &waitForDirectoryCombination)
	var offset int64 = 0
	for {
		buffer := make([]byte, 1024)
		_, err = file.FileHandle.ReadAt(buffer, offset)
		if err == io.EOF {
			err = nil
			break
		}
		bufferChannel <- buffer
		offset += 1024
	}

	close(bufferChannel)
	waitGroup.Wait()
	close(directoryListChannel)
	waitForDirectoryCombination.Wait()
	return
}
