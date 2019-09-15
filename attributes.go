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
	"fmt"
)

type Attributes []attribute

type attribute struct {
	AttributeType  byte
	AttributeBytes []byte
	AttributeSize  uint16
}

// Parse MFT record attributes list.
func (attributes *Attributes) Parse(mftRecord []byte, attributesOffset uint16) (err error) {
	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	//const offsetResidentFlag = 0x08
	//const offsetHeaderNameLength = 0x09
	//
	//const offsetAttributeType = 0x18
	//const lengthAttributeType = 0x04
	//
	//const offsetRecordLength = 0x1c
	//const lengthRecordLength = 0x02
	//
	//const offsetNameLength = 0x1e
	//const offsetNameOffset = 0x1f
	//
	//const offsetStartingVCN = 0x20
	//const lengthStartingVCN = 0x08
	//
	//const offsetBaseFileReference = 0x28
	//const lengthBaseFileReference = 0x08
	//
	//const offsetAttributeId = 0x30
	//const lengthAttributeId = 0x02
	//
	//const offsetName = 0x32

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %s, hex dump: %s", fmt.Sprint(r), hex.EncodeToString(mftRecord))
		}
	}()

	// Init variable that tracks how far to the next attribute
	var distanceToNextAttribute uint16 = 0

	for {
		// Calculate offset to next attribute
		attributesOffset = attributesOffset + distanceToNextAttribute

		// Break if the offset is beyond the byte slice
		lenBytesIn := len(mftRecord)
		if attributesOffset > uint16(lenBytesIn) || attributesOffset+0x04 > uint16(lenBytesIn) {
			break
		}

		// Verify if the byte slice is actually an MFT attribute
		shouldWeContinue := isThisAnAttribute(mftRecord[attributesOffset])
		if shouldWeContinue == false {
			break
		}

		// Pull out information describing the attribute and the attribute bytes
		attribute := attribute{}
		attribute.Parse(mftRecord, attributesOffset)

		// Append the attribute to the attribute struct
		*attributes = append(*attributes, attribute)

		// Track the distance to the next attribute based on the size of the current attribute
		distanceToNextAttribute = attribute.AttributeSize
		if distanceToNextAttribute == 0 {
			break
		}
	}
	return
}

func (attribute *attribute) Parse(mftRecord []byte, attributeOffset uint16) {
	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	attribute.AttributeType = mftRecord[attributeOffset]
	attribute.AttributeSize = binary.LittleEndian.Uint16(mftRecord[attributeOffset+offsetAttributeSize : attributeOffset+offsetAttributeSize+lengthAttributeSize])
	end := attributeOffset + attribute.AttributeSize
	copy(attribute.AttributeBytes, mftRecord[attributeOffset:end])

	return
}

//TODO write a unit test for isThisAnAttribute()
func isThisAnAttribute(attributeHeaderToCheck byte) (result bool) {
	// Init a byte slice that tracks all possible valid MFT attribute types.
	// We'll be used this to verify if what we are looking at is actually an MFT attribute.
	const codeStandardInformation = 0x10
	const codeAttributeList = 0x20
	const codeFileName = 0x30
	const codeVolumeVersion = 0x40
	const codeSecurityDescriptor = 0x50
	const codeVolumeName = 0x60
	const codeVolumeInformation = 0x70
	const codeData = 0x80
	const codeIndexRoot = 0x90
	const codeIndexAllocation = 0xA0
	const codeBitmap = 0xB0
	const codeSymbolicLink = 0xC0
	const codeReparsePoint = 0xD0
	const codeEaInformation = 0xE0
	const codePropertySet = 0xF0

	validAttributeTypes := []byte{
		codeStandardInformation,
		codeAttributeList,
		codeFileName,
		codeVolumeVersion,
		codeSecurityDescriptor,
		codeVolumeName,
		codeVolumeInformation,
		codeData,
		codeIndexRoot,
		codeIndexAllocation,
		codeBitmap,
		codeSymbolicLink,
		codeReparsePoint,
		codeEaInformation,
		codePropertySet,
	}

	// Verify if the byte slice is actually an MFT attribute
	for _, validType := range validAttributeTypes {
		if attributeHeaderToCheck == validType {
			result = true
			break
		} else {
			result = false
		}
	}

	return
}
