package idxfile

import (
	"image"
	"image/png"
	"os"
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
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	return pngFile.Name(), err
}
