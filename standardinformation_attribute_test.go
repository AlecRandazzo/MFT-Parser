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
	"reflect"
	"testing"
	"time"
)

func TestStandardInformationAttributes_Parse(t *testing.T) {
	tests := []struct {
		name                            string
		got                             StandardInformationAttribute
		want                            StandardInformationAttribute
		rawStandardInformationAttribute RawStandardInformationAttribute
		wantErr                         bool
	}{
		{
			name:                            "test 1",
			wantErr:                         false,
			rawStandardInformationAttribute: []byte{16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 49, 147, 66, 169, 237, 209, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 44, 238, 221, 229, 226, 245, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 253, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 168, 220, 169, 88, 0, 0, 0, 0},

			want: StandardInformationAttribute{
				SiCreated:    time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC),
				SiModified:   time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC),
				SiAccessed:   time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC),
				SiChanged:    time.Date(2018, 5, 27, 17, 48, 19, 181726000, time.UTC),
				FlagResident: true,
			},
		},
		{
			name:                            "nil bytes",
			wantErr:                         true,
			rawStandardInformationAttribute: nil,
		},
		{
			name:                            "non-resident",
			wantErr:                         true,
			rawStandardInformationAttribute: []byte{16, 0, 0, 0, 96, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 49, 147, 66, 169, 237, 209, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 44, 238, 221, 229, 226, 245, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 253, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 168, 220, 169, 88, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.got, err = tt.rawStandardInformationAttribute.Parse()
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}
