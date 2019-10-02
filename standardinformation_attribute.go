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
	ts "github.com/AlecRandazzo/Timestamp-Parser"
)

type RawStandardInformationAttribute []byte

type StandardInformationAttribute struct {
	SiCreated    ts.TimeStamp
	SiModified   ts.TimeStamp
	SiAccessed   ts.TimeStamp
	SiChanged    ts.TimeStamp
	FlagResident FlagResidency
}

func (rawStandardInformationAttribute RawStandardInformationAttribute) Parse() (standardInformationAttribute StandardInformationAttribute, err error) {
	const offsetResidentFlag = 0x08

	const offsetSiCreated = 0x18
	const lengthSiCreated = 0x08

	const offsetSiModified = 0x20
	const lengthSiModified = 0x08

	const offsetSiChanged = 0x28
	const lengthSiChanged = 0x08

	const offsetSiAccessed = 0x30
	const lengthSiAccessed = 0x08

	// The standard information Attribute has a minimum length of 0x30
	if len(rawStandardInformationAttribute) < 0x30 {
		err = errors.New("StandardInformationAttributes.Parse() received invalid bytes")
		return
	}
	// Check to see if the standard information Attribute is resident to the MFT or not
	standardInformationAttribute.FlagResident.Parse(rawStandardInformationAttribute[offsetResidentFlag])
	if standardInformationAttribute.FlagResident == false {
		err = errors.New("non resident standard information flag found")
		return
	}

	// Parse timestamps
	_ = standardInformationAttribute.SiCreated.Parse(rawStandardInformationAttribute[offsetSiCreated : offsetSiCreated+lengthSiCreated])
	_ = standardInformationAttribute.SiModified.Parse(rawStandardInformationAttribute[offsetSiModified : offsetSiModified+lengthSiModified])
	_ = standardInformationAttribute.SiChanged.Parse(rawStandardInformationAttribute[offsetSiChanged : offsetSiChanged+lengthSiChanged])
	_ = standardInformationAttribute.SiAccessed.Parse(rawStandardInformationAttribute[offsetSiAccessed : offsetSiAccessed+lengthSiAccessed])
	return
}
