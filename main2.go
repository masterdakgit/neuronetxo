package main

import (
	"math/rand"
	"neuron/nr"
	. "neuron/xo"
)

var (
	Bots [NBots]Bot
)

type Bot struct {
	NeuralNet nr.NeuroNet
	History   HistoryMove
}

func main() {
	Layers = ([]int{9, 37, 9})
	Bots[0].NeuralNet.CreateLayer(Layers)
	Bots[0].NeuralNet.NCorrect = NCorrect

}

func RandomStep() (move int, ret string) {
	NoMove := true
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			NoMove = false
		}
	}
	if NoMove {
		return 101, "Нет свободных клеток."
	}
	r := rand.Intn(9)
	for {
		if XO[r] == 0 {
			XO[r] = -1
			break
		}
		r = (r + 1) % 9
	}
	return r, "Случайный ход."
}
