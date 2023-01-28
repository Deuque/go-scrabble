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

var (
	NoSessionError = "You have not started a session, enter \"scrab new\" to begin"
)

func (se *TerminalScrabbleExecutor) Init() {
	fmt.Println("Hello, Welcome to terminal scrabble, use any of the following commands:" +
		"\n  scrab new : to begin a new session" +
		"\n  srab re : to rearrange the word" +
		"\n  srab ans  <yourAnswer> : to check your answer" +
		"\n  srab tell : to reveal the answer" +
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
		return
	}
	if word == nil || len(*word) == 0 {
		fmt.Println("Error generating word")
		return

	}
	se.Word = *word

	scrabbledWord := se.Scrabbler.ScrabbleWord(*word)
	fmt.Println(scrabbledWord)

}

func (se *TerminalScrabbleExecutor) OnReScrabbleCommand() {
	if !se.SessionStarted() {
		fmt.Println(NoSessionError)
		return
	}

	scrabbledWord := se.Scrabbler.ScrabbleWord(se.Word)
	fmt.Println(scrabbledWord)

}

func (se *TerminalScrabbleExecutor) OnAnswerScrabbleCommand(answer string) {
	if !se.SessionStarted() {
		fmt.Println(NoSessionError)
		return
	}

	correct := se.Scrabbler.CheckAnswer(se.Word, answer)
	if !correct {
		fmt.Println("Ooops incorrect, lets try again!")
	} else {
		fmt.Println("Correct!, You're the boss")
	}

}

func (se *TerminalScrabbleExecutor) OnRevealAnswerCommand() {
	if !se.SessionStarted() {
		fmt.Println(NoSessionError)
		return
	}
	fmt.Printf("The answer is %s\n", se.Word)
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
	} else if strings.EqualFold(text, "scrab tell") {
		se.OnRevealAnswerCommand()
		se.readAndHandleInput()
	} else {
		fmt.Println("Unknown command")
		se.readAndHandleInput()
	}

}
