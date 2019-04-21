package main

import (
	"fmt"
	"log"
	"math/rand"
	"neuron/nr"
)

const (
	Period    = 1000000
	NCorrect  = 0.9
	NNCorrect = 1
)

var (
	XO, XO0 [9]float64
	XA      []float64
	Layers  []int
	N       int
	Byzy    int
	End     bool

	nn0, nn1 GameNeuralNet
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
	//rand.Seed(time.Now().UnixNano())
	Layers = ([]int{9, 37, 9})
	nn0.nn.CreateLayer(Layers)
	nn0.nn.NCorrect = 0.99
	nn0.xo = -1

	Layers = ([]int{9, 37, 9})
	nn1.nn.CreateLayer(Layers)
	nn1.nn.NCorrect = 0.9
	nn1.xo = 1

	Game(&nn0, &nn1)

	fmt.Println(nn0.Win)

	if nn0.Win > nn1.Win {
		nn0.xo = 1
		fmt.Println("Игра с ИИ_0:")
		nn0.GameWithHuman()
	} else {
		fmt.Println("Игра с ИИ_1:")
		nn1.GameWithHuman()
	}
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
	gn.XOToOut()
	gn.nn.Calc()
	gn.LastMove = gn.Move()

	if enemy.LastMove > 100 {
		gn.H = 0
		enemy.H = 0
		XO = XO0
		return
	}

	//Некуда ходить
	if gn.LastMove == 101 {
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

func (gn *GameNeuralNet) GameWithHuman() {
	gn.H = 0
	XO = XO0
	for f := 0; f < 1000; f++ {
	StartX:
		r := 0
		fmt.Print("Ваш ход: ")
		fmt.Scanln(&r)
		Nx := r
		XO[r] = -1
		XOPrint()

		if gn.H >= 4 {
			gn.H = 0
			XO = XO0
			goto StartX
		}

		gn.History[gn.H].XO = XO

		if Winer(-1)[:1] == "w" {
			fmt.Println(N, Nx)
			fmt.Println("Вы победили!")
			gn.CorrectLose(N, Nx)
			gn.H = 0
			goto StartX

		}

	StartO_:
		gn.XOToOut()
		gn.nn.Calc()

		N := gn.Move()
		s := Winer(1)
		XOPrint()

		if s[:1] == "w" {
			//XOPrint()
			gn.CorrectWin(N)
			gn.H = 0
			fmt.Println("Победа ИИ:", s)
			goto StartO_
		}

		if N == 0 {
			//XOPrint()
		} else {
			if N == 101 {
				fmt.Println("Конец.")
				gn.H = 0
			}
			if N == 102 {
				fmt.Println("Ошибка: ИИ сходил на занятую клетку.")
				gn.H = 0
				s := ""
				fmt.Scanln(&s)
			}
		}
		gn.History[gn.H].Move = N
		gn.H++
	}

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

	max0 := float64(0)
	max1 := float64(0)
	max2 := float64(0)
	N := 0
	N0 := 0
	N1 := 0
	N2 := 0
	for n := 0; n < 9; n++ {
		if gn.nn.Layers[len(Layers)-1][n].Out > max0 {
			max2 = max1
			N2 = N1

			max1 = max0
			N1 = N0

			max0 = gn.nn.Layers[len(Layers)-1][n].Out
			N0 = n
		}
	}

	N = N0
	switch rand.Intn(3) {
	case 0:
		N = N0
	case 1:
		if max1 > 0.5 {
			N = N1
		}
	case 2:
		if max2 > 0.5 {
			N = N2
		}
	}

	if gn.H == 0 {
		for {
			N = rand.Intn(9)
			if XO[N] == 0 {
				break
			}
		}
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

	if End {
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

		if End {
			XOPrint()
			fmt.Println(N)
			fmt.Println("-------------------------")
		}

		gn.XOToOut()
		gn.nn.Calc()
		NOld := gn.Move()

		if End {
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
			if End {
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
