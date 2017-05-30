package idxfile

import (
	"log"
	"os"
	"testing"
)

//--------------------------------------------------------------------------------------------------

func TestReadIdxLabel(t *testing.T) {
	// Check image file IDX Header
	// imgIn, imgErr := os.Open(MNIST_LOCAL + MNIST_TEST_IMAGES + IDX)
	imgIn, imgErr := os.Open("../data/" + MNIST_TEST_IMAGES + IDX)
	if imgErr != nil {
		log.Fatal(imgErr)
	}
	defer imgIn.Close()
	imgHeader, imgErrs := ReadIdxHeader(imgIn)
	if imgErrs != nil {
		// log.Printf("imgErrs: %v", imgErrs)
		// log.Fatal(imgErrs)
	}
	log.Printf("imgHeader: %+v\n", imgHeader)
	// TODO: assert

	// Check label file IDX header
	// labelIn, labelErr := os.Open(MNIST_LOCAL + MNIST_TEST_LABELS + IDX)
	labelIn, labelErr := os.Open("../data/" + MNIST_TEST_LABELS + IDX)
	if labelErr != nil {
		log.Fatal(labelErr)
	}
	defer labelIn.Close()
	labelHeader, labelErrs := ReadIdxHeader(labelIn)
	if labelErrs != nil {
		// log.Printf("imgErrs: %v", imgErrs)
		// log.Fatal(labelErrs)
	}
	log.Printf("imgErrs: %+v\n", labelHeader)
	// TODO: assert

	// Check IDX labels
	idx := 0
	label, _ := ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 7

	idx = 1
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 2

	idx = 2
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 1

	idx = 3
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 0

	idx = 4
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 4

	idx = 5
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 1

	idx = 5000
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 3

	idx = 5555
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 3 BUT IS 7!

	idx = 6666
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 7 BUT IS 2!

	idx = 7777
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 5 BUT IS 1!

	idx = 9999
	label, _ = ReadIdxLabel(labelIn, idx)
	log.Printf("idx: %v, label: %v\n", idx, label)
	// 6
}
