package xo

import (
	"fmt"
	"log"
	"math/rand"
	"neuron/nr"
)

var (
	XO, XO0 [9]float64
	XA      []float64
	Layers  []int
)

type HistoryMove struct {
	XO   [9]float64
	Move int
}

type Bot struct {
	NeuralNet nr.NeuroNet
	History   [5]HistoryMove
	NHistory  int
}

func RandomMove() int {
	if NoMove() {
		log.Fatal("RandomMove: Нет свободных клеток.")
	}
	r := rand.Intn(9)
	for {
		if XO[r] == 0 {
			break
		}
		r = (r + 1) % 9
	}
	return r
}

func NoMove() bool {
	NoMove := true
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			NoMove = false
		}
	}
	return NoMove
}

func GameStep(mov int, xo float64) int {
	if mov < 0 || mov > 8 {
		return 202
	}

	if XO[mov] != 0 {
		return 201
	}

	if xo != -1 && xo != 1 {
		return 203
	}

	XO[mov] = xo

	if Winer(xo) {
		if xo == -1 {
			return 1
		}
		if xo == 1 {
			return 2
		}
	}

	if NoMove() {
		return 3
	}

	return 0
}

func Winer(w float64) bool {
	for x := 0; x < 3; x++ {
		o := 0
		for y := 0; y < 3; y++ {
			n := y*3 + x
			if XO[n] == w {
				o++
			}
		}
		if o == 3 {
			return true
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
			return true
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
		return true
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
		return true
	}

	return false
}

func (Bot *Bot) Move() int {
	for n := 0; n < 9; n++ {
		Bot.NeuralNet.Layers[0][n].Out = XO[n]
	}
	Bot.NeuralNet.Calc()
	return Bot.NeuralNet.MaxOutputNumber()
}

func PrintXO() {
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

func (bot *Bot) CorrectByzy() {
	XA = make([]float64, 9)
	if NoMove() {
		log.Fatal("CorrectByzy: ИИ некуда ходить.")
	}
	for n := 0; n < 9; n++ {
		if XO[n] == 0 {
			XA[n] = 1
		} else {
			XA[n] = 0
		}
	}
	bot.NeuralNet.SetAnswers(XA)
	bot.NeuralNet.Correct()
}

func (bot *Bot) CorrectLoseLevel0(B, R int) {
	XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		XA[n] = bot.NeuralNet.Layers[len(Layers)-1][n].Out
	}
	XA[B] = 0
	XA[R] = 1
	bot.NeuralNet.SetAnswers(XA)
	bot.NeuralNet.Correct()
}

func (bot *Bot) ShowHistory() {
	for n := 0; n < bot.NHistory; n++ {
		XO = bot.History[n].XO
		PrintXO()
	}
}

func Results(res int) string {
	switch res {
	case 201:
		return "201: Ход на занятую клетку."
	case 202:
		return "202: Значение должно быть от 0 до 8"
	case 1:
		return "Победили крестики."
	case 2:
		return "Победили нолики."
	case 3:
		return "Ничья."
	case 0:
		return "0: Игра продолжается."
	case 203:
		return "203: \"хо\" должет быть -1 либо 1."
	default:
		return "Незарегистрированный результат."
	}
}
