package main

import (
	"fmt"
	. "neuron/xo"
)

const (
	NCorrect   = 0.15
	NBots      = 2
	LearPeriod = 10000
)

var (
	Bots       [NBots]Bot
	Byzy, Lose int
)

func main() {
	Layers = ([]int{9, 72, 9})
	Bots[0].NeuralNet.CreateLayer(Layers)
	Bots[0].NeuralNet.NCorrect = NCorrect

	Bots[1].NeuralNet.CreateLayer(Layers)
	Bots[1].NeuralNet.NCorrect = NCorrect

	fmt.Println("Обучение нейросети N1:")
	LearnMove(0)
	LearGameLevel0(0)

	fmt.Println("Обучение нейросети N2:")
	LearnMove(1)
	LearGameLevel0(1)

	for {
		HumanOneGame(0)
		HumanOneGame(1)
	}
}

func HumanOneGame(bot int) {
	var Step int

	fmt.Println()
	var r, b int

	if bot%2 == 1 {
		r = RandomMove()
		Step = GameStep(r, 1)
		PrintXO()
	}

	for {
		fmt.Print("Ваш ход (введите номер клетки от 0 до 8): ")
		fmt.Scanln(&r)
		Step = GameStep(r, -1)
		PrintXO()
		if Step != 0 {
			break
		}
		b = Bots[bot].Move()
		Step = GameStep(b, 1)
		PrintXO()
		if Step != 0 {
			break
		}
	}
	fmt.Println(Results(Step))
	if Step == 201 {
		Bots[bot].CorrectByzy()
		Byzy++
	}
	if Step == 1 {
		Bots[bot].CorrectLoseLevel0(b, r)
		Lose++
	}

	XO = XO0
}

func LearGameLevel0(bot int) {
	for {
		Byzy = 0
		Lose = 0
		Step := 0
		for n := 0; n < LearPeriod; n++ {
			var r, b int

			if bot%2 == 1 {
				r = RandomMove()
				Step = GameStep(r, 1)
			}

			for {
				r = RandomMove()
				Step = GameStep(r, -1)
				if Step != 0 {
					break
				}
				b = Bots[bot].Move()
				Step = GameStep(b, 1)
				if Step != 0 {
					break
				}
			}
			if Step == 201 {
				Bots[bot].CorrectByzy()
				Byzy++
			}
			if Step == 1 {
				Bots[bot].CorrectLoseLevel0(b, r)
				Lose++
			}

			if Step == 2 {
				Bots[bot].CorrectWin(b, r)
			}

			XO = XO0
		}
		fmt.Println("Ошибок", Byzy, "и поражений", Lose, "на", LearPeriod)
		if Byzy == 0 && Lose == 0 {
			break
		}
	}

}

func LearnMove(bot int) {
	for {
		Byzy = 0
		for n := 0; n < 1000; n++ {
			for {
				r := RandomMove()
				Step := GameStep(r, -1)
				if Step != 0 {
					break
				}
				b := Bots[bot].Move()
				Step = GameStep(b, 1)
				if Step == 201 {
					Bots[bot].CorrectByzy()
					Byzy++
				}
				if Step != 0 {
					break
				}
			}
			XO = XO0
		}
		fmt.Println("Ошибок", Byzy, "на", LearPeriod)
		if Byzy == 0 {
			break
		}
	}

}
