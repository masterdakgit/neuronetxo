package main

import (
	"fmt"
	"log"
	"math/rand"
	"neuron/nr"
)

const (
	Period    = 1000000
	NCorrect  = 0.99
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
)

type HistoryPlus struct {
	H [9]float64
	N int
}

func main() {
	//rand.Seed(time.Now().UnixNano())
	Layers = ([]int{9, 37, 9})
	nr.CreateLayer(Layers)
	nr.NCorrect = NCorrect

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

		XOToOut()
		nr.Calc()
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
			goto Start
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
			fmt.Println("ИИ проиграл.")
			CorrectLose(N, Nx)
			H = 0
			goto StartX

		}

		XOToOut()
		nr.Calc()

		N := O()
		s := Winer(1)
		XOPrint()

		if s[:1] == "w" {
			//XOPrint()
			CorrectWin(N)
			H = 0
			fmt.Println("Победа ИИ:", s)
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
		if nr.Layers[len(Layers)-1][n].Out > max {
			max = nr.Layers[len(Layers)-1][n].Out
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
		nr.Layers[0][n].Out = XO[n]
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
		XA[n] = nr.Layers[len(Layers)-1][n].Out
	}
	XA[N] = 1
	nr.SetAnswers(XA)
	nr.Correct()

	//ИИ победил, начинаем сначала.
	XO = XO0
}

func CorrectLose(N, Nx int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		XA[n] = nr.Layers[len(Layers)-1][n].Out
	}
	XA[N] = 0
	XA[Nx] = 1
	nr.SetAnswers(XA)
	nr.Correct()

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
		nr.Calc()
		NOld := O()

		if End {
			XOPrint()
			fmt.Println(NOld, N)
			fmt.Printf("%.5f", nr.Layers[len(Layers)-1][N].Out)
			fmt.Println()
			for x := 0; x < 9; x++ {
				fmt.Printf("%.5f", nr.Layers[len(Layers)-1][x].Out)
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
				XA[n] = nr.Layers[len(Layers)-1][n].Out
			}
			XA[N] = 0
			nr.SetAnswers(XA)
			nr.NCorrect *= NNCorrect
			nr.Correct()
		} else {

			continue
		}
	}
	nr.NCorrect = NCorrect

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

	nr.SetAnswers(XA)
	nr.Correct()

	//Сходил на занятую клетку, начинаем сначала.
	XO = XO0
}
