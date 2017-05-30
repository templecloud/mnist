package ann

import (
	"log"
	"testing"
)

//--------------------------------------------------------------------------------------------------

func TestNewNetwork(t *testing.T) {
	layer_spec := []int{1, 1}
	nw, status := New(layer_spec, 1)
	log.Printf("nw: %v\n", nw)
	if status != "ok" {
		t.Error("Test failed")
	}
	input := []float64{0.00}
	output := nw.FeedForward(input)
	log.Printf("output: %v\n", output)
}

func TestNewNetwork1(t *testing.T) {
	layer_spec := []int{2, 3, 2}
	nw, status := New(layer_spec, 1)
	log.Printf("nw: %v\n", nw)
	if status != "ok" {
		t.Error("Test failed")
	}
	input := []float64{0.00, 0.00}
	output := nw.FeedForward(input)
	log.Printf("output: %v\n", output)
}

func TestNewNetwork2(t *testing.T) {
	layer_spec := []int{1000, 9000, 8000, 7000, 6000, 5000, 4000, 3000, 2000, 1000}
	_, status := New(layer_spec, 1)
	if status != "ok" {
		t.Error("Test failed")
	}
}

func TestNewNetwork3(t *testing.T) {
	layer_spec := []int{784, 30, 10}
	_, status := New(layer_spec, 1)
	if status != "ok" {
		t.Error("Test failed")
	}
}

//--------------------------------------------------------------------------------------------------

func TestDotProduct(t *testing.T) {
	v1a := []float64{0.2, 1.0, 0.2, 1.0, 0.2}
	v1b := []float64{1.0, 0.2, 1.0, 0.2, 1.0}

	dp1 := DotProduct(v1a, v1b)

	log.Printf("dotproduct: %f", dp1)
}
