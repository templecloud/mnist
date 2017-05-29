package ann

import (
	"fmt"
	// "os"
	"testing"
)

//--------------------------------------------------------------------------------------------------

func TestNewNetwork(t *testing.T) {
	layer_spec := []int{1, 1}
	nw, status := New(layer_spec, 1)
	fmt.Printf("nw: %v\n", nw)
	if status != "ok" {
		t.Error("Test failed")
	}
	input := []float64{0.00}
	output := nw.FeedForward(input)
	fmt.Printf("output: %v\n", output)
}

func TestNewNetwork1(t *testing.T) {
	layer_spec := []int{2, 3, 2}
	nw, status := New(layer_spec, 1)
	fmt.Printf("nw: %v\n", nw)
	if status != "ok" {
		t.Error("Test failed")
	}
	input := []float64{0.00, 0.00}
	output := nw.FeedForward(input)
	fmt.Printf("output: %v\n", output)
}

func TestNewNetwork2(t *testing.T) {
	layer_spec := []int{1000, 9000, 8000, 7000, 6000, 5000, 4000, 3000, 2000, 1000}
	_, status := New(layer_spec, 1)
	if status != "ok" {
		t.Error("Test failed")
	}
}

//func TestNewNetwork3(t *testing.T){
//    idxFile, _ := os.Open("/Users/Temple/Work/data/mnist/t10k-images-idx3-ubyte")
//    img, _ := sensor.ReadIdxImage(idxFile, 0)
//    layer_spec := []int{784, 30, 10}
//    nw, status := New(layer_spec, 1)
//    fmt.Printf("nw: %v\n", nw)
//    if status != "ok" {
//        t.Error("Test failed")
//    }
//}
