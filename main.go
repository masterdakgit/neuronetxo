package main

import (
	"neuron/nr"
)

func main() {
	nr.CreateLayer([]int{9, 8, 4, 3})
	nr.NCorrect = 0.1

	for x := 0; x < 10000; x++ {
	}
}
