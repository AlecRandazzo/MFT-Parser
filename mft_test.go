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

//func TestRawMasterFileTableRecord_Parse(t *testing.T) {
//	tests := []struct {
//		name          string
//		rawMftRecord  RawMasterFileTableRecord
//		want MasterFileTableRecord
//		got MasterFileTableRecord
//		wantErr       bool
//	}{
//		{
//			name: "test1",
//			rawMftRecord: []byte{70, 73, 76, 69, 48, 0, 3, 0, 113, 250, 76, 78, 8, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 216, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 199, 5, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 81, 3, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 50, 96, 5, 194, 0, 56, 67, 16, 219, 0, 78, 89, 133, 0, 66, 176, 108, 91, 31, 119, 255, 66, 192, 69, 205, 200, 190, 0, 66, 0, 56, 8, 170, 148, 0, 66, 128, 80, 188, 200, 136, 1, 66, 64, 25, 2, 118, 2, 253, 66, 64, 85, 48, 135, 101, 2, 0, 176, 0, 0, 0, 80, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 27, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 192, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 49, 25, 115, 210, 0, 65, 3, 176, 243, 197, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 49, 1, 38, 0, 244, 0, 0, 0, 0, 199, 5, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 199, 5},
//			want: MasterFileTableRecord{
//				BytesPerCluster:               0,
//				RecordHeader:                  RecordHeader{
//					AttributesOffset: 56,
//					RecordNumber:     0,
//					Flags:            RecordHeaderFlags{
//						FlagDeleted:   false,
//						FlagDirectory: false,
//					},
//				},
//				StandardInformationAttributes: StandardInformationAttributes{
//					SiCreated:    ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//					SiModified:   ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//					SiAccessed:   ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//					SiChanged:    ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//					FlagResident: true,
//				},
//				FileNameAttributes:            FileNameAttributes{
//					0: FileNameAttribute{
//						FnCreated:               ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//						FnModified:              ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//						FnAccessed:              ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//						FnChanged:               ts.TimeStamp(time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC)),
//						FlagResident:            true,
//						NameLength:              NameLength{},
//						AttributeSize:           104,
//						ParentDirRecordNumber:   5,
//						ParentDirSequenceNumber: 0,
//						LogicalFileSize:         16384,
//						PhysicalFileSize:        16384,
//						FileNameFlags:           FileNameFlags{
//							ReadOnly:          false,
//							Hidden:            true,
//							System:            true,
//							Archive:           false,
//							Device:            false,
//							Normal:            false,
//							Temporary:         false,
//							Sparse:            false,
//							Reparse:           false,
//							Compressed:        false,
//							Offline:           false,
//							NotContentIndexed: false,
//							Encrypted:         false,
//							Directory:         false,
//							IndexView:         false,
//						},
//						FileNameLength:          8,
//						FileNamespace:           "WIN32 & DOS",
//						FileName:                "$MFT",
//					},
//				},
//				DataAttribute:                DataAttribute{
//					TotalSize:                 0,
//					FlagResident:              false,
//					ResidentDataAttribute:    ResidentDataAttribute{
//						ResidentData:nil,
//					},
//					NonResidentDataAttribute: NonResidentDataAttribute{
//						DataRuns:DataRuns{
//
//						},
//					},
//				},
//				Attributes:                    Attributes{
//					0: Attribute{
//						AttributeType:  16,
//						AttributeBytes: []byte{16,0,0,0,96,0,0,0,0,0,24,0,0,0,0,0,72,0,0,0,24,0,0,0,102,248,4,21,205,173,211,1,102,248,4,21,205,173,211,1,102,248,4,21,205,173,211,1,102,248,4,21,205,173,211,1,6,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
//						AttributeSize:  96,
//					},
//					1: Attribute{
//						AttributeType:  48,
//						AttributeBytes: []byte{48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 0, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
//						AttributeSize:  104,
//					},
//					2: Attribute{
//						AttributeType:  128,
//						AttributeBytes: []byte{128, 0, 0, 0, 128, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 81, 3, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 50, 96, 5, 194, 0, 56, 67, 16, 219, 0, 78, 89, 133, 0, 66, 176, 108, 91, 31, 119, 255, 66, 192, 69, 205, 200, 190, 0, 66, 0, 56, 8, 170, 148, 0, 66, 128, 80, 188, 200, 136, 1, 66, 64, 25, 2, 118, 2, 253, 66, 64, 85, 48, 135, 101, 2, 0},
//						AttributeSize:  128,
//					},
//					3: Attribute{
//						AttributeType:  176,
//						AttributeBytes: []byte{176, 0, 0, 0, 80, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 27, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 192, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 49, 25, 115, 210, 0, 65, 3, 176, 243, 197, 0, 0, 0, 0, 0, 0},
//						AttributeSize:  80,
//					},
//				},
//				UseFullMftFields:              useFullMftFields{},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var err error
//			tt.got, err = tt.rawMftRecord.Parse()
//
//			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
//				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
//			}
//		})
//	}
//}
