package GoFor_MFT_Parser

import (
	"os"
	"sync"
	"testing"
)

func TestMftFile_MftToCSV(t *testing.T) {
	type fields struct {
		FileHandle        *os.File
		MappedDirectories map[uint64]string
		OutputChannel     chan MasterFileTableRecord
	}
	type args struct {
		outFileName string
		waitgroup   *sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
