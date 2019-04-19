package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const nc = 0.1

var (
	W    [4]float64
	X    [4]float64
	N    Neuron
	errN [100]float64
)

type Neuron struct {
	in, out, err float64
}

func main() {
	rand.Seed(time.Now().UnixNano())

	X[0] = 0
	X[1] = 0
	X[2] = 0
	X[3] = 1

	W[0] = 1
	W[1] = 1
	W[2] = 1
	W[3] = 1

	k := 0
	n100 := false

	for {
		RandomX()
		Calculate()

		if k >= 100 {
			k = 0
			n100 = true
		}

		errN[k] = N.err * N.err
		k++

		sn := float64(0)
		for n := 0; n < 100; n++ {
			sn += errN[n]

		}
		if n100 && sn < 1 {
			break
		}
	}

	for {
		RandomX()
		Calculate()
		fmt.Println(Answer(), N.out)
		s := 0
		fmt.Scanln(&s)
	}
}

func Activate(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func Answer() float64 {
	sx := float64(0)
	sx += X[0]
	sx += X[1]
	sx += X[2]
	return sx / 3
}

func RandomX() {
	for x := 0; x < 3; x++ {
		X[x] = float64(rand.Intn(2))
	}
}

func Calculate() {
	N.in = 0
	N.in += X[0] * W[0]
	N.in += X[1] * W[1]
	N.in += X[2] * W[2]
	N.in += X[3] * W[3] // Нейрон смещения

	N.out = Activate(N.in)
	N.err = Answer() - N.out

	W[0] += nc * N.err * X[0] * N.out * (1 - N.out)
	W[1] += nc * N.err * X[1] * N.out * (1 - N.out)
	W[2] += nc * N.err * X[2] * N.out * (1 - N.out)
	W[3] += nc * N.err * X[3] * N.out * (1 - N.out)
}
