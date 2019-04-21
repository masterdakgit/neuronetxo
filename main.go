package main

import (
	"fmt"
	"log"
	"math/rand"
	"neuron/nr"
)

const (
	Period    = 100000
	NCorrect  = 0.3
	NNCorrect = 1
)

var (
	XO, XO0 [9]float64
	History [4]HistoryPlus
	XA      []float64
	Layers  []int
	H       int
	Lose    int
	Byzy    int
	End     bool
	nn0     nr.NeuroNet
)

type NeuralNet struct {
}

type HistoryPlus struct {
	H [9]float64
	N int
}

func main() {
	//rand.Seed(time.Now().UnixNano())
	Layers = ([]int{9, 37, 9})
	nn0.CreateLayer(Layers)
	nn0.NCorrect = NCorrect

	N := 0
	End = false
	for n := 0; n < Period; n++ {
		if n%(Period/100) == 0 {
			fmt.Println(n*100/Period+1, "%", "Поражений:", Lose, "Ошибок:",
				Byzy, "/", Period/100)
			Byzy = 0
			Lose = 0
		}
	Start:
		Nx := X()

		if H >= 4 {
			H = 0
			XO = XO0
			goto Start
		}

		History[H].H = XO

		if Winer(-1)[:1] == "w" {
			CorrectLose(N, Nx)
			Lose++
			H = 0
			goto Start
		}

	StartO:

		XOToOut()
		nn0.Calc()
		N = O()

		if N == 102 {
			Lose++
			H = 0
			goto Start
		}

		History[H].N = N

		s := Winer(1)
		if s[:1] == "w" {
			CorrectWin(N)
			H = 0
			goto StartO
		}

		H++
	}
	fmt.Println()
	End = true

	H = 0
	XO = XO0
	for f := 0; f < 1000; f++ {
	StartX:
		r := 0
		fmt.Print("Ваш ход: ")
		fmt.Scanln(&r)
		Nx := r
		XO[r] = -1
		XOPrint()

		if H >= 4 {
			H = 0
			XO = XO0
			goto StartX
		}

		History[H].H = XO

		if Winer(-1)[:1] == "w" {
			fmt.Println(N, Nx)
			fmt.Println("Вы победили!")
			CorrectLose(N, Nx)
			H = 0
			goto StartX

		}

	StartO_:
		XOToOut()
		nn0.Calc()

		N := O()
		s := Winer(1)
		XOPrint()

		if s[:1] == "w" {
			//XOPrint()
			CorrectWin(N)
			H = 0
			fmt.Println("Победа ИИ:", s)
			goto StartO_
		}

		if N == 0 {
			//XOPrint()
		} else {
			if N == 101 {
				fmt.Println("Конец.")
				H = 0
			}
			if N == 102 {
				fmt.Println("Ошибка: ИИ сходил на занятую клетку.")
				H = 0
				s := ""
				fmt.Scanln(&s)
			}
		}
		History[H].N = N
		H++
	}

}

func O() int {
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

	max := float64(0)
	N := 0
	for n := 0; n < 9; n++ {
		if nn0.Layers[len(Layers)-1][n].Out > max {
			max = nn0.Layers[len(Layers)-1][n].Out
			N = n
		}
	}

	if XO[N] == 0 {
		XO[N] = 1
	} else {
		Byzy++
		CorrectByzy(N)
		H = 0
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
func XOToOut() {
	for n := 0; n < 9; n++ {
		nn0.Layers[0][n].Out = XO[n]
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

func CorrectWin(N int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		XA[n] = nn0.Layers[len(Layers)-1][n].Out
	}
	XA[N] = 1
	nn0.SetAnswers(XA)
	nn0.Correct()

	//ИИ победил, начинаем сначала.
	XO = XO0
}

func CorrectLose(N, Nx int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		XA[n] = nn0.Layers[len(Layers)-1][n].Out
	}
	XA[N] = 0
	XA[Nx] = 1
	nn0.SetAnswers(XA)
	nn0.Correct()

	if End {
		XOPrint()
		fmt.Println(N, Nx)
	}
	for n := H - 1; n >= 0; n-- {
		XO = History[n].H
		N = History[n].N

		if End {
			XOPrint()
			fmt.Println(N)
			fmt.Println("-------------------------")
		}

		XOToOut()
		nn0.Calc()
		NOld := O()

		if End {
			XOPrint()
			fmt.Println(NOld, N)
			fmt.Printf("%.5f", nn0.Layers[len(Layers)-1][N].Out)
			fmt.Println()
			for x := 0; x < 9; x++ {
				fmt.Printf("%.5f", nn0.Layers[len(Layers)-1][x].Out)
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
				XA[n] = nn0.Layers[len(Layers)-1][n].Out
			}
			XA[N] = 0
			nn0.SetAnswers(XA)
			nn0.NCorrect *= NNCorrect
			nn0.Correct()
		} else {

			continue
		}
	}
	nn0.NCorrect = NCorrect

	//ИИ проиград, начинаем сначала.
	XO = XO0
}

func CorrectByzy(N int) {
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

	nn0.SetAnswers(XA)
	nn0.Correct()

	//Сходил на занятую клетку, начинаем сначала.
	XO = XO0
}
