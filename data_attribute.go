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

type RawDataAttribute []byte
type RawResidentDataAttribute []byte
type RawNonResidentDataAttribute []byte
type RawDataRuns []byte
type RawDataRunSplitByte byte
type RawDataRun []byte
type ResidentDataAttribute []byte

type NonResidentDataAttribute struct {
	DataRuns DataRuns
}

type UnresolvedDataRun struct {
	NumberOrder      int
	ClusterOffset    int64
	NumberOfClusters int64
}

type UnresolvedDataRuns map[int]UnresolvedDataRun

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
	TotalSize                uint8
	FlagResident             bool
	ResidentDataAttribute    ResidentDataAttribute
	NonResidentDataAttribute NonResidentDataAttribute
}

func (rawDataAttribute RawDataAttribute) Parse(bytesPerCluster int64) (nonResidentDataAttribute NonResidentDataAttribute, residentDataAttribute ResidentDataAttribute, err error) {
	sizeOfRawDataAttribute := len(rawDataAttribute)
	if sizeOfRawDataAttribute == 0 {
		err = errors.New("received nil bytes")
		return
	}

	if bytesPerCluster == 0 {
		err = errors.New("did not receive a value for bytes per cluster")
		return
	}

	//TODO: handle resident data
	const offsetResidentFlag = 0x08
	if rawDataAttribute[offsetResidentFlag] == 0x00 {
		rawResidentDataAttribute := RawResidentDataAttribute(make([]byte, sizeOfRawDataAttribute))
		copy(rawResidentDataAttribute, rawDataAttribute)
		residentDataAttribute, err = rawResidentDataAttribute.Parse()
		if err != nil {
			err = fmt.Errorf("failed to parse resident data attribute: %v", err)
			return
		}
		return
	} else {
		rawNonResidentDataAttribute := RawNonResidentDataAttribute(make([]byte, sizeOfRawDataAttribute))
		copy(rawNonResidentDataAttribute, rawDataAttribute)
		nonResidentDataAttribute, err = rawNonResidentDataAttribute.Parse(bytesPerCluster)
		if err != nil {
			err = fmt.Errorf("failed to parse non resident data attribute: %v", err)
			return
		}
	}

	return
}

func (rawResidentDataAttribute RawResidentDataAttribute) Parse() (residentDataAttribute ResidentDataAttribute, err error) {
	const offsetResidentData = 0x18
	sizeOfRawResidentDataAttribute := len(rawResidentDataAttribute)
	if sizeOfRawResidentDataAttribute == 0 {
		err = errors.New("received nil bytes")
		return
	} else if sizeOfRawResidentDataAttribute < offsetResidentData {
		err = fmt.Errorf("expected to receive at least 18 bytes, but received %v", sizeOfRawResidentDataAttribute)
		return
	}

	copy(residentDataAttribute, rawResidentDataAttribute[offsetResidentData:])

	return
}

func (rawNonResidentDataAttribute RawNonResidentDataAttribute) Parse(bytesPerCluster int64) (nonResidentDataAttributes NonResidentDataAttribute, err error) {
	const offsetDataRunOffset = 0x20
	sizeOfRawNonResidentDataAttribute := len(rawNonResidentDataAttribute)
	if sizeOfRawNonResidentDataAttribute == 0 {
		err = errors.New("received nil bytes")
		return
	} else if sizeOfRawNonResidentDataAttribute <= offsetDataRunOffset {
		err = fmt.Errorf("expected to receive at least 18 bytes, but received %v", sizeOfRawNonResidentDataAttribute)
		return
	}

	// Identify offset of the data runs in the data Attribute

	dataRunOffset := rawNonResidentDataAttribute[offsetDataRunOffset]

	if sizeOfRawNonResidentDataAttribute < int(dataRunOffset) {
		err = errors.New("data run offset is beyond the size of the byte slice")
		return
	}

	// Pull out the data run bytes
	rawDataRuns := RawDataRuns(make([]byte, sizeOfRawNonResidentDataAttribute))
	copy(rawDataRuns, rawNonResidentDataAttribute[dataRunOffset:])

	// Send the bytes to be parsed
	nonResidentDataAttributes.DataRuns, err = rawDataRuns.Parse(bytesPerCluster)
	if err != nil {
		err = fmt.Errorf("filed to parse data runs: %v", err)
		return
	}

	return
}

func (rawDataRuns RawDataRuns) Parse(bytesPerCluster int64) (dataRuns DataRuns, err error) {
	if rawDataRuns == nil {
		err = errors.New("received null bytes")
		return
	}

	// Initialize a few variables
	UnresolvedDataRun := UnresolvedDataRun{}
	UnresolvedDataRuns := make(UnresolvedDataRuns)
	sizeOfRawDataRuns := len(rawDataRuns)
	dataRuns = make(DataRuns)
	offset := 0
	runCounter := 0

	for {
		// Checks to see if we reached the end of the data runs. If so, break out of the loop.
		if rawDataRuns[offset] == 0x00 || sizeOfRawDataRuns < offset {
			break
		} else {
			// Take the first byte of a data run and send it to get split so we know how many bytes account for the
			// data run's offset and how many account for the data run's length.
			byteToBeSplit := RawDataRunSplitByte(rawDataRuns[offset])
			dataRunSplit := byteToBeSplit.Parse()
			offset += 1

			// Pull out the the bytes that account for the data runs offset2 and length
			var lengthBytes, offsetBytes []byte

			lengthBytes = make([]byte, len(rawDataRuns[offset:(offset+dataRunSplit.lengthByteCount)]))
			copy(lengthBytes, rawDataRuns[offset:(offset+dataRunSplit.lengthByteCount)])
			offsetBytes = make([]byte, len(rawDataRuns[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)]))
			copy(offsetBytes, rawDataRuns[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)])

			// Convert the bytes for the data run offset and length to little endian int64
			UnresolvedDataRun.ClusterOffset, _ = bin.LittleEndianBinaryToInt64(offsetBytes)
			UnresolvedDataRun.NumberOfClusters, _ = bin.LittleEndianBinaryToInt64(lengthBytes)

			// Append the data run to our data run struct
			UnresolvedDataRuns[runCounter] = UnresolvedDataRun

			// Increment the number order in preparation for the next data run.
			runCounter += 1

			// Set the offset tracker to the position of the next data run
			offset = offset + dataRunSplit.lengthByteCount + dataRunSplit.offsetByteCount
		}
	}

	// resolve Data Runs
	dataRunOffset := int64(0)
	for i := 0; i < len(UnresolvedDataRuns); i++ {
		dataRunOffset = dataRunOffset + (UnresolvedDataRuns[i].ClusterOffset * bytesPerCluster)
		dataRuns[i] = DataRun{
			AbsoluteOffset: dataRunOffset,
			Length:         UnresolvedDataRuns[i].NumberOfClusters * bytesPerCluster,
		}
	}
	return
}

/*
	This function will split the first byte of a data run.
	See the following for a good write up on data runs: https://homepage.cs.uri.edu/~thenry/csc487/video/66_NTFS_Data_Runs.pdf
*/
func (rawDataRunSplitByte RawDataRunSplitByte) Parse() (dataRunSplit DataRunSplit) {
	// Convert the byte to a hex string
	hexToSplit := fmt.Sprintf("%x", rawDataRunSplitByte)

	if len(hexToSplit) == 1 {
		dataRunSplit.offsetByteCount = 0
		dataRunSplit.lengthByteCount, _ = strconv.Atoi(string(hexToSplit[0]))
	} else {
		// Split the hex string in half and return each half as an int
		dataRunSplit.offsetByteCount, _ = strconv.Atoi(string(hexToSplit[0]))
		dataRunSplit.lengthByteCount, _ = strconv.Atoi(string(hexToSplit[1]))
	}
	return
}
