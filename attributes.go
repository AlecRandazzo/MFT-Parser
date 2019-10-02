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
	"errors"
	"fmt"
)

type RawAttribute []byte
type RawAttributes []RawAttribute

func (rawAttributes RawAttributes) Parse(bytesPerCluster int64) (fileNameAttributes FileNameAttributes, standardInformationAttribute StandardInformationAttribute, dataAttribute DataAttribute, err error) {
	const codeFileName = 0x30
	const codeStandardInformation = 0x10
	const codeData = 0x80
	for _, rawAttribute := range rawAttributes {
		switch rawAttribute[0x00] {
		case codeFileName:
			rawFileNameAttribute := RawFileNameAttribute(make([]byte, len(rawAttribute)))
			copy(rawFileNameAttribute, rawAttribute)
			fileNameAttribute, err := rawFileNameAttribute.Parse()
			if err != nil {
				continue
			}
			fileNameAttributes = append(fileNameAttributes, fileNameAttribute)
		case codeStandardInformation:
			rawStandardInformationAttribute := RawStandardInformationAttribute(make([]byte, len(rawAttribute)))
			copy(rawStandardInformationAttribute, rawAttribute)
			standardInformationAttribute, err = rawStandardInformationAttribute.Parse()
			if err != nil {
				err = fmt.Errorf("failed to get standard info Attribute %w", err)
				return
			}
		case codeData:
			rawDataAttribute := RawDataAttribute(make([]byte, len(rawAttribute)))
			copy(rawDataAttribute, rawAttribute)
			dataAttribute.NonResidentDataAttribute, dataAttribute.ResidentDataAttribute, err = rawDataAttribute.Parse(bytesPerCluster)
			if err != nil {
				err = fmt.Errorf("failed to get data Attribute %w", err)
				return
			}
		}
	}
	return
}

func (rawMftRecord RawMasterFileTableRecord) GetRawAttributes(recordHeader RecordHeader) (rawAttributes RawAttributes, err error) {
	// Doing some sanity checks
	if len(rawMftRecord) == 0 {
		err = errors.New("received nil bytes")
		return
	}
	if recordHeader.AttributesOffset == 0 {
		err = errors.New("record header argument has an attribute offset value of 0")
		return
	}

	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	// Init variable that tracks how far to the next Attribute
	var distanceToNextAttribute uint16 = 0
	offset := recordHeader.AttributesOffset
	sizeOfRawMftRecord := len(rawMftRecord)

	for {
		// Calculate offset to next Attribute
		offset = offset + distanceToNextAttribute

		// Break if the offset is beyond the byte slice
		if offset > uint16(sizeOfRawMftRecord) || offset+0x04 > uint16(sizeOfRawMftRecord) {
			break
		}

		// Verify if the byte slice is actually an MFT Attribute
		shouldWeContinue := isThisAnAttribute(rawMftRecord[offset])
		if shouldWeContinue == false {
			break
		}

		attributeSize := binary.LittleEndian.Uint16(rawMftRecord[offset+offsetAttributeSize : offset+offsetAttributeSize+lengthAttributeSize])
		end := offset + attributeSize

		rawAttribute := RawAttribute(make([]byte, attributeSize))
		copy(rawAttribute, rawMftRecord[offset:end])

		// Append the rawAttributes to the RawAttributes struct
		rawAttributes = append(rawAttributes, rawAttribute)

		// Track the distance to the next Attribute based on the size of the current Attribute
		distanceToNextAttribute = binary.LittleEndian.Uint16(rawMftRecord[offset+offsetAttributeSize : offset+offsetAttributeSize+lengthAttributeSize])
	}

	return
}

func isThisAnAttribute(attributeHeaderToCheck byte) (result bool) {
	// Init a byte slice that tracks all possible valid MFT Attribute types.
	// We'll be used this to verify if what we are looking at is actually an MFT Attribute.
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

	// Verify if the byte slice is actually an MFT Attribute
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
