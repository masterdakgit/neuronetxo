package xo

import (
	"fmt"
	"log"
	"math/rand"
	"neuron/nr"
	"strconv"
)

const (
	xx = -1
	oo = 1
)

var (
	XO0 [9]float64
)

type HistoryMove struct {
	XO        [9]float64
	Move      int
	EnemyMove int
}

type GameField struct {
	XA         []float64
	Layers     []int
	XO         [9]float64
	XBot, OBot Bot
	xoLast     float64
	movLast    int
	StepRes    int
	NStep      int
}

type Bot struct {
	NeuralNet             nr.NeuroNet
	History               [5]HistoryMove
	NHistory              int
	Lose, Win, Byzy, Draw int
	LastMove              int
	xo                    float64
}

func (game *GameField) Prepare(L []int, defCorrect float64) {
	game.Layers = L

	game.XBot.NeuralNet.CreateLayer(L)
	game.XBot.NeuralNet.NCorrect = defCorrect
	game.XBot.xo = xx

	game.OBot.NeuralNet.CreateLayer(L)
	game.OBot.NeuralNet.NCorrect = defCorrect
	game.OBot.xo = oo

}

func (g *GameField) RandomMove(bot *Bot) {
	if g.NoMove() {
		log.Fatal("RandomMove: Нет свободных клеток.")
	}
	r := rand.Intn(9)
	for {
		if g.XO[r] == 0 {
			break
		}
		r = (r + 1) % 9
	}
	bot.LastMove = r
	g.step(bot)
	g.NStep++

}

func (g *GameField) NoMove() bool {
	NoMove := true
	for n := 0; n < 9; n++ {
		if g.XO[n] == 0 {
			NoMove = false
		}
	}
	return NoMove
}

func (g *GameField) Step() {
	g.StepRes = 0

	if g.NStep == 0 {
		g.XO = XO0
		g.RandomMove(&g.XBot)
		g.NStep++
		return
	}
	if g.NStep%2 == 1 {

		g.Move(&g.OBot)
		if g.StepRes != 0 {
			g.Correcting(&g.OBot, &g.XBot)
			g.NStep = 0
			return
		}

	} else {
		g.Move(&g.XBot)
		if g.StepRes != 0 {
			g.Correcting(&g.XBot, &g.OBot)
			g.NStep = 0
			return
		}
	}
	g.NStep++
}

func (g *GameField) step(bot *Bot) {
	mov := bot.LastMove
	xo := bot.xo
	if mov < 0 || mov > 8 {
		g.StepRes = 202
		return
	}

	if g.XO[mov] != 0 {
		g.StepRes = 201
		return
	}

	if xo != xx && xo != oo {
		g.StepRes = 203
		return
	}

	g.XO[mov] = xo

	if g.Winer(xo) {
		if xo == xx {
			g.StepRes = 1
		} else {
			g.StepRes = 2
		}
		return
	}

	if g.NoMove() {
		g.StepRes = 3
		return
	}

	g.StepRes = 0
}

func (g *GameField) Winer(w float64) bool {
	for x := 0; x < 3; x++ {
		o := 0
		for y := 0; y < 3; y++ {
			n := y*3 + x
			if g.XO[n] == w {
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
			if g.XO[n] == w {
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
		if g.XO[n] == w {
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
		if g.XO[n] == w {
			o++
		}
	}
	if o == 3 {
		return true
	}

	return false
}

func (g *GameField) Move(bot *Bot) {
	for n := 0; n < 9; n++ {
		bot.NeuralNet.Layers[0][n].Out = g.XO[n]
	}
	bot.NeuralNet.Calc()
	bot.LastMove = bot.NeuralNet.MaxOutputNumber()
	g.step(bot)

}

func (g *GameField) PrintXO() {
	for n := 0; n < 9; n++ {
		if n%3 == 0 {
			fmt.Println()
		}
		switch g.XO[n] {
		case 0:
			fmt.Print(". ")
		case xx:
			fmt.Print("x ")
		case oo:
			fmt.Print("o ")
		}
	}
	fmt.Println("\n")
}

func (g *GameField) CorrectByzy(bot *Bot) {
	g.XA = make([]float64, 9)
	if g.NoMove() {
		log.Fatal("CorrectByzy: ИИ некуда ходить.")
	}
	for n := 0; n < 9; n++ {
		if g.XO[n] == 0 {
			g.XA[n] = 1
		} else {
			g.XA[n] = 0
		}
	}
	bot.NeuralNet.SetAnswers(g.XA)
	bot.NeuralNet.Correct()
	bot.Byzy++
}

func (g *GameField) CorrectLose(my, enymy *Bot) {
	g.XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		g.XA[n] = my.NeuralNet.Layers[len(g.Layers)-1][n].Out
	}
	g.XA[my.LastMove] = 0
	g.XA[enymy.LastMove] = 1
	my.NeuralNet.SetAnswers(g.XA)
	my.NeuralNet.Correct()
	my.Lose++
}

func (g *GameField) CorrectWin(bot *Bot) {
	g.XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		g.XA[n] = bot.NeuralNet.Layers[len(g.Layers)-1][n].Out
	}
	g.XA[bot.LastMove] = 1
	bot.NeuralNet.SetAnswers(g.XA)
	bot.NeuralNet.Correct()
	bot.Win++
}

func (g *GameField) CorrectDraw(bot *Bot) {
	g.XA = make([]float64, 9)
	for n := 0; n < 9; n++ {
		g.XA[n] = bot.NeuralNet.Layers[len(g.Layers)-1][n].Out
	}
	g.XA[bot.LastMove] = 0.5
	bot.NeuralNet.SetAnswers(g.XA)
	bot.NeuralNet.Correct()
	bot.Draw++
}

func (g *GameField) Correcting(my, enymy *Bot) {
	if g.StepRes == 201 {
		g.CorrectByzy(my)
	}
	if g.StepRes == 1 {
		g.CorrectLose(enymy, my)
		g.CorrectWin(my)
	}
	if g.StepRes == 2 {
		g.CorrectLose(enymy, my)
		g.CorrectWin(my)
	}
	if g.StepRes == 3 {
		g.CorrectDraw(my)
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
		return "203: \"хо\" должет быть " + strconv.Itoa(xx) + " либо " + strconv.Itoa(oo)
	default:
		return "Незарегистрированный результат."
	}
}
