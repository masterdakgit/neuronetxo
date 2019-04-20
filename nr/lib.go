package nr

import (
	"math"
	"math/rand"
)

var (
	Layers   [][]Neuron
	Weight   [][][]float64
	NCorrect float64
	Answers  []float64
)

type Neuron struct {
	In, Out, Err float64
}

func CreateLayer(L []int) {
	NCorrect = 0.1
	Layers = make([][]Neuron, len(L))
	for n := 0; n < len(L); n++ {
		Layers[n] = make([]Neuron, L[n]+1)
		Layers[n][L[n]].Out = 1
	}

	startWeight()
}

func startWeight() {
	Weight = make([][][]float64, len(Layers)-1)

	for w := range Weight {
		Weight[w] = make([][]float64, len(Layers[w]))

		for l0 := range Weight[w] {
			Weight[w][l0] = make([]float64, len(Layers[w+1])-1)

			for l1 := range Weight[w][l0] {
				Weight[w][l0][l1] = float64((rand.Intn(1000) - 500)) / 1000
			}
		}
	}
}

func Calc() {
	for L := 0; L < len(Layers)-1; L++ {
		for xr := 0; xr < len(Layers[L+1])-1; xr++ {
			Layers[L+1][xr].In = 0
			for xl := 0; xl < len(Layers[L]); xl++ {
				Layers[L+1][xr].In += Layers[L][xl].Out * Weight[L][xl][xr]
				Layers[L+1][xr].Activate()
			}
		}
	}
	layerError()
}

func layerError() {
	Nr := len(Layers) - 1
	for xr := 0; xr < len(Layers[Nr])-1; xr++ {
		Layers[Nr][xr].Err = answer(xr) - Layers[Nr][xr].Out
	}

	for left := Nr - 1; left > 0; left-- {
		right := left + 1
		for xl := 0; xl < len(Layers[left])-1; xl++ {
			Layers[left][xl].Err = 0
			for xr := 0; xr < len(Layers[right])-1; xr++ {
				Layers[left][xl].Err += Layers[right][xr].Err * Weight[left][xl][xr]
			}
		}
	}
	weightCorrect()
}

func weightCorrect() {
	for w := 0; w < len(Weight); w++ {
		for xl := 0; xl < len(Layers[w]); xl++ {
			for xr := 0; xr < len(Layers[w+1])-1; xr++ {
				Weight[w][xl][xr] += weightChange(Layers[w][xl], Layers[w+1][xr])
			}
		}
	}
}

func weightChange(NLeft, NRight Neuron) float64 {
	return NCorrect * NRight.Err * NRight.Out * (1 - NRight.Out) * NLeft.Out
}

func SetAnswers(a []float64) {
	Answers = make([]float64, len(a))
	for n := 0; n < len(a); n++ {
		Answers[n] = a[n]
	}
}

func answer(n int) float64 {
	return Answers[n]
}

func (n *Neuron) Activate() {
	n.Out = 1 / (1 + math.Exp(-n.In))
}
