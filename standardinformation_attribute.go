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
	"encoding/hex"
	"errors"
	"fmt"
)

type StandardInformationAttributes struct {
	SiCreated    TimeStamp
	SiModified   TimeStamp
	SiAccessed   TimeStamp
	SiChanged    TimeStamp
	FlagResident FlagResidency
}

func (standardInfo *StandardInformationAttributes) Parse(attribute attribute) (err error) {
	const offsetResidentFlag = 0x08

	const offsetSiCreated = 0x18
	const lengthSiCreated = 0x08

	const offsetSiModified = 0x20
	const lengthSiModified = 0x08

	const offsetSiChanged = 0x28
	const lengthSiChanged = 0x08

	const offsetSiAccessed = 0x30
	const lengthSiAccessed = 0x08

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to parse standard info attribute")
		}
	}()

	// The standard information attribute has a minimum length of 0x30
	if len(attribute.AttributeBytes) < 0x30 {
		return
	}

	// Check to see if the standard information attribute is resident to the MFT or not
	standardInfo.FlagResident.Parse(attribute.AttributeBytes[offsetResidentFlag])
	if standardInfo.FlagResident == false {
		err = fmt.Errorf("non resident standard information flag found, hex dump: %s", hex.EncodeToString(attribute.AttributeBytes))
		return
	}

	// Parse timestamps
	standardInfo.SiCreated.Parse(attribute.AttributeBytes[offsetSiCreated : offsetSiCreated+lengthSiCreated])
	standardInfo.SiModified.Parse(attribute.AttributeBytes[offsetSiModified : offsetSiModified+lengthSiModified])
	standardInfo.SiChanged.Parse(attribute.AttributeBytes[offsetSiChanged : offsetSiChanged+lengthSiChanged])
	standardInfo.SiAccessed.Parse(attribute.AttributeBytes[offsetSiAccessed : offsetSiAccessed+lengthSiAccessed])
	return
}
