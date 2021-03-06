package idxfile

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

// Set-up and executes tests.
func TestMain(t *testing.M) {
	/// call flag.Parse() here if TestMain uses flags
	// Clean 'test results' directory
	os.RemoveAll(TestOutput)
	os.MkdirAll(TestOutput, os.FileMode(0777))

	os.Exit(t.Run())
}

// Assertions test function.
func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	} else {
		message = fmt.Sprintf("%v != %v\n=> ", a, b) + message
		t.Fatal(message)
	}
}

const (
	MNISTPath   = "../data/"
	IDXFileName = "t10k-images-idx3-ubyte"
	IDXFileExt  = IDX
	IDXFile     = IDXFileName + IDXFileExt
	MNIST       = MNISTPath + IDXFile
	TestOutput  = "../test-results/"
	ImageFile   = TestOutput + IDXFileName
)

func TestReadReadIdxImage(t *testing.T) {
	idxFile, _ := os.Open(MNIST)
	ReadIdxImage(idxFile, 5000)
	fmt.Printf("image created\n")
}

func TestReadReadIdxImages(t *testing.T) {
	idxFile, _ := os.Open(MNIST)
	imageIdxs := []int{0, 999, 1999, 2999, 3999, 4999, 5999, 6999, 7999, 8999, 9999}
	images, _ := ReadIdxImages(idxFile, imageIdxs)
	fmt.Printf("num images created: %v\n", len(images))
	assertEqual(t, len(imageIdxs), len(images), "The correct number of images are extracted.")
}

func TestReadAllIdxImages(t *testing.T) {
	idxFile, _ := os.Open(MNIST)
	header, _ := ReadIdxHeader(idxFile)
	images, _ := ReadAllIdxImages(idxFile)
	fmt.Printf("num images created: %v\n", len(images))
	assertEqual(t, int(header.dimensions[0]), len(images), "All images are extracted.")
}

func TestReadIdxWritePngImage(t *testing.T) {

	idxFile, _ := os.Open(MNIST)
	header, _ := ReadIdxHeader(idxFile)

	imageIdx := 0
	byteImage, _ := ReadIdxImage(idxFile, imageIdx)
	pngFile, _ := os.Create(ImageFile + "-" + strconv.Itoa(imageIdx) + ".png")
	outputImage, _ := WritePngImage(byteImage, pngFile)
	fmt.Printf("outputImage: %v\n", outputImage)

	imageIdx1 := 5000
	byteImage1, _ := ReadIdxImage(idxFile, imageIdx1)
	pngFile1, _ := os.Create(ImageFile + "-" + strconv.Itoa(imageIdx1) + ".png")
	outputImage1, _ := WritePngImage(byteImage1, pngFile1)
	fmt.Printf("outputImage: %v\n", outputImage1)

	imageIdx2 := int(header.dimensions[0] - 1)
	byteImage2, _ := ReadIdxImage(idxFile, imageIdx2)
	pngFile2, _ := os.Create(ImageFile + "-" + strconv.Itoa(imageIdx2) + ".png")
	outputImage2, _ := WritePngImage(byteImage2, pngFile2)
	fmt.Printf("outputImage: %v\n", outputImage2)
}

func TestReadIdxWritePngImage2(t *testing.T) {

	idxFile, _ := os.Open(MNIST)

	imageIdxs := []int{0, 1, 2, 3, 4, 5}
	byteImages, _ := ReadIdxImages(idxFile, imageIdxs)
	for i, byteImage := range byteImages {
		pngFile, _ := os.Create(ImageFile + "-" + strconv.Itoa(imageIdxs[i]) + ".png")
		outputImage, _ := WritePngImage(byteImage, pngFile)
		fmt.Printf("outputImage: %v\n", outputImage)
	}
}

func TestReadIdxWritePngImage3(t *testing.T) {

	idxFile, _ := os.Open(MNIST)

	imageIdxs := []int{5555, 6666, 7777}
	byteImages, _ := ReadIdxImages(idxFile, imageIdxs)
	for i, byteImage := range byteImages {
		pngFile, _ := os.Create(ImageFile + "-" + strconv.Itoa(imageIdxs[i]) + ".png")
		outputImage, _ := WritePngImage(byteImage, pngFile)
		fmt.Printf("outputImage: %v\n", outputImage)
	}
}
