package main

import (
	"fmt"
	"math/rand"
	. "neuron/xo"
	"time"
)

const (
	NCorrect   = 0.15
	NBots      = 2
	LearPeriod = 10000
)

var (
	Bots                  [NBots]Bot
	Byzy, Lose, Win, Draw int
)

func main() {
	rand.Seed(time.Now().UnixNano())
	Layers = ([]int{9, 37, 9})
	Bots[0].NeuralNet.CreateLayer(Layers)
	Bots[0].NeuralNet.NCorrect = NCorrect

	Bots[1].NeuralNet.CreateLayer(Layers)
	Bots[1].NeuralNet.NCorrect = NCorrect

	fmt.Println("Обучение нейросети N1:")
	fmt.Println()
	LearnMove(0)
	fmt.Println()
	LearGameLevel0(0)

	fmt.Println()
	fmt.Println("Обучение нейросети N2:")
	fmt.Println()
	LearnMove(1)
	fmt.Println()
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
		Bots[bot].CorrectLose(b, r)
		Lose++
	}

	XO = XO0
}

func LearGameLevel0(bot int) {
	for g := 0; g < 100; g++ {
		Byzy = 0
		Lose = 0
		Win = 0
		Draw = 0
		Step := 0
		for n := 0; n < LearPeriod; n++ {
			var r, b int
			Bots[bot].NHistory = 0

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

				Bots[bot].History[Bots[bot].NHistory].XO = XO
				Bots[bot].History[Bots[bot].NHistory].Move = b
				Bots[bot].History[Bots[bot].NHistory].EnemyMove = r
				Bots[bot].NHistory++
			}
			if Step == 201 {
				Bots[bot].CorrectByzy()
				Byzy++
			}
			if Step == 1 {
				Bots[bot].CorrectLose(b, r)

				fLose := float64(Lose + LearPeriod/100)
				fN := float64(n)
				if 0.05 > fLose/fN {
					B := Bots[bot].History[0].Move
					XO = Bots[bot].History[0].XO
					XO[B] = 0
					Bots[bot].Move()
					R := RandomMove()
					if XO[4] == 0 {
						R = 4
					}
					Bots[bot].CorrectLose(B, R)
				}

				Lose++
			}
			if Step == 2 {
				Bots[bot].CorrectWin(b)
				Win++
			}
			if Step == 3 {
				Draw++
			}
			XO = XO0
		}
		fmt.Println("Ошибок:", Byzy, " Поражений:", Lose, " Побед:", Win, " Ничьих:", Draw, "/", LearPeriod)

		if Byzy < 10 && Lose < 100 {
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
