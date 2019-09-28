package GoFor_MFT_Parser

import (
	"reflect"
	"testing"
)

func TestRawMasterFileTableRecord_TrimSlackSpace(t *testing.T) {
	tests := []struct {
		name string
		want RawMasterFileTableRecord
		got  RawMasterFileTableRecord
	}{
		{
			name: "test1",
			got:  RawMasterFileTableRecord([]byte{0xba, 0xdb, 0xff, 0xff, 0xff, 0xff, 0x00}),
			want: RawMasterFileTableRecord([]byte{0xba, 0xdb}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.trimSlackSpace()
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}
