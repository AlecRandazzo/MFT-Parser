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
	bin "github.com/AlecRandazzo/BinaryTransforms"
)

// []byte alias for raw record header. Used with the Parse() method.
type RawRecordHeader []byte

// Contains parsed record header values.
type RecordHeader struct {
	AttributesOffset uint16
	RecordNumber     uint32
	Flags            RecordHeaderFlags
}

// byte alias for raw record header flag. Used with the Parse() method.
type RawRecordHeaderFlag byte

// Contains parsed record header flag values.
type RecordHeaderFlags struct {
	FlagDeleted   bool
	FlagDirectory bool
}

// Parses the raw record header receiver and returns a record header.
func (rawRecordHeader RawRecordHeader) Parse() (recordHeader RecordHeader, err error) {
	sizeOfRawRecordHeader := len(rawRecordHeader)

	if sizeOfRawRecordHeader == 0 {
		err = errors.New("RecordHeader.parse() received nil bytes")
		return
	} else if sizeOfRawRecordHeader != 0x38 {
		err = fmt.Errorf("RawRecordHeader.parse() expected 38 bytes, instead it received %v", sizeOfRawRecordHeader)
		return
	}

	const offsetAttributesOffset = 0x14
	const offsetRecordNumber = 0x2C
	const lengthRecordNumber = 0x04

	recordHeader.AttributesOffset = uint16(rawRecordHeader[offsetAttributesOffset])
	rawRecordHeaderFlag, _ := rawRecordHeader.GetRawRecordHeaderFlags()

	recordHeader.Flags = rawRecordHeaderFlag.Parse()
	recordHeader.RecordNumber, _ = bin.LittleEndianBinaryToUInt32(rawRecordHeader[offsetRecordNumber : offsetRecordNumber+lengthRecordNumber])
	return
}

// Parses the raw filename attribute receiver and returns the raw record header flags.
func (rawRecordHeader RawRecordHeader) GetRawRecordHeaderFlags() (rawRecordHeaderFlag RawRecordHeaderFlag, err error) {
	sizeOfRawRecordHeader := len(rawRecordHeader)

	if sizeOfRawRecordHeader == 0 {
		err = errors.New("received a nil bytes")
		return
	} else if sizeOfRawRecordHeader < 0x16 {
		err = fmt.Errorf("expected at least 16 bytes, instead received %v", sizeOfRawRecordHeader)
		return
	}

	const offsetRecordFlag = 0x16
	rawRecordHeaderFlag = RawRecordHeaderFlag(rawRecordHeader[offsetRecordFlag])

	return
}

// Parses the raw record header flag receiver and returns record header flags.
func (rawRecordHeaderFlag RawRecordHeaderFlag) Parse() (recordHeaderFlags RecordHeaderFlags) {
	const codeDeletedFile = 0x00
	//const codeActiveFile = 0x01
	//const codeDeletedDirectory = 0x02
	const codeDirectory = 0x03
	if rawRecordHeaderFlag == codeDeletedFile {
		recordHeaderFlags.FlagDeleted = true
		recordHeaderFlags.FlagDirectory = false
	} else if rawRecordHeaderFlag == codeDirectory {
		recordHeaderFlags.FlagDirectory = true
		recordHeaderFlags.FlagDeleted = false
	} else {
		recordHeaderFlags.FlagDeleted = false
		recordHeaderFlags.FlagDirectory = false
	}
	return
}

// Gets the raw record header from a raw mft record receiver.
func (rawMftRecord RawMasterFileTableRecord) GetRawRecordHeader() (rawRecordHeader RawRecordHeader, err error) {
	sizeOfRawMftRecord := len(rawMftRecord)
	if sizeOfRawMftRecord == 0 {
		err = errors.New("received nil bytes")
		return
	} else if sizeOfRawMftRecord < 0x38 {
		err = fmt.Errorf("expected at least 38 bytes, instead received %v", sizeOfRawMftRecord)
		return
	}

	result, _ := rawMftRecord.IsThisAnMftRecord()
	if result == false {
		err = errors.New("this is not an mft record")
		return
	}

	sizeOfRawRecordHeader := len(rawMftRecord[0:0x38])
	rawRecordHeader = make(RawRecordHeader, sizeOfRawRecordHeader)
	copy(rawRecordHeader, rawMftRecord[0:0x38])
	return
}

// Quick check to see if the raw mft record receiver is a valid mft record.
func (rawMftRecord RawMasterFileTableRecord) IsThisAnMftRecord() (result bool, err error) {
	sizeOfRawMftRecord := len(rawMftRecord)

	if sizeOfRawMftRecord == 0 {
		err = errors.New("received nil bytes")
		result = false
		return
	}
	if sizeOfRawMftRecord < 0x05 {
		err = errors.New("received less than 4 bytes")
		result = false
		return
	}

	const offsetRecordMagicNumber = 0x00
	const lengthRecordMagicNumber = 0x05
	magicNumber := string(rawMftRecord[offsetRecordMagicNumber : offsetRecordMagicNumber+lengthRecordMagicNumber])
	if magicNumber != "FILE0" {
		result = false
		return
	}
	result = true
	return
}
