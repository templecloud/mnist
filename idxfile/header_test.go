//--------------------------------------------------------------------------------------------------
package idxfile

import (
	"log"
	"os"
	"testing"
)

func TestExtractIdxHeader(t *testing.T) {
	idxFile, _ := os.Open(MNIST)
	output, _ := ReadIdxHeader(idxFile)
	log.Printf("output: %v\n", output)
}

func TestConvertType(t *testing.T) {
	var res IdxDataType
	var err error

	res, err = makeIdxDataType(byte(0x08))
	if res.name != "unsigned" && err == nil {
		t.Error("Expected 'unsigned'")
	}

	res, err = makeIdxDataType(byte(0x09))
	if res.name != "signed" && err != nil {
		t.Error("Expected 'signed'")
	}

	res, err = makeIdxDataType(byte(0x0B))
	if res.name != "short" && err != nil {
		t.Error("Expected 'short'")
	}

	res, err = makeIdxDataType(byte(0x0C))
	if res.name != "int" && err != nil {
		t.Error("Expected 'int'")
	}

	res, err = makeIdxDataType(byte(0x0D))
	if res.name != "float" && err != nil {
		t.Error("Expected 'float'")
	}

	res, err = makeIdxDataType(byte(0x0E))
	if res.name != "double" && err != nil {
		t.Error("Expected 'double'")
	}

	res, err = makeIdxDataType(byte(0x01))
	if err.Error() != "Unknown IdxDataType for byte: 1" {
		t.Error("Expected an error for byte 0x01")
	}
}
