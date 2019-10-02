package GoFor_MFT_Parser

import (
	"reflect"
	"testing"
)

func TestDataAttribute_Parse(t *testing.T) {
	type args struct {
		attribute       Attribute
		bytesPerCluster int64
	}
	tests := []struct {
		name    string
		got     DataAttribute
		args    args
		want    DataAttribute
		wantErr bool
	}{
		{
			name: "nonresident data attribute test 1",
			args: args{
				attribute: Attribute{
					AttributeType:  0x80,
					AttributeBytes: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
					AttributeSize:  120,
				},
				bytesPerCluster: 4096,
			},
			want: DataAttribute{
				TotalSize:    0,
				FlagResident: false,
				ResidentDataAttribute: ResidentDataAttribute{
					ResidentData: nil,
				},
				NonResidentDataAttribute: NonResidentDataAttribute{
					DataRuns: DataRuns{
						0: {
							AbsoluteOffset: 3221225472,
							Length:         209846272,
						},
						1: {
							AbsoluteOffset: 132747444224,
							Length:         424071168,
						},
						2: {
							AbsoluteOffset: 784502874112,
							Length:         220422144,
						},
						3: {
							AbsoluteOffset: 30787432448,
							Length:         13619200,
						},
						4: {
							AbsoluteOffset: 641829142528,
							Length:         104091648,
						},
						5: {
							AbsoluteOffset: 763784736768,
							Length:         219549696,
						},
						6: {
							AbsoluteOffset: 873008676864,
							Length:         208510976,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "null bytes data attribute",
			args: args{
				attribute: Attribute{
					AttributeType:  0x80,
					AttributeBytes: nil,
					AttributeSize:  120,
				},
				bytesPerCluster: 4096,
			},
			wantErr: true,
		},
		{
			name:    "resident data attribute test 1",
			wantErr: false,
			args: args{
				attribute: Attribute{
					AttributeType:  128,
					AttributeBytes: []byte{128, 0, 0, 0, 136, 0, 0, 0, 0, 0, 24, 0, 0, 0, 1, 0, 106, 0, 0, 0, 24, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					AttributeSize:  136,
				},
				bytesPerCluster: 4096,
			},
			want: DataAttribute{
				TotalSize:    0,
				FlagResident: true,
				ResidentDataAttribute: ResidentDataAttribute{
					ResidentData: []byte{2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				NonResidentDataAttribute: NonResidentDataAttribute{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.got.Parse(tt.args.attribute, tt.args.bytesPerCluster)
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestDataRuns_Parse(t *testing.T) {
	type args struct {
		dataRunBytes    []byte
		bytesPerCluster int64
	}
	tests := []struct {
		name    string
		got     DataRuns
		args    args
		want    DataRuns
		wantErr bool
	}{
		{
			name: "TestDataRuns_Parse test 1",
			args: args{
				dataRunBytes:    []byte{51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				bytesPerCluster: 4096,
			},
			want: DataRuns{
				0: {
					AbsoluteOffset: 3221225472,
					Length:         209846272,
				},
				1: {
					AbsoluteOffset: 132747444224,
					Length:         424071168,
				},
				2: {
					AbsoluteOffset: 784502874112,
					Length:         220422144,
				},
				3: {
					AbsoluteOffset: 30787432448,
					Length:         13619200,
				},
				4: {
					AbsoluteOffset: 641829142528,
					Length:         104091648,
				},
				5: {
					AbsoluteOffset: 763784736768,
					Length:         219549696,
				},
				6: {
					AbsoluteOffset: 873008676864,
					Length:         208510976,
				},
			},
			wantErr: false,
		},
		{
			name:    "null bytes",
			wantErr: true,
			args: args{
				dataRunBytes:    nil,
				bytesPerCluster: 4096,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got = make(map[int]DataRun)
			err := tt.got.Parse(tt.args.dataRunBytes, tt.args.bytesPerCluster)
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestNonResidentDataAttribute_Parse(t *testing.T) {
	type args struct {
		attribute       Attribute
		bytesPerCluster int64
	}
	tests := []struct {
		name    string
		want    NonResidentDataAttribute
		args    args
		got     NonResidentDataAttribute
		wantErr bool
	}{
		{
			name: "TestNonResidentDataAttribute_Parse test 1",
			args: args{
				attribute: Attribute{
					AttributeType:  0x80,
					AttributeBytes: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
					AttributeSize:  120,
				},
				bytesPerCluster: 4096,
			},
			wantErr: false,
			want: NonResidentDataAttribute{
				DataRuns: DataRuns{
					0: {
						AbsoluteOffset: 3221225472,
						Length:         209846272,
					},
					1: {
						AbsoluteOffset: 132747444224,
						Length:         424071168,
					},
					2: {
						AbsoluteOffset: 784502874112,
						Length:         220422144,
					},
					3: {
						AbsoluteOffset: 30787432448,
						Length:         13619200,
					},
					4: {
						AbsoluteOffset: 641829142528,
						Length:         104091648,
					},
					5: {
						AbsoluteOffset: 763784736768,
						Length:         219549696,
					},
					6: {
						AbsoluteOffset: 873008676864,
						Length:         208510976,
					},
				},
			},
		},
		{
			name:    "null bytes in",
			wantErr: true,
			args: args{
				attribute: Attribute{
					AttributeType:  0,
					AttributeBytes: nil,
					AttributeSize:  0,
				},
				bytesPerCluster: 4096,
			},
		},
		{
			name:    "attribute offset longer than length",
			wantErr: true,
			args: args{
				attribute: Attribute{
					AttributeType:  0,
					AttributeBytes: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
					AttributeSize:  120,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.got.Parse(tt.args.attribute, tt.args.bytesPerCluster)
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestResidentDataAttribute_Parse(t *testing.T) {
	type args struct {
		attribute Attribute
	}
	tests := []struct {
		name    string
		want    ResidentDataAttribute
		args    args
		got     ResidentDataAttribute
		wantErr bool
	}{
		{
			name: "TestResidentDataAttribute_Parse test 1",
			args: args{attribute: Attribute{
				AttributeType:  0x80,
				AttributeBytes: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				AttributeSize:  120,
			}},
			want: ResidentDataAttribute{
				ResidentData: []byte{63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
			},
			wantErr: false,
		},
		{
			name:    "null bytes in",
			wantErr: true,
			args: args{
				attribute: Attribute{
					AttributeType:  0,
					AttributeBytes: nil,
					AttributeSize:  0,
				},
			},
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

func Test_dataRunSplit_Parse(t *testing.T) {
	type args struct {
		dataRun byte
	}
	tests := []struct {
		name string
		got  DataRunSplit
		args args
		want DataRunSplit
	}{
		{
			name: "Split 0x33",
			args: args{dataRun: byte(0x33)},
			want: DataRunSplit{
				offsetByteCount: 3,
				lengthByteCount: 3,
			},
		},
		{
			name: "Split 0x04",
			args: args{dataRun: byte(0x04)},
			want: DataRunSplit{
				offsetByteCount: 0,
				lengthByteCount: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.Parse(tt.args.dataRun)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}
