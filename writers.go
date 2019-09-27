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
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

type OutputWriters interface {
	Write(outputChannel *chan UseFullMftFields, waitGroup *sync.WaitGroup) (err error)
}

type CsvWriter struct {
	OutFile io.Writer
}

func (writer CsvWriter) Write(outputChannel *chan UseFullMftFields, waitGroup *sync.WaitGroup) (err error) {

	csvWriter := csv.NewWriter(writer.OutFile)
	csvWriter.Comma = '|'
	csvHeader := []string{
		"Record Number",
		"directory Flag",
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
		var file UseFullMftFields
		file, openChannel = <-*outputChannel
		csvRow = []string{
			fmt.Sprint(file.RecordNumber),                             //Record Number
			strconv.FormatBool(file.DirectoryFlag),                    //directory Flag
			strconv.FormatBool(file.SystemFlag),                       //System file flag
			strconv.FormatBool(file.HiddenFlag),                       //Hidden flag
			strconv.FormatBool(file.ReadOnlyFlag),                     //Read only flag
			strconv.FormatBool(file.DeletedFlag),                      //Deleted Flag
			file.FilePath,                                             //File directory
			file.FileName,                                             //File Name
			strconv.FormatUint(file.PhysicalFileSize, 10),             // File Size
			time.Time(file.SiCreated).Format("2006-01-02T15:04:05Z"),  //File Created
			time.Time(file.SiModified).Format("2006-01-02T15:04:05Z"), //File Modified
			time.Time(file.SiAccessed).Format("2006-01-02T15:04:05Z"), //File Accessed
			time.Time(file.SiChanged).Format("2006-01-02T15:04:05Z"),  //File entry Modified
			time.Time(file.FnCreated).Format("2006-01-02T15:04:05Z"),  //FileName Created
			time.Time(file.FnModified).Format("2006-01-02T15:04:05Z"), //FileName Modified
			time.Time(file.FnAccessed).Format("2006-01-02T15:04:05Z"), //FileName Accessed
			time.Time(file.FnChanged).Format("2006-01-02T15:04:05Z"),  //FileName Entry Modified
		}
		err = csvWriter.Write(csvRow)
		if err != nil {
			log.Fatal(err)
		}
		break
	}
	waitGroup.Done()
	return
}
