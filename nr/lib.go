package nr

import (
	"math"
	"math/rand"
	"time"
)

var (
	Layers [][]Neuron
	Weight [][][]float64
)

type Neuron struct {
	In, Out, Err float64
}

func CreateLayer(L []int) {
	Layers = make([][]Neuron, len(L))
	for n := 0; n < len(L); n++ {
		Layers[n] = make([]Neuron, L[n]+1)
		Layers[n][L[n]].Out = 1
	}

	startWeight()
}

func startWeight() {
	rand.Seed(time.Now().UnixNano())
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
	//Убираем нейрон смещения у последнего слоя
	Layers[len(Layers)-1] = make([]Neuron, len(Layers[len(Layers)-1])-1)
}

func CalcWeight() {
	for L := 0; L < len(Layers)-1; L++ {
		for xr := 0; xr < len(Layers[L+1]); xr++ {
			Layers[L+1][xr].In = 0
			for xl := 0; xl < len(Weight[L]); xl++ {
				Layers[L+1][xr].In += Layers[L][xl].Out * Weight[L][xl][xr]
			}
		}
	}

}

func (n *Neuron) Activae() {
	n.Out = 1 / (1 + math.Exp(-n.In))
}
