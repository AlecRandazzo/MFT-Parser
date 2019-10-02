package GoFor_MFT_Parser

import (
	"reflect"
	"testing"
)

func TestRawRecordHeaderFlag_Parse(t *testing.T) {
	tests := []struct {
		name                string
		got                 RecordHeaderFlags
		want                RecordHeaderFlags
		rawRecordHeaderFlag RawRecordHeaderFlag
	}{
		{
			name:                "deleted file 0x00",
			rawRecordHeaderFlag: 0x00,
			want: RecordHeaderFlags{
				FlagDeleted:   true,
				FlagDirectory: false,
			},
		},
		{
			name:                "directory 0x03",
			rawRecordHeaderFlag: 0x03,
			want: RecordHeaderFlags{
				FlagDeleted:   false,
				FlagDirectory: true,
			},
		},
		{
			name:                "other 0x06",
			rawRecordHeaderFlag: 0x06,
			want: RecordHeaderFlags{
				FlagDeleted:   false,
				FlagDirectory: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got = tt.rawRecordHeaderFlag.Parse()
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestRawRecordHeader_Parse(t *testing.T) {
	tests := []struct {
		name            string
		rawRecordHeader RawRecordHeader
		got             RecordHeader
		wantErr         bool
		want            RecordHeader
	}{
		{
			name:            "valid raw record header",
			rawRecordHeader: RawRecordHeader([]byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0}),
			want: RecordHeader{
				AttributesOffset: 56,
				RecordNumber:     0,
				Flags: RecordHeaderFlags{
					FlagDeleted:   false,
					FlagDirectory: false,
				},
			},
			wantErr: false,
		},
		{
			name:            "nil bytes",
			rawRecordHeader: nil,
			wantErr:         true,
		},
		{
			name:            "bytes does not start with FILE0",
			rawRecordHeader: RawRecordHeader([]byte{0, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0}),
			wantErr:         true,
		},
		{
			name:            "not 38 bytes",
			rawRecordHeader: RawRecordHeader([]byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0}),
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.got, err = tt.rawRecordHeader.Parse()
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestRawMasterFileTableRecord_GetRawRecordHeader(t *testing.T) {
	tests := []struct {
		name         string
		rawMftRecord RawMasterFileTableRecord
		want         RawRecordHeader
		got          RawRecordHeader
		wantErr      bool
	}{
		{
			name:         "nil bytes",
			rawMftRecord: nil,
			wantErr:      true,
		},
		{
			name:         "not enough bytes",
			rawMftRecord: RawMasterFileTableRecord([]byte{0x01}),
			wantErr:      true,
		},
		{
			name:         "valid rawMftRecord",
			rawMftRecord: RawMasterFileTableRecord([]byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0, 176, 0, 0, 0, 72, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 176, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 49, 43, 103, 244, 2, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 17, 1, 255, 0, 0, 0, 0, 0, 0, 29, 7, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 29, 7}),
			wantErr:      false,
			want:         RawRecordHeader([]byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.got, err = tt.rawMftRecord.GetRawRecordHeader()
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestRawRecordHeader_GetRawRecordHeaderFlags(t *testing.T) {
	tests := []struct {
		name            string
		rawRecordHeader RawRecordHeader
		got             RawRecordHeaderFlag
		want            RawRecordHeaderFlag
		wantErr         bool
	}{
		{
			name:            "valid raw record header",
			rawRecordHeader: RawRecordHeader([]byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0}),
			wantErr:         false,
			want:            RawRecordHeaderFlag(0x01),
		},
		{
			name:            "nil bytes",
			rawRecordHeader: nil,
			wantErr:         true,
		},
		{
			name:            "not enough bytes",
			rawRecordHeader: RawRecordHeader([]byte{0x00}),
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.got, err = tt.rawRecordHeader.GetRawRecordHeaderFlags()
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}
