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
	bin "github.com/AlecRandazzo/BinaryTransforms"
)

type RecordHeader struct {
	AttributesOffset uint16
	RecordNumber     uint32
	Flags            RecordHeaderFlags
}

type RecordHeaderFlags struct {
	FlagDeleted   bool
	FlagDirectory bool
}

func (recordHeader *RecordHeader) Parse(mftRecord []byte) (err error) {
	if len(mftRecord) == 0 {
		err = errors.New("RecordHeader.Parse() received nil bytes")
		return
	}

	const offsetRecordMagicNumber = 0x00
	const lengthRecordMagicNumber = 0x05
	magicNumber := string(mftRecord[offsetRecordMagicNumber : offsetRecordMagicNumber+lengthRecordMagicNumber])
	if magicNumber != "FILE0" {
		err = errors.New("mftrecord missing magic number FILE0")
		return
	}

	const offsetAttributesOffset = 0x14
	const offsetRecordNumber = 0x2c
	const lengthRecordNumber = 0x04
	const offsetRecordFlag = 0x16

	recordHeader.AttributesOffset = uint16(mftRecord[offsetAttributesOffset])
	recordHeader.Flags.Parse(mftRecord[offsetRecordFlag])
	recordHeader.RecordNumber, _ = bin.LittleEndianBinaryToUInt32(mftRecord[offsetRecordNumber : offsetRecordNumber+lengthRecordNumber])
	return
}

func (recordHeaderFlags *RecordHeaderFlags) Parse(inByte byte) {
	const codeDeletedFile = 0x00
	//const codeActiveFile = 0x01
	//const codeDeletedDirectory = 0x02
	const codeDirectory = 0x03
	if inByte == codeDeletedFile {
		recordHeaderFlags.FlagDeleted = true
		recordHeaderFlags.FlagDirectory = false
	} else if inByte == codeDirectory {
		recordHeaderFlags.FlagDirectory = true
		recordHeaderFlags.FlagDeleted = false
	} else {
		recordHeaderFlags.FlagDeleted = false
		recordHeaderFlags.FlagDirectory = false
	}
	return
}
