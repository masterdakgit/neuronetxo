package nr

import (
	"math"
	"math/rand"
)

type NeuroNet struct {
	Layers   [][]Neuron
	Weight   [][][]float64
	NCorrect float64
	answers  []float64
}

type Neuron struct {
	In, Out, Err float64
}

type Output struct {
	Out float64
	N   int
}

func (nn *NeuroNet) CreateLayer(L []int) {
	nn.NCorrect = 0.3
	nn.Layers = make([][]Neuron, len(L))
	for n := 0; n < len(L); n++ {
		nn.Layers[n] = make([]Neuron, L[n]+1)
		nn.Layers[n][L[n]].Out = 1
	}

	nn.startWeight()
}

func (nn *NeuroNet) startWeight() {
	nn.Weight = make([][][]float64, len(nn.Layers)-1)

	for w := range nn.Weight {
		nn.Weight[w] = make([][]float64, len(nn.Layers[w]))

		for l0 := range nn.Weight[w] {
			nn.Weight[w][l0] = make([]float64, len(nn.Layers[w+1])-1)

			for l1 := range nn.Weight[w][l0] {
				nn.Weight[w][l0][l1] = float64((rand.Intn(1000) - 500)) / 1000
			}
		}
	}
}

func (nn *NeuroNet) Calc() {
	for L := 0; L < len(nn.Layers)-1; L++ {
		for xr := 0; xr < len(nn.Layers[L+1])-1; xr++ {
			nn.Layers[L+1][xr].In = 0
			for xl := 0; xl < len(nn.Layers[L]); xl++ {
				nn.Layers[L+1][xr].In += nn.Layers[L][xl].Out * nn.Weight[L][xl][xr]
				nn.Layers[L+1][xr].Activate()
			}
		}
	}
}

func (nn *NeuroNet) Correct() {
	nn.layerError()
}

func (nn *NeuroNet) layerError() {
	Nr := len(nn.Layers) - 1
	for xr := 0; xr < len(nn.Layers[Nr])-1; xr++ {
		nn.Layers[Nr][xr].Err = nn.answer(xr) - nn.Layers[Nr][xr].Out
	}

	for left := Nr - 1; left > 0; left-- {
		right := left + 1
		for xl := 0; xl < len(nn.Layers[left])-1; xl++ {
			nn.Layers[left][xl].Err = 0
			for xr := 0; xr < len(nn.Layers[right])-1; xr++ {
				nn.Layers[left][xl].Err += nn.Layers[right][xr].Err * nn.Weight[left][xl][xr]
			}
		}
	}
	nn.weightCorrect()
}

func (nn *NeuroNet) weightCorrect() {
	for w := 0; w < len(nn.Weight); w++ {
		for xl := 0; xl < len(nn.Layers[w]); xl++ {
			for xr := 0; xr < len(nn.Layers[w+1])-1; xr++ {
				nn.Weight[w][xl][xr] += nn.weightChange(nn.Layers[w][xl], nn.Layers[w+1][xr])
			}
		}
	}
}

func (nn *NeuroNet) weightChange(NLeft, NRight Neuron) float64 {
	return nn.NCorrect * NRight.Err * NRight.Out * (1 - NRight.Out) * NLeft.Out
}

func (nn *NeuroNet) SetAnswers(a []float64) {
	nn.answers = make([]float64, len(a))
	for n := 0; n < len(a); n++ {
		nn.answers[n] = a[n]
	}
}

func (nn *NeuroNet) answer(n int) float64 {
	return nn.answers[n]
}

func (n *Neuron) Activate() {
	n.Out = 1 / (1 + math.Exp(-n.In))
}

func (nn *NeuroNet) SortOutput() []Output {
	Out := make([]Output, len(nn.Layers[len(nn.Layers)-1])-1)

	for n := 0; n < len(Out); n++ {
		Out[n].Out = nn.Layers[len(nn.Layers)-1][n].Out
		Out[n].N = n
	}

	for n := 0; n < len(Out)-1; n++ {
		for m := n; m < len(Out); m++ {
			if Out[n].Out < Out[m].Out {
				O := Out[m].Out
				N := Out[m].N
				Out[m].Out = Out[n].Out
				Out[m].N = Out[n].N
				Out[n].Out = O
				Out[n].N = N
			}
		}
	}

	return Out
}
