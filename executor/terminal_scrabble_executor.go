package executor

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Deuque/slack-scrabble/controllers"
)

type TerminalScrabbleExecutor struct {
	Word      string
	Scrabbler controllers.Scrabbler
}

func (se *TerminalScrabbleExecutor) Init() {
	fmt.Println("Hello, Welcome to terminal scrabble, use any of the following commands:" +
		"\n  scrab new  : to begin a new session" +
		"\n  srab re : to rearrange the word" +
		"\n  srab ans  <yourAnswer> : to check your result" +
		"\nLet's go!")

	se.readAndHandleInput()
}

func (se *TerminalScrabbleExecutor) SessionStarted() bool {
	return len(se.Word) > 0
}
func (se *TerminalScrabbleExecutor) OnNewScrabbleCommand() {
	word, err := se.Scrabbler.FetchWord()
	if err != nil {
		fmt.Println(err)
	} else if word == nil || len(*word) == 0 {
		fmt.Println("Error generating word")
	} else {
		se.Word = *word

		scrabbledWord := se.Scrabbler.ScrabbleWord(*word)
		fmt.Println(scrabbledWord)
	}
}

func (se *TerminalScrabbleExecutor) OnReScrabbleCommand() {
	if !se.SessionStarted() {
		fmt.Println("You have not started a session, enter \"scrab new\" to begin")
	} else {
		scrabbledWord := se.Scrabbler.ScrabbleWord(se.Word)
		fmt.Println(scrabbledWord)
	}
}

func (se *TerminalScrabbleExecutor) OnAnswerScrabbleCommand(answer string) {
	if !se.SessionStarted() {
		fmt.Println("You have not started a session, enter \"scrab new\" to begin")
	} else {
		correct := se.Scrabbler.CheckAnswer(se.Word, answer)
		if !correct {
			fmt.Println("Ooops incorrect, lets try again!")
		} else {
			fmt.Println("Correct!, You're the boss")
		}
	}
}

func (se *TerminalScrabbleExecutor) readAndHandleInput() {
	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		fmt.Println("Error reading input")
	}

	text := scanner.Text()

	if strings.EqualFold(text, "scrab new") {
		se.OnNewScrabbleCommand()
		se.readAndHandleInput()
	} else if strings.EqualFold(text, "scrab re") {
		se.OnReScrabbleCommand()
		se.readAndHandleInput()
	} else if strings.HasPrefix(text, "scrab ans") {
		split := strings.Split(text, " ")
		if len(split) < 3 {
			fmt.Println("Attach the answer after the word \"ans")
		} else {
			se.OnAnswerScrabbleCommand(split[2])
		}
		se.readAndHandleInput()
	} else {
		fmt.Println("Unknown command")
		se.readAndHandleInput()
	}

}
