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

type RawRecordHeader []byte

type RecordHeader struct {
	AttributesOffset uint16
	RecordNumber     uint32
	Flags            RecordHeaderFlags
}

type RawRecordHeaderFlag byte

type RecordHeaderFlags struct {
	FlagDeleted   bool
	FlagDirectory bool
}

func (rawRecordHeader RawRecordHeader) Parse() (recordHeader RecordHeader, err error) {
	sizeOfRawRecordHeader := len(rawRecordHeader)

	if sizeOfRawRecordHeader == 0 {
		err = errors.New("RecordHeader.Parse() received nil bytes")
		return
	} else if sizeOfRawRecordHeader != 0x38 {
		err = fmt.Errorf("RawRecordHeader.Parse() expected 38 bytes, instead it received %v", sizeOfRawRecordHeader)
		return
	}

	const offsetRecordMagicNumber = 0x00
	const lengthRecordMagicNumber = 0x05
	magicNumber := string(rawRecordHeader[offsetRecordMagicNumber : offsetRecordMagicNumber+lengthRecordMagicNumber])
	if magicNumber != "FILE0" {
		err = errors.New("mftrecord missing magic number FILE0")
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

func (rawMftRecord RawMasterFileTableRecord) GetRawRecordHeader() (rawRecordHeader RawRecordHeader, err error) {
	sizeOfRawMftRecord := len(rawMftRecord)
	if sizeOfRawMftRecord == 0 {
		err = errors.New("received nil bytes")
		return
	} else if sizeOfRawMftRecord < 0x38 {
		err = fmt.Errorf("expected at least 38 bytes, instead received %v", sizeOfRawMftRecord)
		return
	}
	sizeOfRawRecordHeader := len(rawMftRecord[0:0x38])
	rawRecordHeader = make(RawRecordHeader, sizeOfRawRecordHeader)
	copy(rawRecordHeader, rawMftRecord[0:0x38])
	return
}
