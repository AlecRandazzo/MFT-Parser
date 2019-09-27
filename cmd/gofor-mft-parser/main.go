/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package main

import (
	"fmt"
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log configuration
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
}

func main() {
	outFileName := "out.csv"
	inFileName := "MFT"

	outFile, err := os.Create(outFileName)
	if err != nil {
		err = fmt.Errorf("failed to create output file %s: %w", outFileName, err)
		return
	}
	defer outFile.Close()

	inFile, err := os.Open(inFileName)
	if err != nil {
		err = fmt.Errorf("failed to open file %s: %w", inFileName, err)
		return
	}
	defer inFile.Close()

	writer := mft.CsvWriter{
		OutFile: outFile,
	}

	err = mft.ParseMFT(inFile, writer, 4)
	if err != nil {
		log.Fatal(err)
	}

}
