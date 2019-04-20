package nr

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
	Weight = make([][][]float64, len(Layers)-1)

	for w := range Weight {
		Weight[w] = make([][]float64, len(Layers[w]))
		for l0 := range Layers[w] {
			Weight[w][l0] = make([]float64, len(Layers[w+1])-1)
		}
	}

	Layers[len(Layers)-1] = make([]Neuron, len(Layers[len(Layers)-1])-1)
}
