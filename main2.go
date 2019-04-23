package main

import (
	"fmt"
	"math/rand"
	. "neuron/xo"
	"time"
)

const (
	DefCorrect = 0.15
)

var (
	g GameField
)

func main() {
	rand.Seed(time.Now().UnixNano())
	g.Prepare([]int{9, 36, 9}, DefCorrect)

	for {
		for n := 0; n < 10000; n++ {
			g.Step()
		}
		fmt.Println(g.OBot.Byzy, g.XBot.Byzy)
		if g.XBot.Byzy == 0 && g.OBot.Byzy == 0 {
			break
		}
		g.XBot.Byzy = 0
		g.OBot.Byzy = 0
	}

	for n := 0; n < 10; n++ {
		g.Step()
		fmt.Println(Results(g.StepRes))
		g.PrintXO()
	}

}
