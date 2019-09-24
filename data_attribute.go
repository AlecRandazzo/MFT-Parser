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
	"errors"
	"fmt"
	bin "github.com/AlecRandazzo/BinaryTransforms"
	"strconv"
)

type ResidentDataAttribute struct {
	ResidentData []byte
}

type NonResidentDataAttribute struct {
	DataRuns DataRuns
}

type UnparsedDataRun struct {
	NumberOrder      int
	ClusterOffset    int64
	NumberOfClusters int64
}

type UnparsedDataRuns map[int]UnparsedDataRun

type DataRuns map[int]DataRun

type DataRun struct {
	AbsoluteOffset int64
	Length         int64
}

type DataRunSplit struct {
	offsetByteCount int
	lengthByteCount int
}

type DataAttribute struct {
	TotalSize                 uint8
	FlagResident              bool
	ResidentDataAttributes    ResidentDataAttribute
	NonResidentDataAttributes NonResidentDataAttribute
}

func (dataAttribute *DataAttribute) Parse(attribute Attribute, bytesPerCluster int64) (err error) {
	const offsetResidentFlag = 0x08

	if len(attribute.AttributeBytes) < 0x18 {
		err = errors.New("DataAttribute.Parse() did not receive valid bytes")
		return
	}

	//TODO: handle resident data
	if attribute.AttributeBytes[offsetResidentFlag] == 0x00 {
		dataAttribute.FlagResident = true
		_ = dataAttribute.ResidentDataAttributes.Parse(attribute)
		return
	} else {
		dataAttribute.FlagResident = false
		_ = dataAttribute.NonResidentDataAttributes.Parse(attribute, bytesPerCluster)
	}

	return
}

func (residentDataAttribute *ResidentDataAttribute) Parse(attribute Attribute) (err error) {
	if len(attribute.AttributeBytes) < 0x18 {
		err = errors.New("ResidentDataAttribute.Parse() did not receive valid bytes")
		return
	}

	const offsetResidentData = 0x18
	residentDataAttribute.ResidentData = attribute.AttributeBytes[offsetResidentData:]

	return
}

func (nonResidentDataAttributes *NonResidentDataAttribute) Parse(attribute Attribute, bytesPerCluster int64) (err error) {

	attributeLength := len(attribute.AttributeBytes)
	if attributeLength <= 0x20 {
		err = errors.New("NonResidentDataAttribute.Parse() did not receive valid bytes")
		return
	}

	// Identify offset of the data runs in the data Attribute
	const offsetDataRunOffset = 0x20
	dataRunOffset := attribute.AttributeBytes[offsetDataRunOffset]

	if attributeLength < int(dataRunOffset) {
		err = errors.New("attribute offset longer than length")
		return
	}

	// Pull out the data run bytes
	dataRunsBytes := make([]byte, attributeLength)
	copy(dataRunsBytes, attribute.AttributeBytes[dataRunOffset:])

	// Send the bytes to be parsed
	nonResidentDataAttributes.DataRuns = make(map[int]DataRun)
	_ = nonResidentDataAttributes.DataRuns.Parse(dataRunsBytes, bytesPerCluster)

	return
}

func (dataRuns *DataRuns) Parse(dataRunBytes []byte, bytesPerCluster int64) (err error) {
	if dataRunBytes == nil {
		err = errors.New("dataruns.Parse() received null bytes")
		*dataRuns = nil
		return
	}

	// Initialize a few variables
	UnparsedDataRun := UnparsedDataRun{}
	UnparsedDataRuns := make(UnparsedDataRuns)
	offset := 0
	runCounter := 0

	for {
		// Checks to see if we reached the end of the data runs. If so, break out of the loop.
		if dataRunBytes[offset] == 0x00 || len(dataRunBytes) < offset {
			break
		} else {
			// Take the first byte of a data run and send it to get split so we know how many bytes account for the
			// data run's offset and how many account for the data run's length.
			byteToBeSplit := dataRunBytes[offset]
			dataRunSplit := DataRunSplit{}
			dataRunSplit.Parse(byteToBeSplit)
			offset += 1

			// Pull out the the bytes that account for the data runs offset2 and length
			var lengthBytes, offsetBytes []byte

			lengthBytes = make([]byte, len(dataRunBytes[offset:(offset+dataRunSplit.lengthByteCount)]))
			copy(lengthBytes, dataRunBytes[offset:(offset+dataRunSplit.lengthByteCount)])
			offsetBytes = make([]byte, len(dataRunBytes[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)]))
			copy(offsetBytes, dataRunBytes[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)])

			// Convert the bytes for the data run offset and length to little endian int64
			UnparsedDataRun.ClusterOffset, _ = bin.LittleEndianBinaryToInt64(offsetBytes)
			UnparsedDataRun.NumberOfClusters, _ = bin.LittleEndianBinaryToInt64(lengthBytes)

			// Append the data run to our data run struct
			UnparsedDataRuns[runCounter] = UnparsedDataRun

			// Increment the number order in preparation for the next data run.
			runCounter += 1

			// Set the offset tracker to the position of the next data run
			offset = offset + dataRunSplit.lengthByteCount + dataRunSplit.offsetByteCount
		}
	}

	// Resolve Data Runs
	dataRunOffset := int64(0)
	for i := 0; i < len(UnparsedDataRuns); i++ {
		dataRunOffset = dataRunOffset + (UnparsedDataRuns[i].ClusterOffset * bytesPerCluster)
		(*dataRuns)[i] = DataRun{
			AbsoluteOffset: dataRunOffset,
			Length:         UnparsedDataRuns[i].NumberOfClusters * bytesPerCluster,
		}
	}
	return
}

func (dataRunSplit *DataRunSplit) Parse(dataRun byte) {
	/*
		This function will split the first byte of a data run.
		See the following for a good write up on data runs: https://homepage.cs.uri.edu/~thenry/csc487/video/66_NTFS_Data_Runs.pdf
	*/
	// Convert the byte to a hex string
	hexToSplit := fmt.Sprintf("%x", dataRun)

	// Split the hex string in half and return each half as an int
	dataRunSplit.offsetByteCount, _ = strconv.Atoi(string(hexToSplit[0]))
	dataRunSplit.lengthByteCount, _ = strconv.Atoi(string(hexToSplit[1]))

	return
}
