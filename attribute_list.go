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
	bin "github.com/Go-Forensics/BinaryTransforms"
)

// RawAttributeListAttribute is a []byte alias for raw attribute list attributes. Used with the Parse() method
type RawAttributeListAttribute []byte

// AttributeListAttribute contains information about a attribute list attribute
type AttributeListAttribute struct {
	Type                     byte
	MFTReferenceRecordNumber uint32
}

// AttributeListAttributes is a slice of AttributeListAttribute
type AttributeListAttributes []AttributeListAttribute

// Parse a raw attribute list attribute
func (rawAttributeListAttribute RawAttributeListAttribute) Parse() (attributeListAttributes AttributeListAttributes, err error) {
	const offsetAttributeType = 0x00

	const offsetRecordLength = 0x04
	const lengthRecordLength = 0x02

	const offsetFirstSubAttribute = 0x18

	const offsetMFTReferenceRecordNumber = 0x10
	const lengthMFTReferenceRecordNumber = 0x04

	// Sanity checking
	sizeOfRawAttribute := len(rawAttributeListAttribute)
	if sizeOfRawAttribute == 0 {
		err = errors.New("RawAttributeListAttribute.Parse() received nil bytes")
		attributeListAttributes = AttributeListAttributes{}
		return
	} else if rawAttributeListAttribute[offsetAttributeType] != 0x20 {
		err = fmt.Errorf("RawAttributeListAttribute.Parse() receive an attribute thats not an attribute list. Attribute magic number is %x", rawAttributeListAttribute[offsetAttributeType])
		attributeListAttributes = AttributeListAttributes{}
		return
	}

	recordLength, _ := bin.LittleEndianBinaryToUInt16(rawAttributeListAttribute[offsetRecordLength : offsetRecordLength+lengthRecordLength])
	if int(recordLength) != sizeOfRawAttribute {
		err = fmt.Errorf("RawAttributeListAttribute.Parse() received a byte slice thats not equal to the expected attribute length. Size received was %d but expected %d", sizeOfRawAttribute, recordLength)
		attributeListAttributes = AttributeListAttributes{}
		return
	}

	pointerToSubAttribute := offsetFirstSubAttribute
	for pointerToSubAttribute < sizeOfRawAttribute {
		result := isThisAnAttribute(rawAttributeListAttribute[pointerToSubAttribute])
		if result == false {
			return
		}
		attributeListAttribute := AttributeListAttribute{}
		attributeListAttribute.Type = rawAttributeListAttribute[pointerToSubAttribute]
		sizeOfSubAttribute, _ := bin.LittleEndianBinaryToUInt16(rawAttributeListAttribute[pointerToSubAttribute+offsetRecordLength : pointerToSubAttribute+offsetRecordLength+lengthRecordLength])
		attributeListAttribute.MFTReferenceRecordNumber, _ = bin.LittleEndianBinaryToUInt32(rawAttributeListAttribute[pointerToSubAttribute+offsetMFTReferenceRecordNumber : pointerToSubAttribute+offsetMFTReferenceRecordNumber+lengthMFTReferenceRecordNumber])
		attributeListAttributes = append(attributeListAttributes, attributeListAttribute)
		pointerToSubAttribute += int(sizeOfSubAttribute)
	}

	return
}
