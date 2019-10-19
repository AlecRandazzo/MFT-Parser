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
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type directory struct {
	recordNumber       uint64
	directoryName      string
	parentRecordNumber uint64
}

// UnresolvedDirectoryTree contains a slice of directories that need to be joined to create a directory tree.
type UnresolvedDirectoryTree map[uint64]directory

// DirectoryTree contains a directory tree.
type DirectoryTree map[uint64]string

// IsThisADirectory will quickly check the bytes of an MFT record to determine if it is a directory or not.
func (rawMftRecord RawMasterFileTableRecord) IsThisADirectory() (result bool, err error) {
	// Sanity checks that the method received good data
	const offsetRecordFlag = 0x16
	const codeDirectory = 0x03
	sizeOfRawMFTRecord := len(rawMftRecord)
	if sizeOfRawMFTRecord == 0 {
		result = false
		err = errors.New("RawMasterFileTableRecord.IsThisADirectory() received nil bytes ")
		return
	}
	if sizeOfRawMFTRecord <= offsetRecordFlag {
		result = false
		err = errors.New("RawMasterFileTableRecord.IsThisADirectory() received not enough bytes ")
		return
	}

	// Skip straight to the offset where the directory flag resides and check if it has the directory flag or not.
	recordFlag := rawMftRecord[offsetRecordFlag]
	if recordFlag == codeDirectory {
		result = true
	} else {
		result = false
	}
	return
}

func convertRawMFTRecordToDirectory(rawMftRecord RawMasterFileTableRecord) (directory directory, err error) {
	// Sanity checks that the raw mft record is a directory or not
	result, err := rawMftRecord.IsThisADirectory()
	if result == false {
		err = errors.New("this is not a directory")
		return
	}

	// Get record header bytes
	rawRecordHeader, err := rawMftRecord.GetRawRecordHeader()
	if err != nil {
		err = fmt.Errorf("failed to parse get record header: %w", err)
		return
	}

	// Parse the raw record header
	recordHeader, _ := rawRecordHeader.Parse()

	// Get the raw mft attributes
	rawAttributes, err := rawMftRecord.GetRawAttributes(recordHeader)
	if err != nil {
		err = fmt.Errorf("failed to get raw attributes: %w", err)
		return
	}
	doesntMatter := int64(4096)

	// Find the filename attribute and parse it for its record number, directory name, and parent record number.
	fileNameAttributes, _, _, err := rawAttributes.Parse(doesntMatter)
	for _, fileNameAttribute := range fileNameAttributes {
		if strings.Contains(fileNameAttribute.FileNamespace, "WIN32") == true || strings.Contains(fileNameAttribute.FileNamespace, "POSIX") {
			directory.recordNumber = uint64(recordHeader.RecordNumber)
			directory.directoryName = fileNameAttribute.FileName
			directory.parentRecordNumber = fileNameAttribute.ParentDirRecordNumber
			break
		}
	}
	return
}

// BuildUnresolvedDirectoryTree takes an MFT and does a first pass to find all the directories listed in it. These will form an unresolved directory tree that need to be stitched together.
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
		unresolvedDirectoryTree[directory.recordNumber] = directory
	}

	return
}

// Resolve combines a running list of directories from a channel in order to create the systems directory trees.
func (unresolvedDirectoryTree UnresolvedDirectoryTree) Resolve(volumeLetter string) (directoryTree DirectoryTree, err error) {
	err = volumeLetterCheck(volumeLetter)
	if err != nil {
		err = fmt.Errorf("failed to build directory tree due to invalid volume letter: %w", err)
		return
	}
	directoryTree = make(DirectoryTree)
	for recordNumber, directoryMetadata := range unresolvedDirectoryTree {
		// Sanity check
		if directoryMetadata.directoryName == "" && directoryMetadata.parentRecordNumber == 0 && directoryMetadata.recordNumber == 0 {
			continue
		}

		mappingDirectory := directoryMetadata.directoryName
		parentRecordNumberPointer := directoryMetadata.parentRecordNumber
		for {
			if _, ok := unresolvedDirectoryTree[parentRecordNumberPointer]; ok {
				if recordNumber == 5 {
					mappingDirectory = fmt.Sprintf("%s:\\", volumeLetter)
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = fmt.Sprintf("%s:\\%s", volumeLetter, mappingDirectory)
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = fmt.Sprintf("%s\\%s", unresolvedDirectoryTree[parentRecordNumberPointer].directoryName, mappingDirectory)
				parentRecordNumberPointer = unresolvedDirectoryTree[parentRecordNumberPointer].parentRecordNumber
				continue
			}
			directoryTree[recordNumber] = fmt.Sprintf("%s:\\$ORPHANFILE\\%s", volumeLetter, mappingDirectory)
			break
		}
	}
	return
}

// BuildDirectoryTree takes an MFT and creates a directory tree where the slice keys are the mft record number of the directory. This record number is importable because files will reference it as its parent mft record number.
func BuildDirectoryTree(reader io.Reader, volumeLetter string) (directoryTree DirectoryTree, err error) {
	err = volumeLetterCheck(volumeLetter)
	if err != nil {
		err = fmt.Errorf("failed to build directory tree due to invalid volume letter: %w", err)
		return
	}
	directoryTree = make(DirectoryTree)
	unresolvedDirectoryTree, _ := BuildUnresolvedDirectoryTree(reader)
	directoryTree, _ = unresolvedDirectoryTree.Resolve(volumeLetter)
	return
}

func volumeLetterCheck(volumeLetter string) (err error) {
	volumeLetterRune := []rune(volumeLetter)
	if volumeLetter == "" {
		err = errors.New("volume letter was blank")
		return
	} else if len(volumeLetterRune) != 1 {
		err = errors.New("volume letter contained more than one character")
		return
	} else if !unicode.IsLetter(volumeLetterRune[0]) {
		err = errors.New("volume letter was not a letter")
		return
	}
	err = nil
	return
}
