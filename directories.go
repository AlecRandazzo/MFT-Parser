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
	"errors"
	"io"
	"strings"
	"sync"
)

type directory struct {
	DirectoryName      string
	ParentRecordNumber uint64
}

type unresolvedDirectoryList map[uint64]directory

type DirectoryTree map[uint64]string

// Quickly checks the bytes of an MFT record to determine if it is a directory or not.
func (rawMftRecord *RawMasterFileTableRecord) IsThisADirectory() (result bool, err error) {

	const offsetRecordFlag = 0x16
	const codeDirectory = 0x03
	if len(*rawMftRecord) == 0 {
		result = false
		err = errors.New("RawMasterFileTableRecord.Parse() received nil bytes ")
		return
	}
	if len(*rawMftRecord) <= offsetRecordFlag {
		result = false
		err = errors.New("RawMasterFileTableRecord.Parse() received not enough bytes ")
		return
	}
	recordFlag := []byte(*rawMftRecord)[offsetRecordFlag]
	if recordFlag == codeDirectory {
		result = true
	} else {
		result = false
	}
	return
}

// Creates a list of directories from a channel of MFR record bytes.
func (unresolvedDirectoryList *unresolvedDirectoryList) create(inboundBuffer *chan RawMasterFileTableRecord, directoryListChannel *chan unresolvedDirectoryList, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	openChannel := true
	const codeFileName = 0x30
	//*unresolvedDirectoryList = make(map[uint64]directory)
	for openChannel == true {
		var rawMftRecord RawMasterFileTableRecord
		rawMftRecord, openChannel = <-*inboundBuffer
		result, err := rawMftRecord.IsThisADirectory()
		if result == false {
			continue
		}
		var mftRecord MasterFileTableRecord
		err = mftRecord.RecordHeader.Parse(rawMftRecord)
		if err != nil {
			continue
		}

		mftRecord.Attributes.Parse(rawMftRecord, mftRecord.RecordHeader.AttributesOffset)

		for _, attribute := range mftRecord.Attributes {
			switch attribute.AttributeType {
			case codeFileName:
				fileNameAttribute := FileNameAttribute{}
				err = fileNameAttribute.Parse(attribute)
				if err != nil {
					continue
				}
				if strings.Contains(fileNameAttribute.FileNamespace, "WIN32") == true || strings.Contains(fileNameAttribute.FileNamespace, "POSIX") {
					(*unresolvedDirectoryList)[uint64(mftRecord.RecordHeader.RecordNumber)] = directory{
						DirectoryName:      fileNameAttribute.FileName,
						ParentRecordNumber: fileNameAttribute.ParentDirRecordNumber,
					}
					break
				}
			default:
				continue
			}
		}

	}
	*directoryListChannel <- *unresolvedDirectoryList
	return
}

// Combines a running list of directories from a channel in order to create the systems directory trees.
func (directoryTree *DirectoryTree) resolve(unresolvedDirectoryListChannel *chan unresolvedDirectoryList, waitForDirectoryResolution *sync.WaitGroup) {
	defer waitForDirectoryResolution.Done()

	// Merge lists
	var masterUnresolvedDirectoryList unresolvedDirectoryList
	masterUnresolvedDirectoryList = make(unresolvedDirectoryList)
	openChannel := true

	for openChannel == true {
		var unresolvedDirectoryList unresolvedDirectoryList
		unresolvedDirectoryList = make(map[uint64]directory)
		unresolvedDirectoryList, openChannel = <-*unresolvedDirectoryListChannel
		for key, value := range unresolvedDirectoryList {
			masterUnresolvedDirectoryList[key] = value
		}
	}

	directoryTree.findParentDirectories(&masterUnresolvedDirectoryList)
	return
}

func (directoryTree *DirectoryTree) findParentDirectories(masterUnresolvedDirectoryList *unresolvedDirectoryList) {
	for recordNumber, directoryMetadata := range *masterUnresolvedDirectoryList {
		mappingDirectory := directoryMetadata.DirectoryName
		parentRecordNumberPointer := directoryMetadata.ParentRecordNumber
		for {
			if _, ok := (*masterUnresolvedDirectoryList)[parentRecordNumberPointer]; ok {
				if recordNumber == 5 {
					mappingDirectory = "\\"
					(*directoryTree)[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = "\\" + mappingDirectory
					(*directoryTree)[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = (*masterUnresolvedDirectoryList)[parentRecordNumberPointer].DirectoryName + "\\" + mappingDirectory
				parentRecordNumberPointer = (*masterUnresolvedDirectoryList)[parentRecordNumberPointer].ParentRecordNumber
				continue
			}
			(*directoryTree)[recordNumber] = "$ORPHANFILE\\" + mappingDirectory
			break
		}
	}
}

// Builds a list of directories for the purpose of of mapping MFT records to their parent directories.
func (directoryTree *DirectoryTree) Build(reader io.Reader, numberOfWorkers int) (err error) {
	var waitGroup sync.WaitGroup
	*directoryTree = make(DirectoryTree)
	bufferChannel := make(chan RawMasterFileTableRecord, 100)
	unresolvedDirectoryListChannel := make(chan unresolvedDirectoryList, numberOfWorkers)
	for i := 0; i <= numberOfWorkers; i++ {
		waitGroup.Add(1)
		unresolvedDirectoryList := unresolvedDirectoryList{}
		go unresolvedDirectoryList.create(&bufferChannel, &unresolvedDirectoryListChannel, &waitGroup)
	}

	var waitForDirectoryResolution sync.WaitGroup
	waitForDirectoryResolution.Add(1)

	go directoryTree.resolve(&unresolvedDirectoryListChannel, &waitForDirectoryResolution)
	for {
		buffer := make(RawMasterFileTableRecord, 1024)
		_, err = reader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}
		bufferChannel <- buffer
	}

	close(bufferChannel)
	waitGroup.Wait()
	close(unresolvedDirectoryListChannel)
	waitForDirectoryResolution.Wait()
	return
}
