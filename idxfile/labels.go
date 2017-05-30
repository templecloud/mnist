package idxfile

import (
	"os"
)

func ReadIdxLabel(idxFile *os.File, imageIdx int) (int, error) {
	// extract header
	idxHeader, _ := ReadIdxHeader(idxFile)
	// compute images values - downcast to int for simplicity...
	fileHeaderOffset := int(idxHeader.fileHeaderOffset())
	dataWidth := int(idxHeader.idxMagic.dataType.numBytes)
	// calculate offsets (as int64 to obey first seek result)
	currentOffset, _ := idxFile.Seek(0, 1)
	startOffset := int64(fileHeaderOffset + imageIdx)
	endOffset := startOffset + int64(dataWidth)
	// extract images bytes
	idxFile.Seek(startOffset-currentOffset, 1)
	labelBytes := make([]byte, endOffset-startOffset)
	_, err := idxFile.Read(labelBytes)
	// build result
	// TODO: either simplify or make generic

	label := int(labelBytes[0])
	// return the label
	return label, err
}
