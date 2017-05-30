//--------------------------------------------------------------------------------------------------
//==================================================================================================
// MNIST IDX file format
//--------------------------------------------------------------------------------------------------
package idxfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

//--------------------------------------------------------------------------------------------------
// The data type of an IdxMagic number.

type IdxDataType struct {
	code     int
	name     string
	numBytes int
}

// Create an IdxDataType from the specified byte.
func makeIdxDataType(dataType byte) (IdxDataType, error) {
	var idxDataType IdxDataType
	var err error
	switch {
	case dataType == 0x08:
		idxDataType = IdxDataType{0x08, "unsigned", 1}
	case dataType == 0x09:
		idxDataType = IdxDataType{0x09, "signed", 1}
	case dataType == 0x0B:
		idxDataType = IdxDataType{0x0B, "short", 2}
	case dataType == 0x0C:
		idxDataType = IdxDataType{0x0C, "int", 4}
	case dataType == 0x0D:
		idxDataType = IdxDataType{0x0D, "float", 4}
	case dataType == 0x0E:
		idxDataType = IdxDataType{0x0E, "double", 4}
	case true:
		err = fmt.Errorf("Unknown IdxDataType for byte: %x", dataType)
	}
	return idxDataType, err
}

//--------------------------------------------------------------------------------------------------
// IDX Magic Number

// The magic number of an IdxHeader.
type IdxMagic struct {
	nilBytes      bool
	dataType      IdxDataType
	numDimensions int
}

// Create an IdxMagic from the specified magic bytes.
func makeIdxMagic(magicBytes []byte) IdxMagic {
	nilBytes, _ := makeNilBytes(magicBytes[0:2])
	idxDataType, _ := makeIdxDataType(magicBytes[2:3][0])
	numDimensions, _ := makeNumDimensions(magicBytes[3:])
	return IdxMagic{nilBytes, idxDataType, numDimensions}
}

// Return true if the specified magic bytes have the required number of nil bytes.
func makeNilBytes(nilBytes []byte) (bool, error) {
	var result bool
	var err error
	value, _ := binary.Uvarint(nilBytes)
	if len(nilBytes) == 2 && value == 0 {
		result = true
	} else {
		err = fmt.Errorf("Bad nilByte in header: %x", nilBytes)
	}
	return result, err
}

// Return the number of dimensions specified in the dimension byte of the magic number header.
func makeNumDimensions(numDimByte []byte) (int, error) {
	// return int(numDimByte[0]), nil
	var numDimensions uint8
	err := binary.Read(bytes.NewReader(numDimByte), binary.BigEndian, &numDimensions)
	return int(numDimensions), err
}

//--------------------------------------------------------------------------------------------------
// IDX Header

// The header structure of an IDX format file
type IdxHeader struct {
	idxMagic   IdxMagic
	dimensions []uint32
}

func makeDimension(dimensionBytes []byte) (uint32, error) {
	var dimension uint32
	err := binary.Read(bytes.NewReader(dimensionBytes), binary.BigEndian, &dimension)
	return dimension, err
}

func (idx *IdxHeader) fileHeaderOffset() int {
	return 4 + (idx.idxMagic.numDimensions * 4)
}

// Read the IDXHeader from an IDX data format file
func ReadIdxHeader(idxFile *os.File) (IdxHeader, []error) {
	// collect all errors
	errors := make([]error, 0)
	// reset to file start
	idxFile.Seek(0, 0)
	// read magic bytes
	magicBytes, err := readBytes(idxFile, 4)
	idxMagic := makeIdxMagic(magicBytes)
	errors = appendError(errors, err)
	// read dimensional meta data
	dimensions := make([]uint32, idxMagic.numDimensions)
	for n := 0; n < idxMagic.numDimensions; n++ {
		dimNBytes, err := readBytes(idxFile, 4)
		errors = appendError(errors, err)
		dimensions[n], err = makeDimension(dimNBytes)
		errors = appendError(errors, err)
	}
	// return the header
	return IdxHeader{idxMagic, dimensions}, errors
}

func readBytes(file *os.File, numBytes int) ([]byte, error) {
	bytes := make([]byte, numBytes)
	_, err := file.Read(bytes)
	return bytes, err
}

func appendError(errors []error, error error) []error {
	return append(errors, error)
}
