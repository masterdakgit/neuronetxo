package xo

const (
	NCorrect = 0.15
	NBots    = 1
)

var (
	XO, XO0 [9]float64
	Layers  []int
)

type HistoryMove struct {
	XO   [9]float64
	Move int
}
