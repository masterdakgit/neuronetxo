package main

import (
	"fmt"
	"neuron/nr"
)

func main() {
	nr.CreateLayer([]int{2, 1})
	nr.NCorrect = 0.1
	nr.SetAnswers([]float64{0.5})

	for x := 0; x < 100; x++ {
		nr.Calc()
		fmt.Println(nr.Layers[1][0].Err)
	}

}
