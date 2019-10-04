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
	"fmt"
	"io"
	"strings"
)

type directory struct {
	RecordNumber       uint64
	DirectoryName      string
	ParentRecordNumber uint64
}

type UnresolvedDirectoryTree map[uint64]directory

type DirectoryTree map[uint64]string

// Quickly checks the bytes of an MFT record to determine if it is a directory or not.
func (rawMftRecord RawMasterFileTableRecord) IsThisADirectory() (result bool, err error) {

	const offsetRecordFlag = 0x16
	const codeDirectory = 0x03
	sizeOfRawMFTRecord := len(rawMftRecord)
	if sizeOfRawMFTRecord == 0 {
		result = false
		err = errors.New("RawMasterFileTableRecord.Parse() received nil bytes ")
		return
	}
	if sizeOfRawMFTRecord <= offsetRecordFlag {
		result = false
		err = errors.New("RawMasterFileTableRecord.Parse() received not enough bytes ")
		return
	}
	recordFlag := rawMftRecord[offsetRecordFlag]
	if recordFlag == codeDirectory {
		result = true
	} else {
		result = false
	}
	return
}

func convertRawMFTRecordToDirectory(rawMftRecord RawMasterFileTableRecord) (directory directory, err error) {
	result, err := rawMftRecord.IsThisADirectory()
	if result == false {
		err = errors.New("this is not a directory")
		return
	}
	rawRecordHeader, err := rawMftRecord.GetRawRecordHeader()
	if err != nil {
		err = fmt.Errorf("failed to parse get record header: %v", err)
		return
	}
	recordHeader, _ := rawRecordHeader.Parse()

	rawAttributes, err := rawMftRecord.GetRawAttributes(recordHeader)
	if err != nil {
		err = fmt.Errorf("failed to get raw attributes: %v", err)
		return
	}
	doesntMatter := int64(4096)
	fileNameAttributes, _, _, err := rawAttributes.Parse(doesntMatter)
	for _, fileNameAttribute := range fileNameAttributes {
		if strings.Contains(fileNameAttribute.FileNamespace, "WIN32") == true || strings.Contains(fileNameAttribute.FileNamespace, "POSIX") {
			directory.RecordNumber = uint64(recordHeader.RecordNumber)
			directory.DirectoryName = fileNameAttribute.FileName
			directory.ParentRecordNumber = fileNameAttribute.ParentDirRecordNumber
		}
		break
	}
	return
}

func BuildUnresolvedDirectoryTree(reader io.Reader) (unresolvedDirectoryTree UnresolvedDirectoryTree, err error) {
	unresolvedDirectoryTree = make(UnresolvedDirectoryTree)
	for {
		buffer := make(RawMasterFileTableRecord, 1024)
		_, err = reader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}

		directory, err := convertRawMFTRecordToDirectory(buffer)
		if err != nil {
			continue
		}
		unresolvedDirectoryTree[directory.RecordNumber] = directory
	}

	return
}

// Combines a running list of directories from a channel in order to create the systems directory trees.
func (unresolvedDirectoryTree UnresolvedDirectoryTree) resolve() (directoryTree DirectoryTree) {
	directoryTree = make(DirectoryTree)
	for recordNumber, directoryMetadata := range unresolvedDirectoryTree {
		mappingDirectory := directoryMetadata.DirectoryName
		parentRecordNumberPointer := directoryMetadata.ParentRecordNumber
		for {
			if _, ok := unresolvedDirectoryTree[parentRecordNumberPointer]; ok {
				if recordNumber == 5 {
					mappingDirectory = "\\"
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = "\\" + mappingDirectory
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = unresolvedDirectoryTree[parentRecordNumberPointer].DirectoryName + "\\" + mappingDirectory
				parentRecordNumberPointer = unresolvedDirectoryTree[parentRecordNumberPointer].ParentRecordNumber
				continue
			}
			directoryTree[recordNumber] = "$ORPHANFILE\\" + mappingDirectory
			break
		}
	}
	return
}

// Builds a list of directories for the purpose of of mapping MFT records to their parent directories.
func BuildDirectoryTree(reader io.Reader) (directoryTree DirectoryTree, err error) {
	directoryTree = make(DirectoryTree)
	unresolvedDirectoryTree, _ := BuildUnresolvedDirectoryTree(reader)
	directoryTree = unresolvedDirectoryTree.resolve()
	return
}
