package neat

import (
	"fmt"
	"math/rand"
	"testing"
)

func NEATUnitTest() {
	fmt.Println("===== NEAT Unit Test =====")

	fmt.Println("\x1b[32m=Testing config JSON file import...\x1b[0m")
	config, err := NewConfigJSON("config.json")
	if err != nil {
		fmt.Println("\x1b[31mFAIL\x1b[0m")
	}
	config.Summarize()

	fmt.Println("\x1b[32m=Testing creating and running NEAT...\x1b[0m")
	New(config, func(n *NeuralNetwork) float64 {
		return rand.Float64()
	}).Run()
}

func TestNEAT(t *testing.T) {
	rand.Seed(0)
	NEATUnitTest()
}