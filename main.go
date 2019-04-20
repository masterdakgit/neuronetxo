package main

import (
	"fmt"
	"math/rand"
	"neuron/nr"
)

var (
	XO, XO0 [9]float64
	XA      []float64
)

func main() {
	nr.CreateLayer([]int{9, 25, 9})
	nr.NCorrect = 0.1
	for x := 0; x < 10000; x++ {
		XOPrepare()
		nr.Calc()
		err := float64(0)
		for n := 0; n < 9; n++ {
			err += nr.Layers[2][n].Err * nr.Layers[2][n].Err
		}
	}
	XOPrint()
	NeuroAnswer()
}

func XOPrepare() {
	XO = XO0
	for x := 0; x < 8; x++ {
		r := rand.Intn(9)
		XO[r] = 1
	}

	for n := 0; n < 9; n++ {
		nr.Layers[0][n].Out = XO[n]
	}

	XA = make([]float64, 9)

	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			XA[n] = 1
		}
	}

	nr.SetAnswers(XA)
}

func XOPrint() {
	for n := 0; n < 9; n++ {
		if n%3 == 0 {
			fmt.Println()
		}
		fmt.Print(XO[n])
	}
}

func NeuroAnswer() {
	for n := 0; n < 9; n++ {
		if n%3 == 0 {
			fmt.Println()
		}
		fmt.Printf("%.2f", nr.Layers[2][n].Out)
	}

}
