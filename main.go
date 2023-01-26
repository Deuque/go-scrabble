package main

import (
	"github.com/Deuque/slack-scrabble/controllers"
	"github.com/Deuque/slack-scrabble/executor"
)

func main() {
	var se executor.ScrabbleExecutor = &executor.TerminalScrabbleExecutor{
		Word:      "",
		Scrabbler: controllers.NewMockScrabbler(),
	}

	se.Init()
}
