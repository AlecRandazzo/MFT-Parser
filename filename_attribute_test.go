package GoFor_MFT_Parser

import (
	ts "github.com/AlecRandazzo/Timestamp-Parser"
	"reflect"
	"testing"
	"time"
)

func TestFileNameAttribute_Parse(t *testing.T) {
	type args struct {
		attribute Attribute
	}
	tests := []struct {
		name    string
		got     FileNameAttribute
		args    args
		want    FileNameAttribute
		wantErr bool
	}{
		{
			name:    "TestFileNameAttribute_Parse test 1",
			wantErr: false,
			args: args{attribute: Attribute{
				AttributeType:  0x30,
				AttributeBytes: []byte{48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
				AttributeSize:  120,
			}},
			want: FileNameAttribute{
				FnCreated:    ts.TimeStamp(time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC)),
				FnModified:   ts.TimeStamp(time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC)),
				FnAccessed:   ts.TimeStamp(time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC)),
				FnChanged:    ts.TimeStamp(time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC)),
				FlagResident: true,
				NameLength: NameLength{
					FlagNamed: false,
					NamedSize: 0,
				},
				AttributeSize:           104,
				ParentDirRecordNumber:   5,
				ParentDirSequenceNumber: 0,
				LogicalFileSize:         16384,
				PhysicalFileSize:        16384,
				FileNameFlags: FileNameFlags{
					ReadOnly:          false,
					Hidden:            true,
					System:            true,
					Archive:           false,
					Device:            false,
					Normal:            false,
					Temporary:         false,
					Sparse:            false,
					Reparse:           false,
					Compressed:        false,
					Offline:           false,
					NotContentIndexed: false,
					Encrypted:         false,
					Directory:         false,
					IndexView:         false,
				},
				FileNameLength: 8,
				FileNamespace:  "WIN32 & DOS",
				FileName:       "$MFT",
			},
		},
		{
			name:    "null bytes in",
			wantErr: true,
			args: args{attribute: Attribute{
				AttributeType:  0x30,
				AttributeBytes: nil,
				AttributeSize:  120,
			}},
		},
		{
			name:    "non-resident",
			wantErr: true,
			args: args{attribute: Attribute{
				AttributeType:  0x30,
				AttributeBytes: []byte{48, 0, 0, 0, 104, 0, 0, 1, 1, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
				AttributeSize:  120,
			}},
			want: FileNameAttribute{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.got.Parse(tt.args.attribute)
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestFileNameFlags_Parse(t *testing.T) {
	type args struct {
		flagBytes []byte
	}
	tests := []struct {
		name string
		got  FileNameFlags
		args args
		want FileNameFlags
	}{
		{
			name: "flag test 1",
			args: args{
				flagBytes: []byte{6, 0, 0, 0},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            true,
				System:            true,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name: "flag test 2",
			args: args{
				flagBytes: []byte{1, 0, 0, 16},
			},
			want: FileNameFlags{
				ReadOnly:          true,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name: "flag test 3",
			args: args{
				flagBytes: []byte{32, 0, 0, 0},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name: "flag test 4",
			args: args{
				flagBytes: []byte{32, 33, 0, 0},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         true,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: true,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name: "flag test 5",
			args: args{
				flagBytes: []byte{32, 6, 0, 0},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            true,
				Reparse:           true,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name: "flag test 6",
			args: args{
				flagBytes: []byte{6, 36, 0, 16},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            true,
				System:            true,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           true,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: true,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name: "flag test 7",
			args: args{
				flagBytes: []byte{0, 8, 0, 16},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        true,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name: "flag test 8",
			args: args{
				flagBytes: []byte{32, 18, 64, 0},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            true,
				Reparse:           false,
				Compressed:        false,
				Offline:           true,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name: "flag test 9",
			args: args{
				flagBytes: []byte{0, 32, 0, 16},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: true,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name: "flag test 10",
			args: args{
				flagBytes: []byte{0, 0, 0, 16},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name: "flag test 10",
			args: args{
				flagBytes: []byte{38, 0, 0, 32},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            true,
				System:            true,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         true,
			},
		},
		{
			name: "flag test 11",
			args: args{
				flagBytes: []byte{0x40, 0x40, 0x10, 0x10},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            true,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         true,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name: "flag test 12",
			args: args{
				flagBytes: []byte{0x80, 0x00, 0x00, 0x00},
			},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            true,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.Parse(tt.args.flagBytes)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestFlagResidency_Parse(t *testing.T) {
	type args struct {
		byteToCheck byte
	}
	tests := []struct {
		name          string
		flagResidency FlagResidency
		args          args
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_identifyFileNamespace(t *testing.T) {
	type args struct {
		fileNamespaceFlag byte
	}
	tests := []struct {
		name              string
		args              args
		wantFileNameSpace string
	}{
		{
			name:              "POSIX",
			args:              args{fileNamespaceFlag: 0x00},
			wantFileNameSpace: "POSIX",
		},
		{
			name:              "WIN32",
			args:              args{fileNamespaceFlag: 0x01},
			wantFileNameSpace: "WIN32",
		},
		{
			name:              "DOS",
			args:              args{fileNamespaceFlag: 0x02},
			wantFileNameSpace: "DOS",
		},
		{
			name:              "WIN32 & DOS",
			args:              args{fileNamespaceFlag: 0x03},
			wantFileNameSpace: "WIN32 & DOS",
		},
		{
			name:              "null",
			args:              args{fileNamespaceFlag: 0x04},
			wantFileNameSpace: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFileNameSpace := identifyFileNamespace(tt.args.fileNamespaceFlag); gotFileNameSpace != tt.wantFileNameSpace {
				t.Errorf("identifyFileNamespace() = %v, want %v", gotFileNameSpace, tt.wantFileNameSpace)
			}
		})
	}
}
