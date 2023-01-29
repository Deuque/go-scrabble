package main

import (
	"github.com/Deuque/slack-scrabble/controllers"
	"github.com/Deuque/slack-scrabble/executor"
)

func main() {
	// var se executor.ScrabbleExecutor = &executor.TerminalScrabbleExecutor{
	// 	Scrabbler: controllers.NewHttpScrabbler(),
	// }

	var se executor.ScrabbleExecutor = &executor.SlackScrabbleExecutor{
		Scrabbler: controllers.NewMockScrabbler(),
	}

	se.Init()
}
