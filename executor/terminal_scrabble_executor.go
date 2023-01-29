package executor

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Deuque/slack-scrabble/controllers"
)

type TerminalScrabbleExecutor struct {
	Word *string
	*controllers.Scrabbler
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

	writer := ScrabbleWriter{
		Write: func(reply string) error {
			fmt.Println(reply)
			return nil
		},
	}
	se.readAndHandleInput(&writer)
}

func (se *TerminalScrabbleExecutor) sessionStarted() bool {
	return se.Word != nil && len(*se.Word) > 0
}
func (se *TerminalScrabbleExecutor) onNewScrabbleCommand(writer *ScrabbleWriter) {
	word, err := se.FetchWord()
	if err != nil {
		writer.Write(err.Error())
		return
	}
	if word == nil || len(*word) == 0 {
		writer.Write("Error generating word")
		return

	}
	se.Word = word

	scrabbledWord := se.ScrabbleWord(*word)
	writer.Write(scrabbledWord)

}

func (se *TerminalScrabbleExecutor) onReShuffleCommand(writer *ScrabbleWriter) {
	if !se.sessionStarted() {
		writer.Write(NoSessionError)
		return
	}

	scrabbledWord := se.ScrabbleWord(*se.Word)
	writer.Write(scrabbledWord)

}

func (se *TerminalScrabbleExecutor) onAnswerScrabbleCommand(answer string, writer *ScrabbleWriter) {
	if !se.sessionStarted() {
		writer.Write(NoSessionError)
		return
	}

	correct := se.CheckAnswer(*se.Word, answer)
	if !correct {
		writer.Write("Ooops incorrect, lets try again!")
	} else {
		writer.Write("Correct!, You're the boss")
	}

}

func (se *TerminalScrabbleExecutor) onRevealAnswerCommand(writer *ScrabbleWriter) {
	if !se.sessionStarted() {
		writer.Write(NoSessionError)
		return
	}
	fmt.Printf("The answer is %s\n", *se.Word)
}

func (se *TerminalScrabbleExecutor) readAndHandleInput(writer *ScrabbleWriter) {
	writer.Write("")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		writer.Write("Error reading input")
	}

	text := scanner.Text()

	if strings.EqualFold(text, "scrab new") {
		se.onNewScrabbleCommand(writer)
		se.readAndHandleInput(writer)
	} else if strings.EqualFold(text, "scrab re") {
		se.onReShuffleCommand(writer)
		se.readAndHandleInput(writer)
	} else if strings.HasPrefix(text, "scrab ans") {
		split := strings.Split(text, " ")
		if len(split) < 3 {
			writer.Write("Attach the answer after the word \"ans\"")
		} else {
			se.onAnswerScrabbleCommand(split[2], writer)
		}
		se.readAndHandleInput(writer)
	} else if strings.EqualFold(text, "scrab tell") {
		se.onRevealAnswerCommand(writer)
		se.readAndHandleInput(writer)
	} else {
		writer.Write("Unknown command")
		se.readAndHandleInput(writer)
	}

}
