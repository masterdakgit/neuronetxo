package main

import (
	"fmt"
	"neuron/nr"
)

var (
	XO [3][]float64
)

func main() {
	nr.CreateLayer([]int{9, 8, 4, 3})
	nr.NCorrect = 0.1
	XOPrepare()

	for x := 0; x < 1000; x++ {
		e := float64(0)
		for n := 0; n < 3; n++ {
			for x := 0; x < 9; x++ {
				nr.Layers[0][x].Out = XO[n][x]
			}
			nr.SetAnswers(XO[n])
			nr.Calc()
			e += nr.Layers[3][n].Err * nr.Layers[3][n].Err
		}
		fmt.Println(e)
	}
}

func XOPrepare() {
	XO[0] = make([]float64, 9)
	XO[0][0] = 1
	XO[0][4] = 1
	XO[0][8] = 1

	XO[1] = make([]float64, 9)
	XO[1][3] = 1
	XO[1][4] = 1
	XO[1][5] = 1

	XO[2] = make([]float64, 9)
	XO[2][2] = 1
	XO[2][5] = 1
	XO[2][8] = 1
}
