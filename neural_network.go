package neat

import (
	"fmt"
	"sort"
)

// Neuron is an implementation of a single neuron of a neural network.
type Neuron struct {
	ID         int                 // neuron ID
	Type       string              // neuron type
	Activated  bool                // true if it has been activated
	Signal     float64             // signal held by this neuron
	Synapses   map[*Neuron]float64 // synapse from input neurons
	Activation *ActivationFunc     // activation function
}

// NewNeuron returns a new instance of neuron, given a node gene.
func NewNeuron(nodeGene *NodeGene) *Neuron {
	return &Neuron{
		ID:         nodeGene.ID,
		Type:       nodeGene.Type,
		Activated:  false,
		Signal:     0.0,
		Synapses:   make(map[*Neuron]float64),
		Activation: nodeGene.Activation,
	}
}

// String returns the string representation of Neuron.
func (n *Neuron) String() string {
	if len(n.Synapses) == 0 {
		return fmt.Sprintf("[%s(%d, %s)]", n.Type, n.ID, n.Activation.Name)
	}
	str := fmt.Sprintf("[%s(%d, %s)] (\n", n.Type, n.ID, n.Activation.Name)
	for neuron, weight := range n.Synapses {
		str += fmt.Sprintf("  <--{%.3f}--[%s(%d, %s)]\n",
			weight, neuron.Type, neuron.ID, neuron.Activation.Name)
	}
	return str + ")"
}

// Activate retrieves signal from neurons that are connected to this neuron and
// return its signal.
func (n *Neuron) Activate() float64 {
	// if the neuron's already activated, or it isn't connected from any neurons,
	// return its current signal.
	if n.Activated || len(n.Synapses) == 0 {
		return n.Signal
	}
	n.Activated = true

	inputSum := 0.0
	for neuron, weight := range n.Synapses {
		inputSum += neuron.Activate() * weight
	}
	n.Signal = n.Activation.Fn(inputSum)
	return n.Signal
}

// NeuralNetwork is an implementation of the phenotype neural network that is
// decoded from a genome.
type NeuralNetwork struct {
	NumInputs  int       // number of inputs
	NumOutputs int       // number of outputs
	Neurons    []*Neuron // neurons in the neural network
}

// NewNeuralNetwork returns a new instance of NeuralNetwork given a genome to
// decode from.
func NewNeuralNetwork(g *Genome) *NeuralNetwork {
	sort.Slice(g.NodeGenes, func(i, j int) bool {
		return g.NodeGenes[i].ID < g.NodeGenes[j].ID
	})

	numInputs := 0
	numOutputs := 0

	neurons := make([]*Neuron, 0, len(g.NodeGenes))
	for _, nodeGene := range g.NodeGenes {
		if nodeGene.Type == "input" {
			numInputs++
		} else if nodeGene.Type == "output" {
			numOutputs++
		}
		neurons = append(neurons, NewNeuron(nodeGene))
	}

	for _, connGene := range g.ConnGenes {
		if !connGene.Disabled {
			if in := sort.Search(len(neurons), func(i int) bool {
				return neurons[i].ID >= connGene.From.ID
			}); in < len(neurons) && neurons[in].ID == connGene.From.ID {
				if out := sort.Search(len(neurons), func(i int) bool {
					return neurons[i].ID >= connGene.To.ID
				}); out < len(neurons) && neurons[out].ID == connGene.To.ID {
					neurons[out].Synapses[neurons[in]] = connGene.Weight
				}
			}
		}
	}
	return &NeuralNetwork{numInputs, numOutputs, neurons}
}

// String returns the string representation of NeuralNetwork.
func (n *NeuralNetwork) String() string {
	str := fmt.Sprintf("NeuralNetwork(%d, %d):\n", n.NumInputs, n.NumOutputs)
	for _, neuron := range n.Neurons {
		str += neuron.String() + "\n"
	}
	return str[:len(str)-1]
}

// Feedforward propagates inputs signals from input neurons to output neurons,
// and return output signals.
func (n *NeuralNetwork) FeedForward(inputs []float64) ([]float64, error) {
	if len(inputs) != n.NumInputs {
		errStr := "Invalid number of inputs: %d != %d"
		return nil, fmt.Errorf(errStr, n.NumInputs, len(inputs))
	}

	// register sensor inputs
	for i := 0; i < n.NumInputs; i++ {
		n.Neurons[i].Signal = inputs[i]
	}

	// recursively propagate from input neurons to output neurons
	outputs := make([]float64, 0, n.NumOutputs)
	for i := n.NumInputs; i < n.NumInputs+n.NumOutputs; i++ {
		outputs = append(outputs, n.Neurons[i].Activate())
	}

	// reset all neurons
	for _, neuron := range n.Neurons {
		neuron.Activated = false
		neuron.Signal = 0.0
	}

	return outputs, nil
}