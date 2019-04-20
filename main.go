package main

import (
	"fmt"
	"math/rand"
	"neuron/nr"
	"time"
)

var (
	XO, XO0 [9]float64
	XA      []float64
	Layers  []int
)

func main() {
	rand.Seed(time.Now().UnixNano())
	Layers = ([]int{9, 27, 9})
	nr.CreateLayer(Layers)
	nr.NCorrect = 0.2

	N := 0
	for n := 0; n < 1000000; n++ {
		Nx := X()
		if Winer(-1)[:1] == "w" {
			CorrectLose(N, Nx)
		}

		XOToOut()
		nr.Calc()
		N = O()

		s := Winer(1)
		if s[:1] == "w" {
			CorrectWin(N)
		}

	}

	XO = XO0
	for f := 0; f < 100; f++ {
	Start:
		Nx := X()
		XOPrint()
		if Winer(-1)[:1] == "w" {
			fmt.Println(N, Nx)
			fmt.Println("ИИ проиграл.")
			CorrectLose(N, Nx)
			s := ""
			fmt.Scanln(&s)
			goto Start

		}
		XOToOut()
		nr.Calc()

		N := O()
		s := Winer(1)

		if s[:1] == "w" {
			XOPrint()
			CorrectWin(N)
			fmt.Println("Победа ИИ:", s)
			//fmt.Scanln(&s)
		}

		if N == 0 {
			XOPrint()
		} else {
			if N == 101 {
				fmt.Println("Конец.")
			}
			if N == 102 {
				fmt.Println("Ошибка: ИИ сходил на занятую клетку.")
				break
			}
		}
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
		CorrectByzy(N)
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

	//ИИ проиград, начинаем сначала.
	XO = XO0
}

func CorrectByzy(N int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			XA[n] = 1
		} else {
			XA[n] = 0
		}
	}
	nr.SetAnswers(XA)
	nr.Correct()

	//Сходил на занятую клетку, начинаем сначала.
	XO = XO0
}
