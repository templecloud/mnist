package idxfile

import (
	"os"
)

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
)

//==================================================================================================
// ByteImage

// An in-memory image stored as a consecuative array of bytes.
// NB: The IDX dictates the height and width dimensions hsould actuall be uint32 types. but, that
// is a bit pointless and a pain in the arse....
type ByteImage struct {
	width  int
	height int
	bytes  []byte
}

//  Get the image bytes as a formatted matrix.
func (img *ByteImage) pixels() [][]byte {
	byteIdx := 0
	pixels := make([][]byte, img.height)
	for rowIdx := 0; rowIdx < img.height; rowIdx++ {
		pixels[rowIdx] = make([]byte, img.width)
		for colIdx := 0; colIdx < img.width; colIdx++ {
			pixels[rowIdx][colIdx] = img.bytes[byteIdx]
			byteIdx++
		}
	}
	return pixels
}

//func (img *ByteImage) normalise() []float64 {
//
//	normalised := make([]float64, len(img.bytes))
//	for i := 0; i < len(img.bytes); i++ {
//		img.bytes[i] = float64(img.bytes[i] / 256)
//	}
//
//	return normalised
//}

//==================================================================================================
// MNIST IDX file format

//--------------------------------------------------------------------------------------------------
// IDX Data Type

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

//--------------------------------------------------------------------------------------------------
// IDX ByeImage store functions

// Read a ByteImage from an IDX format data file.
func ReadIdxImage(idxFile *os.File, imageIdx int) (ByteImage, error) {
	// extract header
	idxHeader, _ := ReadIdxHeader(idxFile)
	// compute images values - downcast to int for simplicity...
	fileHeaderOffset := int(idxHeader.fileHeaderOffset())
	dataWidth := int(idxHeader.idxMagic.dataType.numBytes)
	imgWidth := int(idxHeader.dimensions[1])
	imgHeight := int(idxHeader.dimensions[2])
	imgSize := imgWidth * imgHeight * dataWidth
	// calculate offsets (as int64 to obey first seek result)
	currentOffset, _ := idxFile.Seek(0, 1)
	startOffset := int64(fileHeaderOffset + (imgSize * imageIdx))
	endOffset := startOffset + int64(imgSize)
	// extract images bytes
	idxFile.Seek(startOffset-currentOffset, 1)
	imageBytes := make([]byte, endOffset-startOffset)
	_, err := idxFile.Read(imageBytes)
	// build result
	img := ByteImage{imgWidth, imgHeight, imageBytes}
	// return the image
	return img, err
}

// Read a set of ByteImages from an IDX format data file.
func ReadIdxImages(idxFile *os.File, imageIdxs []int) ([]ByteImage, error) {
	numImages := len(imageIdxs)
	images := make([]ByteImage, numImages)
	for i := 0; i < numImages; i++ {
		images[i], _ = ReadIdxImage(idxFile, i)
	}
	return images, nil
}

// Read all ByteImages from an IDX format data file.
func ReadAllIdxImages(idxFile *os.File) ([]ByteImage, error) {
	// extract header and numImages -  - downcast to int for simplicity...
	idxHeader, _ := ReadIdxHeader(idxFile)
	numImages := int(idxHeader.dimensions[0])
	// extract images
	images := make([]ByteImage, numImages)
	for i := 0; i < numImages; i++ {
		images[i], _ = ReadIdxImage(idxFile, i)
	}
	return images, nil
}

// Generate a PNG format image file from the speicified ByteImage.
func WritePngImage(byteImage ByteImage, pngFile *os.File) (string, error) {
	// compute images values
	imgWidth := byteImage.width
	imgHeight := byteImage.height
	// create image
	imgRect := image.Rect(0, 0, imgWidth, imgHeight)
	img := image.NewGray(imgRect)
	img.Pix = byteImage.bytes
	img.Stride = imgWidth
	// write image
	err := png.Encode(pngFile, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pngFile.Name(), err
}

//---------------------------------------------------------------------------------------------------------------------
//==================================================================================================
// MNIST IDX ByteImageStore downloader

//--------------------------------------------------------------------------------------------------
func gunzip(source string, destination string) (err error) {
	// open an input stream on the source path
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	// open a gzip reader on the source stream
	archive, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer archive.Close()
	// open an ouput stream on the destination path
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	// copy the unarchived stream to the destination file
	_, err = io.Copy(out, archive)
	if err != nil {
		return err
	}
	return nil
}

type Download struct {
	Source, Destination string
}

func download(url string, filepath string) (err error) {
	// create file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// download mnist
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// save mnsit
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	return nil
}

// http://yann.lecun.com/exdb/mnist/
const (
	MNIST_LOCAL        = "./data/"
	MNIST_REMOTE       = "http://yann.lecun.com/exdb/mnist/"
	MNIST_TRAIN_IMAGES = "train-images-idx3-ubyte"
	MNIST_TRAIN_LABELS = "train-labels-idx1-ubyte"
	MNIST_TEST_IMAGES  = "t10k-images-idx3-ubyte"
	MNIST_TEST_LABELS  = "t10k-labels-idx1-ubyte"
	GZIP               = ".gz"
	IDX                = ".idx"
)

func DownloadMNIST() (err error) {
	// ensure the specified local directory exist
	err = os.MkdirAll(MNIST_LOCAL, os.FileMode(0777))
	if err != nil {
		log.Fatal(err)
	}
	// download and unarchive each mnist file
	var mnistFiles = [...]string{
		MNIST_TRAIN_IMAGES, MNIST_TRAIN_LABELS, MNIST_TEST_IMAGES, MNIST_TEST_LABELS}
	for _, v := range mnistFiles {
		// download zipped file
		source := MNIST_REMOTE + v + GZIP
		destination := MNIST_LOCAL + v + GZIP
		unarchived := strings.TrimSuffix(destination, GZIP) + IDX
		if _, err := os.Stat(destination); os.IsNotExist(err) {
			toDownload := Download{source, destination}
			log.Printf("Downloading File: %+v\n", toDownload)
			err := download(source, destination)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Already downloaded: %s\n", destination)
		}
		// unarchive
		if _, err := os.Stat(unarchived); os.IsNotExist(err) {
			log.Printf("Local File     : %s\n", destination)
			log.Printf("Unarchived File: %s\n", unarchived)
			err := gunzip(destination, unarchived)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Already unarchived: %s\n", unarchived)
		}
	}
	return nil
}
