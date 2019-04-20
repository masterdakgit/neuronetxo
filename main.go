package main

import (
	"fmt"
	"math"
	"math/rand"
	"neuron/nr"
)

const (
	nc = 0.1

	n0 = 3
	n1 = 9
	n2 = 1
)

var (
	NS0 [n0 + 1]Neuron
	NS1 [n1 + 1]Neuron
	NS2 [n2]Neuron

	W0 [n0 + 1][n1 + 1]float64
	W1 [n1 + 1][n2 + 1]float64
)

type Neuron struct {
	in, out, err float64
}

func main() {
	nr.CreateLayer([]int{5, 3, 1})

	for l := range nr.Layers {
		for n := range nr.Layers[l] {
			fmt.Print(nr.Layers[l][n].Out, " ")
		}
		fmt.Println()
	}

	for w := range nr.Weight {
		for n := range nr.Weight[w] {
			for x := range nr.Weight[w][n] {
				fmt.Printf("%.2f", nr.Weight[w][n][x])
				fmt.Print(" ")
			}
			fmt.Print("| ")
		}
		fmt.Println()
	}

	/*
		rand.Seed(time.Now().UnixNano())
		StartWeight()

		for n := 0; n < 100000; n++{
			InputPrepare()
			Calculate()
		}

		for{
			InputPrepare()
			Calculate()

			fmt.Println(Answer())
			for n := 0; n < n2; n++{
				fmt.Print(NS2[n].out, " ")
			}


			s := 0
			fmt.Scanln(&s)
		}*/
}

func InputPrepare() {
	for x := 0; x < n0; x++ {
		NS0[x].out = float64(rand.Intn(2))
	}
}

func StartWeight() {
	for x0 := 0; x0 < n0+1; x0++ {
		for x1 := 0; x1 < n1; x1++ {
			W0[x0][x1] = float64((rand.Intn(1000) - 500) / 10000)
		}
	}
	for x0 := 0; x0 < n1+1; x0++ {
		for x1 := 0; x1 < n2; x1++ {
			W1[x0][x1] = float64((rand.Intn(1000) - 500) / 10000)
		}
	}

	NS0[n0].out = 1
	NS1[n1].out = 1

}

func Calculate() {
	//Подсчет входов слой NS1
	for x1 := 0; x1 < n1; x1++ {
		NS1[x1].in = 0
		for x0 := 0; x0 < n0+1; x0++ {
			NS1[x1].in += NS0[x0].out * W0[x0][x1]
		}
		NS1[x1].out = Activate(NS1[x1].in)
	}

	//Подсчет входов слой NS2 и ошибки
	for x2 := 0; x2 < n2; x2++ {
		NS2[x2].in = 0
		for x1 := 0; x1 < n1+1; x1++ {
			NS2[x2].in += NS1[x1].out * W1[x1][x2]
		}
		NS2[x2].out = Activate(NS2[x2].in)
		NS2[x2].err = Answer() - NS2[x2].out
	}

	//Подсчет ошибки для слоя NS1
	for x1 := 0; x1 < n1+1; x1++ {
		NS1[x1].err = 0
		for x2 := 0; x2 < n2; x2++ {
			NS1[x1].err += NS2[x2].err * W1[x1][x2]
		}
	}

	Correct()
}

func Correct() {
	for x1 := 0; x1 < n1; x1++ {
		for x0 := 0; x0 < n0+1; x0++ {
			W0[x0][x1] += WeightChange(NS0[x0], NS1[x1])
		}
	}
	for x1 := 0; x1 < n1+1; x1++ {
		for x0 := 0; x0 < n2; x0++ {
			W1[x1][x0] += WeightChange(NS1[x1], NS2[x0])
		}
	}

}

func WeightChange(NLeft, NRight Neuron) float64 {
	return nc * NRight.err * NRight.out * (1 - NRight.out) * NLeft.out
}

func Activate(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func Answer() float64 {
	a := float64(0)
	for n := 0; n < n0; n++ {
		a += NS0[n].out
	}

	return a / 3
}
