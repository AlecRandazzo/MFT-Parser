package GoFor_MFT_Parser

import (
	ts "github.com/AlecRandazzo/Timestamp-Parser"
	"reflect"
	"testing"
	"time"
)

func TestStandardInformationAttributes_Parse(t *testing.T) {
	type args struct {
		attribute Attribute
	}
	tests := []struct {
		name    string
		got     StandardInformationAttributes
		args    args
		want    StandardInformationAttributes
		wantErr bool
	}{
		{
			name:    "test 1",
			wantErr: false,
			args: args{attribute: Attribute{
				AttributeType:  0x10,
				AttributeBytes: []byte{16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 49, 147, 66, 169, 237, 209, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 44, 238, 221, 229, 226, 245, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 253, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 168, 220, 169, 88, 0, 0, 0, 0},
				AttributeSize:  0,
			}},
			want: StandardInformationAttributes{
				SiCreated:    ts.TimeStamp(time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC)),
				SiModified:   ts.TimeStamp(time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC)),
				SiAccessed:   ts.TimeStamp(time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC)),
				SiChanged:    ts.TimeStamp(time.Date(2018, 5, 27, 17, 48, 19, 181726000, time.UTC)),
				FlagResident: true,
			},
		},
		{
			name:    "nil bytes",
			wantErr: true,
			args: args{attribute: Attribute{
				AttributeType:  0x10,
				AttributeBytes: nil,
				AttributeSize:  0,
			}},
		},
		{
			name:    "non-resident",
			wantErr: true,
			args: args{attribute: Attribute{
				AttributeType:  0x10,
				AttributeBytes: []byte{16, 0, 0, 0, 96, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 49, 147, 66, 169, 237, 209, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 44, 238, 221, 229, 226, 245, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 253, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 168, 220, 169, 88, 0, 0, 0, 0},
				AttributeSize:  0,
			}},
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
