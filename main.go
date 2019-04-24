package main

import (
	"fmt"
	"math/rand"
	. "neuron/xo"
	"time"
)

const (
	DefCorrect = 0.2
	Cores      = 2
)

var (
	g          [Cores]GameField
	NBotsLearn int
	MaxLose    int
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Укажите допустимое количество поражений ИИ на 1000 игр: ")
	fmt.Scanln(&MaxLose)

	fmt.Println("Начинаем обучение...")

	for core := 0; core < Cores; core++ {
		g[core].Prepare([]int{9, 57, 3, 1}, DefCorrect, core)
		go RandomBattle(core)
	}

	fmt.Scanln()

	for {
		g[0].NStep = 0
		for {
			g[0].StepHuman()
			PrintXO(g[0].XO)
			if g[0].StepRes != 0 {
				fmt.Println(Results(g[0].StepRes))
				break
			}
		}

		g[1].NStep = 0
		for {
			g[1].StepHuman()
			PrintXO(g[1].XO)
			if g[1].StepRes != 0 {
				fmt.Println(Results(g[1].StepRes))
				break
			}
		}
	}
}

func RandomBattle(core int) {
	for r := 0; r < 100; r++ {
		for n := 0; n < 1000; n++ {
			for {
				g[core].StepRandom()
				if g[core].StepRes != 0 {
					break
				}
			}
		}

		fmt.Println("Bot", core, "- Lose:", g[core].OBot.Lose, " Win:", g[core].OBot.Win,
			" Draw:", g[core].OBot.Draw)
		if g[core].OBot.Lose <= MaxLose {
			break
		}

		g[core].OBot.Lose = 0
		g[core].OBot.Win = 0
		g[core].OBot.Draw = 0
	}

	fmt.Println("Bot", core, "закончил обучение.")

	NBotsLearn++

	if NBotsLearn > 1 {
		fmt.Println("Обучение завершено, нажмите Enter для игры с ботами.")
	}
}
