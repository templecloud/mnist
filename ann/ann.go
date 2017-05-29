package ann

import (
	"log"
	"math"
	"math/rand"
)

//--------------------------------------------------------------------------------------------------

type Network struct {
	num_layers int           // num layers
	layer_spec []int         // the number of neurons in each layer
	weights    [][][]float64 // weight matrices for each neuron
	biases     [][]float64   // bias matrices for each neuron
}

type Layer interface {
	Layer(idx int) []int
}

type Neuron interface {
	Layer(idx int) []int
}

// Generate a randomised neural network for the specification.
func New(layer_spec []int, seed int64) (Network, string) {
	r := rand.New(rand.NewSource(seed))
	// create the network and allocate space for the layers
	num_layers := len(layer_spec)
	weights := make([][][]float64, num_layers)
	biases := make([][]float64, num_layers) // bias matrices for each neuron
	// for each layer in the n/w...
	for l := 1; l < num_layers; l++ {
		// make a layer of neurons
		num_neurons := layer_spec[l]
		num_neuron_inputs := layer_spec[l-1]
		neuron_weights := make([][]float64, num_neurons)
		neuron_biases := make([]float64, num_neurons)
		// for each neuron in the layer...
		for n := 0; n < num_neurons; n++ {
			// make a list of weights for each incoming input
			input_weights := make([]float64, num_neuron_inputs)
			// for each incoming input from the previous layer...
			for w := 0; w < num_neuron_inputs; w++ {
				// set a random weight between 0..1
				input_weights[w] = r.Float64()
			}
			// set the input weights
			neuron_weights[n] = input_weights
			// add a random bias between 0..1
			neuron_biases[n] = r.Float64()
		}
		// set the weights and bias for the neurons in the layer
		weights[l] = neuron_weights
		biases[l] = neuron_biases
	}
	return Network{num_layers, layer_spec, weights, biases}, "ok"
}

// Calculate the output of the specified network from the input
func (nw *Network) FeedForward(nw_input []float64) (nw_output []float64) {
	// allocate space for the result
	var layer_output []float64
	// process each layer and propgate the results through the network
	num_layers := len(nw.layer_spec)
	layer_input := nw_input
	for l := 1; l < num_layers; l++ {
		num_neurons := nw.layer_spec[l]
		layer_output = make([]float64, num_neurons)
		log.Printf("layer: %v - num_neurons: %v\n", l, num_neurons)
		log.Printf("layer: %v - input      : %v\n", l, layer_input)
		// for each neuron in the layer, calculate the output from the inputs in the previous layer
		for n := 0; n < num_neurons; n++ {
			// dereference the weights and biases for this layer from the network
			neuron_weights := nw.weights[l][n]
			neuron_bias := nw.biases[l][n]
			log.Printf("neuron[%v-%v] - input  : %v\n", l, n, layer_input)
			log.Printf("neuron[%v-%v] - weights: %v\n", l, n, neuron_weights)
			log.Printf("neuron[%v-%v] - bias   : %v\n", l, n, neuron_bias)
			// calculate the ouput of the neuron from its inputs in the previous layer
			num_neuron_inputs := len(neuron_weights)
			neuron_output := 0.00
			for w := 0; w < num_neuron_inputs; w++ {
				neuron_output += neuron_weights[w] * layer_input[w]
			}
			neuron_output += neuron_bias
			neuron_output = 1.0 / (1.0 + math.Exp(-neuron_output))
			log.Printf("neuron[%v-%v] - output : %v\n", l, n, neuron_output)
			// add the result to the output of the layer
			layer_output[n] = neuron_output
		}
		// set the newly calculated output as the input to the next iteration
		log.Printf("layer: %v - ouput      : %v\n", l, layer_output)
		layer_input = layer_output
	}
	// Return the final layer ouput
	return layer_output
}

// Calculate the sigmoid output of a vlaue between 0..1
func Sigmoid(k float64) float64 {
	return 1.0 / (1.0 + math.Exp(-k))
}

func StochasitcGradientDescent() {

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
