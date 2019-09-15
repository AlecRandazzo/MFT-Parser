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
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func (file MftFile) MftToCSV(outFileName string, waitgroup *sync.WaitGroup) (err error) {
	outFile, err := os.Create(outFileName)
	if err != nil {
		err = fmt.Errorf("ParseMFT(): failed to create output file %s: %w", outFileName, err)
	}
	defer outFile.Close()
	csvWriter := csv.NewWriter(outFile)
	csvWriter.Comma = '|'
	csvHeader := []string{
		"Record Number",
		"Directory Flag",
		"System File Flag",
		"Hidden Flag",
		"Read-only Flag",
		"Deleted Flag",
		"File Path",
		"File Name",
		"File Size",
		"File Created",
		"File Modified",
		"File Accessed",
		"File Entry Modified",
		"FileName Created ",
		"FileName Modified ",
		"Filename Accessed ",
		"Filename Entry Modified ",
	}
	err = csvWriter.Write(csvHeader)
	if err != nil {
		log.Fatal(err)
	}

	openChannel := true
	for openChannel != false {
		var csvRow []string
		var mftRecord MasterFileTableRecord
		mftRecord, openChannel = <-file.OutputChannel
		for _, record := range mftRecord.FileNameAttributes {
			if strings.Contains(record.FileNamespace, "WIN32") == true || strings.Contains(record.FileNamespace, "POSIX") {
				var fileDirectory string
				if value, ok := file.MappedDirectories[record.ParentDirRecordNumber]; ok {
					fileDirectory = value
				} else {
					fileDirectory = "$ORPHANFILE"
				}
				recordNumber := fmt.Sprint(mftRecord.RecordHeader.RecordNumber)
				physicalFileSize := strconv.FormatUint(record.PhysicalFileSize, 10)
				csvRow = []string{
					recordNumber, //Record Number
					strconv.FormatBool(mftRecord.RecordHeader.Flags.FlagDirectory), //Directory Flag
					strconv.FormatBool(record.FileNameFlags.System),                //System file flag
					strconv.FormatBool(record.FileNameFlags.Hidden),                //Hidden flag
					strconv.FormatBool(record.FileNameFlags.ReadOnly),              //Read only flag
					strconv.FormatBool(mftRecord.RecordHeader.Flags.FlagDeleted),   //Deleted Flag
					fileDirectory,    //File Directory
					record.FileName,  //File Name
					physicalFileSize, // File Size
					string(mftRecord.StandardInformationAttributes.SiCreated),  //File Created
					string(mftRecord.StandardInformationAttributes.SiModified), //File Modified
					string(mftRecord.StandardInformationAttributes.SiAccessed), //File Accessed
					string(mftRecord.StandardInformationAttributes.SiChanged),  //File entry Modified
					string(record.FnCreated),                                   //FileName Created
					string(record.FnModified),                                  //FileName Modified
					string(record.FnAccessed),                                  //FileName Accessed
					string(record.FnChanged),                                   //FileName Entry Modified
				}
				err = csvWriter.Write(csvRow)
				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}
	}
	waitgroup.Done()
	return
}
