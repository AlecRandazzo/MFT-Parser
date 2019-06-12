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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMasterFileTableRecord_getAttributeList(t *testing.T) {
	mftBytes, _ := hex.DecodeString("46494C45300003009B1565BC210000000100010038000100C801000000040000000000000000000007000000000000001D07000000000000100000006000000000001800000000004800000018000000EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D10106000000000000000000000000000000000000000001000000000000000000000000000000000000300000006800000000001800000003004A000000180001000500000000000500EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D101004000000000000000400000000000000600000000000000040324004D00460054000000000000008000000078000000010040000000060000000000000000003F3705000000000040000000000000000000745300000000000074530000000000007453000000003320C80000000C436D9401D485E2014336D2006AFA7B0942FD0CF13008F5424563C94EE4084361D100EB51C60143DAC60011E49601000000B000000048000000010040000000050000000000000000002A00000000000000400000000000000000B002000000000008A002000000000008A0020000000000312B67F402000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF0000000008100000000000003101FFFF0B1101FF0000000000001D07FFFFFFFF00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001D07")

	tests := []struct {
		name          string
		mftRecord     *masterFileTableRecord
		wantMftRecord *masterFileTableRecord
	}{
		{
			name: "Test1",
			mftRecord: &masterFileTableRecord{
				MftRecordBytes: mftBytes,
				RecordHeader: recordHeader{
					AttributesOffset: 0x38,
				},
				AttributeInfo: []AttributeInfo{},
			},
			wantMftRecord: &masterFileTableRecord{
				AttributeInfo: []AttributeInfo{
					{
						AttributeType:  16,
						AttributeBytes: []byte{16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.mftRecord.GetAttributeList()
			assert.Equal(t, tt.wantMftRecord.AttributeInfo[0], tt.mftRecord.AttributeInfo[0], "The attribute info should be equal.")
		})
	}
}
