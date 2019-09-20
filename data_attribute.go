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
	"encoding/hex"
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
	const codeData = 0x80
	const offsetResidentFlag = 0x08

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to parse data Attribute")
		}
	}()

	if len(attribute.AttributeBytes) < 0x18 {
		return
	}

	//TODO: handle resident data
	if attribute.AttributeBytes[offsetResidentFlag] == 0x00 {
		dataAttribute.FlagResident = true
		err = dataAttribute.ResidentDataAttributes.Parse(attribute)
		if err != nil {
			err = fmt.Errorf("failed to parse resident data Attribute: %w", err)
			return
		}
		return
	} else {
		dataAttribute.FlagResident = false
		err = dataAttribute.NonResidentDataAttributes.Parse(attribute, bytesPerCluster)
		if err != nil {
			err = fmt.Errorf("failed to parse non resident data Attribute: %w", err)
			return
		}
	}

	return
}

func (residentDataAttribute *ResidentDataAttribute) Parse(attribute Attribute) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovery %s, hex dump: %s", fmt.Sprint(r), hex.EncodeToString(attribute.AttributeBytes))
		}
	}()

	if len(attribute.AttributeBytes) < 0x18 {
		return
	}

	const offsetResidentData = 0x18
	residentDataAttribute.ResidentData = attribute.AttributeBytes[offsetResidentData:]

	return
}

func (nonResidentDataAttributes *NonResidentDataAttribute) Parse(attribute Attribute, bytesPerCluster int64) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovery %s, hex dump: %s", fmt.Sprint(r), hex.EncodeToString(attribute.AttributeBytes))
		}
	}()

	attributeLength := len(attribute.AttributeBytes)
	if attributeLength <= 0x20 {
		return
	}

	// Identify offset of the data runs in the data Attribute
	const offsetDataRunOffset = 0x20
	dataRunOffset := attribute.AttributeBytes[offsetDataRunOffset]

	if attributeLength < int(dataRunOffset) {
		return
	}

	// Pull out the data run bytes
	dataRunsBytes := make([]byte, attributeLength)
	copy(dataRunsBytes, attribute.AttributeBytes[dataRunOffset:])

	// Send the bytes to be parsed
	nonResidentDataAttributes.DataRuns = make(map[int]DataRun)
	nonResidentDataAttributes.DataRuns.Parse(dataRunsBytes, bytesPerCluster)
	if nonResidentDataAttributes.DataRuns == nil {
		err = fmt.Errorf("failed to identify data runs: %w", err)
		return
	}

	return
}

func (dataRuns *DataRuns) Parse(dataRunBytes []byte, bytesPerCluster int64) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	/*
		This function will parse the data runs from an MFT record.
		See the following for a good write up on data runs: https://homepage.cs.uri.edu/~thenry/csc487/video/66_NTFS_Data_Runs.pdf
	*/
	// Initialize a few variables
	UnparsedDataRun := UnparsedDataRun{}
	UnparsedDataRuns := make(UnparsedDataRuns)
	offset := 0
	runCounter := 0

	for {
		// Checks to see if we reached the end of the data runs. If so, break out of the loop.
		if dataRunBytes[offset] == 0x00 {
			break
		} else {
			// Take the first byte of a data run and send it to get split so we know how many bytes account for the
			// data run's offset and how many account for the data run's length.
			byteToBeSplit := dataRunBytes[offset]
			dataRunSplit := DataRunSplit{}
			dataRunSplit.Parse(byteToBeSplit)
			if dataRunSplit.offsetByteCount == 0 && dataRunSplit.lengthByteCount == 0 {
				*dataRuns = nil
				return
			}
			offset += 1

			// Pull out the the bytes that account for the data runs offset2 and length
			var lengthBytes, offsetBytes []byte

			lengthBytes = make([]byte, len(dataRunBytes[offset:(offset+dataRunSplit.lengthByteCount)]))
			copy(lengthBytes, dataRunBytes[offset:(offset+dataRunSplit.lengthByteCount)])
			offsetBytes = make([]byte, len(dataRunBytes[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)]))
			copy(offsetBytes, dataRunBytes[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)])

			// Convert the bytes for the data run offset and length to little endian int64
			UnparsedDataRun.ClusterOffset = bin.LittleEndianBinaryToInt64(offsetBytes)
			if UnparsedDataRun.ClusterOffset == 0 {
				*dataRuns = nil
				return
			}

			UnparsedDataRun.NumberOfClusters = bin.LittleEndianBinaryToInt64(lengthBytes)
			if UnparsedDataRun.NumberOfClusters == 0 {
				*dataRuns = nil
				return
			}
			// Append the data run to our data run struct
			UnparsedDataRuns[runCounter] = UnparsedDataRun

			// Increment the number order in preparation for the next data run.
			runCounter += 1

			// Set the offset2 tracker to the position of the next data run
			offset = offset + dataRunSplit.lengthByteCount + dataRunSplit.offsetByteCount
			if len(dataRunBytes) < offset {
				break
			}
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
	if len(hexToSplit) != 2 {
		dataRunSplit.offsetByteCount = 0
		dataRunSplit.lengthByteCount = 0
		return
	}

	// Split the hex string in half and return each half as an int
	var err error
	dataRunSplit.offsetByteCount, err = strconv.Atoi(string(hexToSplit[0]))
	if err != nil {
		dataRunSplit.offsetByteCount = 0
		dataRunSplit.lengthByteCount = 0
		return
	}
	dataRunSplit.lengthByteCount, err = strconv.Atoi(string(hexToSplit[1]))
	if err != nil {
		dataRunSplit.offsetByteCount = 0
		dataRunSplit.lengthByteCount = 0
		return
	}
	return
}
