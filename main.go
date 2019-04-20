package main

import (
	"fmt"
	"math/rand"
	"neuron/nr"
)

var (
	XO, XO0       [9]float64
	XA            []float64
	XO3x3, XO3x30 [3][3]float64
)

func main() {
	nr.CreateLayer([]int{9, 25, 9})
	nr.NCorrect = 0.8
	for x := 0; x < 10000; x++ {
		XOPrepare()
		XA = make([]float64, 9)
		XA[Checking()] = 1
		nr.SetAnswers(XA)
		nr.Calc()

		err := float64(0)
		for n := 0; n < 9; n++ {
			err += nr.Layers[2][n].Err * nr.Layers[2][n].Err
		}
	}

	for n := 0; n < 30; n++ {
		XOPrepare()
		XOPrint()
		XA = make([]float64, 9)
		XA[Checking()] = 1
		nr.SetAnswers(XA)
		nr.Calc()
		NeuroAnswer()
		fmt.Println("------------------------------")
	}

}

func XOPrepare() {
	XO = XO0
	XO3x3 = XO3x30
	for n := 0; n < 2; n++ {
		r := rand.Intn(9)
		XO[r] = 1

		x := r % 3
		y := r / 3
		XO3x3[x][y] = 1

	}

	for n := 0; n < 9; n++ {
		nr.Layers[0][n].Out = XO[n]
	}

}

func XOPrint() {
	for n := 0; n < 9; n++ {
		if n%3 == 0 {
			fmt.Println()
		}
		fmt.Print(XO[n], " ")
	}
	fmt.Println()
}

func NeuroAnswer() {
	for n := 0; n < 9; n++ {
		if n%3 == 0 {
			fmt.Println()
		}
		fmt.Printf("%.0f", nr.Layers[2][n].Out)
		fmt.Print(" ")
	}
	fmt.Println()

}

func Checking() (a int) {
	var xa, ya, t int

	if XO3x3[1][1] == 0 {
		xa = 1
		ya = 1
	}

	for x := 0; x < 3; x++ {
		z := 0
		for y := 0; y < 3; y++ {
			if XO3x3[x][y] > 0.5 {
				z++
			}
			if XO3x3[x][y] < 0.5 {
				t = y
			}
		}
		if z == 2 {
			ya = t
			xa = x
			return ya*3 + xa
		}
	}

	for y := 0; y < 3; y++ {
		z := 0
		for x := 0; x < 3; x++ {
			if XO3x3[x][y] > 0.5 {
				z++
			}
			if XO3x3[x][y] < 0.5 {
				t = x
			}
		}
		if z == 2 {
			ya = y
			xa = t
			return ya*3 + xa
		}
	}

	z := 0
	for xy := 0; xy < 3; xy++ {
		if XO3x3[xy][xy] > 0.5 {
			z++
		}
		if z == 2 {
			for xy2 := 0; xy2 < 3; xy2++ {
				if XO3x3[xy2][xy2] < 0.5 {
					xa = xy2
					ya = xy2
					return ya*3 + xa
				}
			}
		}
	}
	z = 0
	for x := 2; x >= 0; x-- {
		y := 2 - x
		if XO3x3[x][y] > 0.5 {
			z++
		}
		if z == 2 {
			for x := 2; x >= 0; x-- {
				y := 2 - x
				if XO3x3[x][y] < 0.5 {
					xa = x
					ya = y
					return ya*3 + xa
				}
			}
		}
	}

	return ya*3 + xa
}
