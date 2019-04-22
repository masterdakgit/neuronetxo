package main

import (
	"fmt"
	. "neuron/xo"
)

const LearPeriod = 10000

var (
	Bots       [NBots]Bot
	Byzy, Lose int
)

func main() {
	Layers = ([]int{9, 37, 9})
	Bots[0].NeuralNet.CreateLayer(Layers)
	Bots[0].NeuralNet.NCorrect = NCorrect
	fmt.Println("Обучение нейросети:")
	LearnMove()
	LearGameLevel0()
	HumanGame()
}

func HumanGame() {
	var Step int
	for {
		fmt.Println()
		var r, b int
		for {
			fmt.Print("Ваш ход (введите номер клетки от 0 до 8): ")
			fmt.Scanln(&r)
			Step = GameStep(r, -1)
			PrintXO()
			if Step != 0 {
				break
			}
			b = Bots[0].Move()
			Step = GameStep(b, 1)
			PrintXO()
			if Step != 0 {
				break
			}
		}
		fmt.Println(Results(Step))
		if Step == 201 {
			Bots[0].CorrectByzy()
			Byzy++
		}
		if Step == 1 {
			Bots[0].CorrectLoseLevel0(b, r)
			Lose++
		}

		XO = XO0

	}
}

func LearGameLevel0() {
	for {
		Byzy = 0
		Lose = 0
		Step := 0
		for n := 0; n < LearPeriod; n++ {
			var r, b int
			for {
				r = RandomMove()
				Step = GameStep(r, -1)
				if Step != 0 {
					break
				}
				b = Bots[0].Move()
				Step = GameStep(b, 1)
				if Step != 0 {
					break
				}
			}
			if Step == 201 {
				Bots[0].CorrectByzy()
				Byzy++
			}
			if Step == 1 {
				Bots[0].CorrectLoseLevel0(b, r)
				Lose++
			}

			XO = XO0
		}
		fmt.Println("Ошибок", Byzy, "и поражений", Lose, "на", LearPeriod)
		if Byzy == 0 && Lose == 0 {
			break
		}
	}

}

func LearnMove() {
	for {
		Byzy = 0
		for n := 0; n < 1000; n++ {
			for {
				r := RandomMove()
				Step := GameStep(r, -1)
				if Step != 0 {
					break
				}
				b := Bots[0].Move()
				Step = GameStep(b, 1)
				if Step == 201 {
					Bots[0].CorrectByzy()
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
