package main

import (
	"fmt"
	"log"
	"math/rand"
	"neuron/nr"
	"time"
)

const (
	Period    = 1000000
	NCorrect  = 0.9
	NNCorrect = 1
	NBots     = 2
)

var (
	XO, XO0   [9]float64
	XA        []float64
	Layers    []int
	N         int
	Byzy      int
	HumanStep bool
	End       bool

	HumanNet GameNeuralNet
	bn       [NBots]GameNeuralNet
)

type GameNeuralNet struct {
	nn                       nr.NeuroNet
	History                  [5]HistoryPlus
	xo                       float64
	H                        int
	Lose, Win, Draw, ErrByzy int
	LastMove                 int
}

type HistoryPlus struct {
	XO   [9]float64
	Move int
}

func main() {

	rand.Seed(time.Now().UnixNano())
	Layers = ([]int{9, 37, 9})
	for n := 0; n < len(bn); n++ {
		bn[n].nn.CreateLayer(Layers)
		bn[n].nn.NCorrect = 0.15
	}

	for n0 := 0; n0 < len(bn)-1; n0++ {
		for n1 := n0 + 1; n1 < len(bn); n1++ {
			bn[n0].xo = 1
			bn[n1].xo = -1
			XO = XO0
			bn[n0].H = 0
			bn[n1].H = 0
			Game(&bn[n0], &bn[n1])
			fmt.Println(100*(n1+n0*(NBots-1))/((NBots-1)*(NBots-1)), "%")
		}
	}

	maxWin := 0
	NWiner := 0
	for n := 0; n < len(bn); n++ {
		if maxWin < bn[n].Win {
			maxWin = bn[n].Win
			NWiner = n
		}
	}

	fmt.Println("Победил бот N", NWiner, "набрав", maxWin, "побед.")

	HumanNet = bn[0]
	HumanNet.xo = -1
	bn[NWiner].xo = 1
	GameWithHuman(&HumanNet, &bn[NWiner])
}

func Game(gn0, gn1 *GameNeuralNet) {
	End = false
	for n := 0; n < Period; n++ {
		if n%(Period/100) == 0 {
			fmt.Println()
			fmt.Println(n*100/Period+1, "%")
			fmt.Println("ИИ_0  Win:", gn0.Win, " Lose:", gn0.Lose, " Draw:", gn0.Draw,
				" ErrByzy:", gn0.ErrByzy, "/", Period/100)
			fmt.Println("ИИ_1  Win:", gn1.Win, " Lose:", gn1.Lose, " Draw:", gn1.Draw,
				" ErrByzy:", gn1.ErrByzy, "/", Period/100)

			//fmt.Println("ИИ0 ErrByzy:", gn0.ErrByzy, "/", Period/5)
			//fmt.Println("ИИ1 ErrByzy:", gn1.ErrByzy, "/", Period/5)
			if n > 0 && (gn0.Win+gn1.Win) == 0 {
				break
			}

			//	gn0.ErrByzy = 0
			//	gn1.ErrByzy = 0

			gn0.ErrByzy = 0
			gn0.Win = 0
			gn0.Draw = 0
			gn0.Lose = 0

			gn1.ErrByzy = 0
			gn1.Win = 0
			gn1.Draw = 0
			gn1.Lose = 0

		}
		gn1.GameMove(gn0)
		gn0.GameMove(gn1)

	}
	End = true
}

func (gn *GameNeuralNet) GameMove(enemy *GameNeuralNet) {
	if HumanStep {
		HumanMove()
	} else {
		gn.XOToOut()
		gn.nn.Calc()
		gn.LastMove = gn.Move()
	}

	if enemy.LastMove > 100 || HumanNet.LastMove > 100 {
		HumanNet.LastMove = 50
		gn.H = 0
		enemy.H = 0
		XO = XO0
		return
	}

	if End {
		XOPrint()
		if !HumanStep {
			gn.PrintOut()
		}
	}

	//Некуда ходить
	if gn.LastMove == 101 {
		if End {
			fmt.Println()
			fmt.Println("Ничья!")
		}
		gn.Draw++
		enemy.Draw++

		XO = gn.History[gn.H-1].XO
		enemy.XOToOut()
		enemy.nn.Calc()
		enemy.CorrectDraw(enemy.LastMove)

		XO = enemy.History[gn.H-1].XO
		gn.XOToOut()
		gn.nn.Calc()
		gn.CorrectDraw(gn.History[gn.H-1].Move)

		XO = XO0
		return
	}

	//Сходил на занятую клетку
	if gn.LastMove == 102 {
		if End {
			fmt.Println()
			fmt.Println("Ошибка: ИИ сходил на занятую клетку.")
		}

		gn.ErrByzy++
		gn.Lose++
		enemy.Win++

		XO = XO0
		return
	}

	gn.History[gn.H].XO = XO
	gn.History[gn.H].Move = gn.LastMove

	//Если победил
	if Winer(gn.xo)[:1] == "w" {
		if End {
			fmt.Println()
			if HumanStep {
				fmt.Println("Вы победили!")
			} else {
				fmt.Println("ИИ победил, Вы проиграли.")
			}
		}

		gn.CorrectWin(gn.LastMove)
		gn.Win++
		enemy.Lose++
		enemy.CorrectLose(enemy.LastMove, gn.LastMove)

		XO = XO0
		gn.LastMove = 103
		return
	}

	//Если проиграл
	W := float64(0)
	if gn.xo == 1 {
		W = -1
	} else {
		W = 1
	}
	if Winer(W)[:1] == "w" {
		if End {
			fmt.Println()
			if HumanStep {
				fmt.Println("ИИ победил, Вы проиграли.")
			} else {
				fmt.Println("Вы победили!")
			}
		}
		gn.CorrectLose(gn.LastMove, enemy.LastMove)
		gn.Lose++
		enemy.Win++
		enemy.CorrectWin(enemy.LastMove)

		XO = XO0
		gn.LastMove = 104
		return
	}

	gn.H++

}

func (gn *GameNeuralNet) Move() int {
	b := true
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			b = false
		}
	}

	if b {
		XO = XO0
		//log.Fatal("Ошибка: ИИ некуда ходить.")
		return 101
	}

	N := gn.nn.SortOutput()[0].N

	if gn.H == 0 && XO == XO0 {
		r := rand.Intn(9)
		N = r
	}

	if XO[N] == 0 {
		XO[N] = gn.xo
	} else {
		Byzy++
		gn.CorrectByzy(N)
		gn.H = 0
		return 102
	}
	return N
}

func X() int {
	b := true
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			b = false
		}
	}

	if b {
		return 103
	}

	r := 0
	for {
		r = rand.Intn(9)
		if XO[r] == 0 {
			XO[r] = -1
			break
		}
	}
	return r
}
func (gn *GameNeuralNet) XOToOut() {
	for n := 0; n < 9; n++ {
		gn.nn.Layers[0][n].Out = XO[n]
	}
}

func XOPrint() {
	for n := 0; n < 9; n++ {
		if n%3 == 0 {
			fmt.Println()
		}
		switch XO[n] {
		case 0:
			fmt.Print(". ")
		case -1:
			fmt.Print("x ")
		case 1:
			fmt.Print("o ")
		}
	}
	fmt.Println("\n")
}

func Winer(w float64) string {
	for x := 0; x < 3; x++ {
		o := 0
		for y := 0; y < 3; y++ {
			n := y*3 + x
			if XO[n] == w {
				o++
			}
		}
		if o == 3 {
			return "w|"
		}
	}

	for y := 0; y < 3; y++ {
		o := 0
		for x := 0; x < 3; x++ {
			n := y*3 + x
			if XO[n] == w {
				o++
			}
		}
		if o == 3 {
			return "w-"
		}
	}

	o := 0
	for m := 0; m < 3; m++ {
		x := m
		y := m
		n := y*3 + x
		if XO[n] == w {
			o++
		}
	}
	if o == 3 {
		return "wX"
	}

	o = 0
	for x := 2; x >= 0; x-- {
		y := 2 - x
		n := y*3 + x
		if XO[n] == w {
			o++
		}
	}
	if o == 3 {
		return "w/"
	}

	return "--"
}

func (gn *GameNeuralNet) CorrectWin(N int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		XA[n] = gn.nn.Layers[len(Layers)-1][n].Out
	}
	XA[N] = 1
	gn.nn.SetAnswers(XA)
	gn.nn.Correct()

	//ИИ победил, начинаем сначала.
	XO = XO0
}

func (gn *GameNeuralNet) CorrectLose(N, Nx int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		XA[n] = gn.nn.Layers[len(Layers)-1][n].Out
	}
	XA[N] = 0
	XA[Nx] = 1
	gn.nn.SetAnswers(XA)
	gn.nn.Correct()

	EndH := false

	if EndH {
		XOPrint()
		fmt.Println(N, Nx)
	}
	for n := gn.H - 1; n >= 0; n-- {
		XO = gn.History[n].XO
		N = gn.History[n].Move

		if N > 100 {
			fmt.Println("Ошибка: ", N)
			continue
		}

		if EndH {
			XOPrint()
			fmt.Println(N)
			fmt.Println("-------------------------")
		}

		gn.XOToOut()
		gn.nn.Calc()
		NOld := gn.Move()

		if EndH {
			XOPrint()
			fmt.Println(NOld, N)
			fmt.Printf("%.5f", gn.nn.Layers[len(Layers)-1][N].Out)
			fmt.Println()
			for x := 0; x < 9; x++ {
				fmt.Printf("%.5f", gn.nn.Layers[len(Layers)-1][x].Out)
				fmt.Print(" ")
			}
			fmt.Println()
		}

		if N == NOld {
			if EndH {
				fmt.Println("Ходит по-старому.")
				qwe := 0
				fmt.Scanln(&qwe)
			}

			XA = make([]float64, 9)
			for n := 0; n < 9; n++ {
				XA[n] = gn.nn.Layers[len(Layers)-1][n].Out
			}
			XA[N] = 0
			gn.nn.SetAnswers(XA)
			gn.nn.NCorrect *= NNCorrect
			gn.nn.Correct()
		} else {

			continue
		}
	}
	gn.nn.NCorrect = NCorrect

	//ИИ проиград, начинаем сначала.
	XO = XO0
}

func (gn *GameNeuralNet) CorrectByzy(N int) {
	XA = make([]float64, 9)
	E := true
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			XA[n] = 1
			E = false
		} else {
			XA[n] = 0
		}
	}
	if E {
		log.Fatal("Ошибка: ИИ некуда ходить.")
	}

	gn.nn.SetAnswers(XA)
	gn.nn.Correct()

	//Сходил на занятую клетку, начинаем сначала.
	XO = XO0
}

func (gn *GameNeuralNet) CorrectDraw(N int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		XA[n] = gn.nn.Layers[len(Layers)-1][n].Out
	}

	/*	XOPrint()
		fmt.Println(N)
		r := 0
		fmt.Scanln(&r)
	*/
	XA[N] = 0.5
	gn.nn.SetAnswers(XA)
	gn.nn.Correct()

	//Ничья, начинаем сначала.
	XO = XO0
}

func HumanMove() {
	NoMove := true
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			NoMove = false
			break
		}
	}

	if NoMove {
		fmt.Print("Ничья.")
		HumanNet.LastMove = 203
		return
	}

	fmt.Println()
	fmt.Print("Ваш ход: ")
	r := 0
	fmt.Scanln(&r)

	if (r > 8) || (r < 0) {
		fmt.Println("Ошибка: Значение должно быть от 0 до 8.")
		HumanNet.LastMove = 201
		return
	}

	if XO[r] != 0 {
		fmt.Println("Ошибка: Клетка занята.")
		HumanNet.LastMove = 202
		return
	}

	XO[r] = -1
	HumanNet.LastMove = r
}

func GameWithHuman(gnHum, gnII *GameNeuralNet) {
	XO = XO0
	gnHum.H = 0
	gnII.H = 0

	for {
		HumanStep = true
		gnHum.GameMove(gnII)

		HumanStep = false
		gnII.GameMove(gnHum)
	}
}

func (gn *GameNeuralNet) PrintOut() {
	for n := 0; n < 9; n++ {
		fmt.Print(gn.nn.SortOutput()[n].N, " - ")
		fmt.Printf("%.3f", gn.nn.SortOutput()[n].Out)
		if n < 8 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
}
