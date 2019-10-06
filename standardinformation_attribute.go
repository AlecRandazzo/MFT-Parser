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
	"time"
)

type RawStandardInformationAttribute []byte

type StandardInformationAttribute struct {
	SiCreated    time.Time
	SiModified   time.Time
	SiAccessed   time.Time
	SiChanged    time.Time
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
	rawSiCreated := ts.RawTimestamp(rawStandardInformationAttribute[offsetSiCreated : offsetSiCreated+lengthSiCreated])
	rawSiModified := ts.RawTimestamp(rawStandardInformationAttribute[offsetSiModified : offsetSiModified+lengthSiModified])
	rawSiChanged := ts.RawTimestamp(rawStandardInformationAttribute[offsetSiChanged : offsetSiChanged+lengthSiChanged])
	rawSiAccessed := ts.RawTimestamp(rawStandardInformationAttribute[offsetSiAccessed : offsetSiAccessed+lengthSiAccessed])

	standardInformationAttribute.SiCreated, _ = rawSiCreated.Parse()
	standardInformationAttribute.SiModified, _ = rawSiModified.Parse()
	standardInformationAttribute.SiChanged, _ = rawSiChanged.Parse()
	standardInformationAttribute.SiAccessed, _ = rawSiAccessed.Parse()
	return
}
