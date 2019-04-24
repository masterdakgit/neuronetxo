package main

import (
	"fmt"
	"math/rand"
	. "neuron/xo"
	"time"
)

const (
	DefCorrect = 0.3
	Cores      = 4
)

var (
	g     [Cores]GameField
	rscan int
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for core := 0; core < Cores; core++ {
		g[core].Prepare([]int{9, 57, 3, 1}, DefCorrect)
		go RandomBattle(core)
	}
	fmt.Scanln(&rscan)
}

func RandomBattle(core int) {
	for {
		for n := 0; n < 10000; n++ {
			g[core].StepRandom()
		}

		fmt.Println("Core:", core, "OLose:", g[core].OBot.Lose, "Win:", g[core].OBot.Win)
		if g[core].OBot.Lose == 0 {
			break
		}
		g[core].OBot.Lose = 0
		g[core].OBot.Win = 0
	}
	for {
		fmt.Println("Core: ", core, " - найдено решение.")
		time.Sleep(2000000000)
	}
}
