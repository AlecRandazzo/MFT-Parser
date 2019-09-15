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
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	bin "github.com/AlecRandazzo/BinaryTransforms"
)

type FileNameAttributes []FileNameAttribute

type FlagResidency bool

type FileNameAttribute struct {
	FnCreated               TimeStamp
	FnModified              TimeStamp
	FnAccessed              TimeStamp
	FnChanged               TimeStamp
	FlagResident            FlagResidency
	NameLength              NameLength
	AttributeSize           uint32
	ParentDirRecordNumber   uint64
	ParentDirSequenceNumber uint16
	LogicalFileSize         uint64
	PhysicalFileSize        uint64
	FileNameFlags           FileNameFlags
	FileNameLength          byte
	FileNamespace           string
	FileName                string
}

type NameLength struct {
	FlagNamed bool
	NamedSize byte
}

type FileNameFlags struct {
	ReadOnly          bool
	Hidden            bool
	System            bool
	Archive           bool
	Device            bool
	Normal            bool
	Temporary         bool
	Sparse            bool
	Reparse           bool
	Compressed        bool
	Offline           bool
	NotContentIndexed bool
	Encrypted         bool
	Directory         bool
	IndexView         bool
}

func (filenameAttribute *FileNameAttribute) Parse(attribute attribute) (err error) {
	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	const offsetResidentFlag = 0x08

	const offsetParentRecordNumber = 0x18
	const lengthParentRecordNumber = 0x06

	const offsetParentDirSequenceNumber = 0x1e
	const lengthParentDirSequenceNumber = 0x02

	const offsetFnCreated = 0x20
	const lengthFnCreated = 0x08

	const offsetFnModified = 0x28
	const lengthFnModified = 0x08

	const offsetFnChanged = 0x30
	const lengthFnChanged = 0x08

	const offsetFnAccessed = 0x38
	const lengthFnAccessed = 0x08

	const offsetLogicalFileSize = 0x40
	const lengthLogicalFileSize = 0x08

	const offSetPhysicalFileSize = 0x48
	const lengthPhysicalFileSize = 0x08

	const offsetFnFlags = 0x50
	const lengthFnFlags = 0x04

	const offsetFileNameLength = 0x58
	const offsetFileNameSpace = 0x59
	const offsetFileName = 0x5a

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to parse filename attribute")
		}
	}()

	// The filename attribute has a minimum length of 0x44
	if len(attribute.AttributeBytes) < 0x44 {
		return
	}
	filenameAttribute.AttributeSize, _ = bin.LittleEndianBinaryToUInt32(attribute.AttributeBytes[offsetAttributeSize : offsetAttributeSize+lengthAttributeSize])
	filenameAttribute.FlagResident.Parse(attribute.AttributeBytes[offsetResidentFlag])
	if filenameAttribute.FlagResident == false {
		err = fmt.Errorf("parseFileNameAttribute(): non-resident filename attribute encountered, hex dump: %s", hex.EncodeToString(attribute.AttributeBytes))
		return
	}

	filenameAttribute.ParentDirRecordNumber = bin.LittleEndianBinaryToUInt64(attribute.AttributeBytes[offsetParentRecordNumber : offsetParentRecordNumber+lengthParentRecordNumber])
	filenameAttribute.ParentDirSequenceNumber, _ = bin.LittleEndianBinaryToUInt16(attribute.AttributeBytes[offsetParentDirSequenceNumber : offsetParentDirSequenceNumber+lengthParentDirSequenceNumber])
	filenameAttribute.FnCreated.Parse(attribute.AttributeBytes[offsetFnCreated : offsetFnCreated+lengthFnCreated])
	filenameAttribute.FnModified.Parse(attribute.AttributeBytes[offsetFnModified : offsetFnModified+lengthFnModified])
	filenameAttribute.FnChanged.Parse(attribute.AttributeBytes[offsetFnChanged : offsetFnChanged+lengthFnChanged])
	filenameAttribute.FnAccessed.Parse(attribute.AttributeBytes[offsetFnAccessed : offsetFnAccessed+lengthFnAccessed])
	filenameAttribute.LogicalFileSize = bin.LittleEndianBinaryToUInt64(attribute.AttributeBytes[offsetLogicalFileSize : offsetLogicalFileSize+lengthLogicalFileSize])
	filenameAttribute.PhysicalFileSize = bin.LittleEndianBinaryToUInt64(attribute.AttributeBytes[offSetPhysicalFileSize : offSetPhysicalFileSize+lengthPhysicalFileSize])
	filenameAttribute.FileNameFlags.Parse(attribute.AttributeBytes[offsetFnFlags : offsetFnFlags+lengthFnFlags])
	filenameAttribute.FileNameLength = attribute.AttributeBytes[offsetFileNameLength] * 2 // times two to account for unicode characters
	filenameAttribute.FileNamespace = identifyFileNamespace(attribute.AttributeBytes[offsetFileNameSpace])
	filenameAttribute.FileName = bin.UnicodeBytesToASCII(attribute.AttributeBytes[offsetFileName : offsetFileName+filenameAttribute.FileNameLength])
	return
}

func (flagResidency *FlagResidency) Parse(byteToCheck byte) {
	switch byteToCheck {
	case 0x00:
		*flagResidency = true
	default:
		*flagResidency = false
	}
	return
}

func identifyFileNamespace(fileNamespaceFlag byte) (fileNameSpace string) {
	switch fileNamespaceFlag {
	case 0x00:
		fileNameSpace = "POSIX"
	case 0x01:
		fileNameSpace = "WIN32"
	case 0x02:
		fileNameSpace = "DOS"
	case 0x03:
		fileNameSpace = "WIN32 & DOS"
	default:
		fileNameSpace = ""
	}

	return
}

func (fileNameFlags *FileNameFlags) Parse(flagBytes []byte) {
	unparsedFlags := binary.LittleEndian.Uint32(flagBytes)

	//init values
	fileNameFlags.ReadOnly = false
	fileNameFlags.Hidden = false
	fileNameFlags.System = false
	fileNameFlags.Archive = false
	fileNameFlags.Device = false
	fileNameFlags.Normal = false
	fileNameFlags.Temporary = false
	fileNameFlags.Sparse = false
	fileNameFlags.Reparse = false
	fileNameFlags.Compressed = false
	fileNameFlags.Offline = false
	fileNameFlags.NotContentIndexed = false
	fileNameFlags.Encrypted = false
	fileNameFlags.Directory = false
	fileNameFlags.IndexView = false

	if unparsedFlags&0x0001 != 0 {
		fileNameFlags.ReadOnly = true
	}
	if unparsedFlags&0x0002 != 0 {
		fileNameFlags.Hidden = true
	}
	if unparsedFlags&0x0004 != 0 {
		fileNameFlags.System = true
	}
	if unparsedFlags&0x0010 != 0 {
		fileNameFlags.Directory = true
	}
	if unparsedFlags&0x0020 != 0 {
		fileNameFlags.Archive = true
	}
	if unparsedFlags&0x0040 != 0 {
		fileNameFlags.Device = true
	}
	if unparsedFlags&0x0080 != 0 {
		fileNameFlags.Normal = true
	}
	if unparsedFlags&0x0100 != 0 {
		fileNameFlags.Temporary = true
	}
	if unparsedFlags&0x0200 != 0 {
		fileNameFlags.Sparse = true
	}
	if unparsedFlags&0x0400 != 0 {
		fileNameFlags.Reparse = true
	}
	if unparsedFlags&0x0800 != 0 {
		fileNameFlags.Compressed = true
	}
	if unparsedFlags&0x1000 != 0 {
		fileNameFlags.Offline = true
	}
	if unparsedFlags&0x2000 != 0 {
		fileNameFlags.NotContentIndexed = true
	}
	if unparsedFlags&0x4000 != 0 {
		fileNameFlags.Encrypted = true
	}
	if unparsedFlags&0x10000000 != 0 {
		fileNameFlags.Directory = true
	}
	if unparsedFlags&0x20000000 != 0 {
		fileNameFlags.IndexView = true
	}
	return
}
