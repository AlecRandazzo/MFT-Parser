// Copyright (c) 2020 Alec Randazzo

package main

import (
	"flag"
	"fmt"
	mft "github.com/Go-Forensics/MFT-Parser"
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
	inFileName := flag.String("mft", "", "Input MFT file to parse.")
	outFileName := flag.String("output", "parsed_mft.csv", "Output file.")
	bytesPerCluster := flag.Int64("c", 4096, "Bytes per cluster. This is typically 4096.")
	volumeLetter := flag.String("volume", "", "Volume letter. This will prepend the volume letter to all directory paths.")
	flag.Parse()

	outFile, err := os.Create(*outFileName)
	if err != nil {
		err = fmt.Errorf("failed to create output file %s: %w", *outFileName, err)
		return
	}
	defer outFile.Close()

	inFile, err := os.Open(*inFileName)
	if err != nil {
		err = fmt.Errorf("failed to open file %s: %w", *inFileName, err)
		return
	}
	defer inFile.Close()

	writer := mft.CsvResultWriter{}
	mft.ParseMFT(*volumeLetter, inFile, &writer, outFile, *bytesPerCluster)

}
