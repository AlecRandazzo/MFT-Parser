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
	"fmt"
	"io"
	"strconv"
	"sync"
)

type ResultWriter interface {
	ResultWriter(streamer io.Writer, outputChannel *chan UsefulMftFields, waitGroup *sync.WaitGroup)
}

type CsvResultWriter struct{}

func (csvResultWriter *CsvResultWriter) ResultWriter(streamer io.Writer, outputChannel *chan UsefulMftFields, waitGroup *sync.WaitGroup) {
	delimiter := "|"
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
		"FileName Created",
		"FileName Modified",
		"Filename Accessed",
		"Filename Entry Modified",
		"\n",
	}

	// Write CSV header
	headerSize := len(csvHeader)
	for index, header := range csvHeader {
		_, _ = streamer.Write([]byte(header))
		if index < headerSize-2 {
			_, _ = streamer.Write([]byte(delimiter))
		}
	}

	openChannel := true
	for {
		var file UsefulMftFields
		file, openChannel = <-*outputChannel
		if openChannel == false {
			break
		}
		csvRow := []string{
			fmt.Sprint(file.RecordNumber),                  //Record Number
			strconv.FormatBool(file.DirectoryFlag),         //directory Flag
			strconv.FormatBool(file.SystemFlag),            //System file flag
			strconv.FormatBool(file.HiddenFlag),            //Hidden flag
			strconv.FormatBool(file.ReadOnlyFlag),          //Read only flag
			strconv.FormatBool(file.DeletedFlag),           //Deleted Flag
			file.FilePath,                                  //File directory
			file.FileName,                                  //File Name
			strconv.FormatUint(file.PhysicalFileSize, 10),  // File Size
			file.SiCreated.Format("2006-01-02T15:04:05Z"),  //File Created
			file.SiModified.Format("2006-01-02T15:04:05Z"), //File Modified
			file.SiAccessed.Format("2006-01-02T15:04:05Z"), //File Accessed
			file.SiChanged.Format("2006-01-02T15:04:05Z"),  //File entry Modified
			file.FnCreated.Format("2006-01-02T15:04:05Z"),  //FileName Created
			file.FnModified.Format("2006-01-02T15:04:05Z"), //FileName Modified
			file.FnAccessed.Format("2006-01-02T15:04:05Z"), //FileName Accessed
			file.FnChanged.Format("2006-01-02T15:04:05Z"),  //FileName Entry Modified
			"\n", // Newline
		}

		csvRowSize := len(csvRow)
		for index, item := range csvRow {
			_, _ = streamer.Write([]byte(item))
			if index < csvRowSize-2 {
				_, _ = streamer.Write([]byte(delimiter))
			}
		}

	}
	waitGroup.Done()
	return
}
