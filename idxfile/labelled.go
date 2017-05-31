package idxfile

import (
	"fmt"
	"os"
)

type LabelledImage struct {
	img   ByteImage
	label int
}

func ReadLabelledImage(imageFile *os.File, labelFile *os.File, idx int) (LabelledImage, error) {
	var labelledImage LabelledImage
	// read image
	image, err := ReadIdxImage(imageFile, idx)
	if err != nil {
		return labelledImage, err
	}
	// read label
	label, err := ReadIdxLabel(labelFile, idx)
	if err != nil {
		return labelledImage, err
	}
	// return result
	labelledImage = LabelledImage{image, label}
	return labelledImage, err
}

func ReadLabelledImages(imageFile *os.File, labelFile *os.File, idxs []int) ([]LabelledImage, error) {
	var err error
	numLabelledImages := len(idxs)
	labelledImages := make([]LabelledImage, numLabelledImages)
	for i := 0; i < numLabelledImages; i++ {
		labelledImages[i], err = ReadLabelledImage(imageFile, labelFile, i)
		if err != nil {
			return labelledImages, err
		}
	}
	return labelledImages, err
}

// Read all LabelledImages from IDX format data files.
func ReadAllLabelledImages(imageFile *os.File, labelFile *os.File) ([]LabelledImage, error) {
	var labelledImages []LabelledImage
	var err error
	// extract header and numImages -  - downcast to int for simplicity...
	imageIdxHeader, _ := ReadIdxHeader(imageFile)
	numImages := int(imageIdxHeader.dimensions[0])
	labelIdxHeader, _ := ReadIdxHeader(labelFile)
	numLabels := int(labelIdxHeader.dimensions[0])
	if numImages != numLabels {
		return labelledImages, fmt.Errorf("Inconsistent number of iamges and labels %i != %i", numImages, numLabels)
	}
	// extract images
	labelledImages = make([]LabelledImage, numImages)
	for i := 0; i < numImages; i++ {
		labelledImages[i], _ = ReadLabelledImage(imageFile, labelFile, i)
	}
	return labelledImages, err
}
