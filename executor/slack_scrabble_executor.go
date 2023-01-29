package executor

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Deuque/slack-scrabble/controllers"
	"github.com/shomali11/slacker"
)

type SlackScrabbleExecutor struct {
	Word *string
	*controllers.Scrabbler
	*slacker.Slacker
}

func (se *SlackScrabbleExecutor) Init() {
	se.Slacker = slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	se.handleBotInput()

	err := se.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func (se *SlackScrabbleExecutor) sessionStarted() bool {
	return se.Word != nil && len(*se.Word) > 0
}

func (se *SlackScrabbleExecutor) onNewScrabbleCommand(writer *ScrabbleWriter) {

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

func (se *SlackScrabbleExecutor) onReShuffleCommand(writer *ScrabbleWriter) {
	if !se.sessionStarted() {
		writer.Write(NoSessionError)
		return
	}

	scrabbledWord := se.ScrabbleWord(*se.Word)
	writer.Write(scrabbledWord)

}

func (se *SlackScrabbleExecutor) onAnswerScrabbleCommand(answer string, writer *ScrabbleWriter) {
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

func (se *SlackScrabbleExecutor) onRevealAnswerCommand(writer *ScrabbleWriter) {
	if !se.sessionStarted() {
		writer.Write(NoSessionError)
		return
	}
	writer.Write(fmt.Sprintf("The answer is %s\n", *se.Word))
}

func (se *SlackScrabbleExecutor) handleBotInput() {
	se.Command("new", &slacker.CommandDefinition{
		Description: "Start a scrabble session",
		Examples:    []string{"new"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			writer := botScrabbleWriter(response)
			se.onNewScrabbleCommand(writer)
		},
	})

	se.Command("shuffle", &slacker.CommandDefinition{
		Description: "Shuffle a scrabble word again",
		Examples:    []string{"shuffle"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			writer := botScrabbleWriter(response)
			se.onReShuffleCommand(writer)
		},
	})

	se.Command("answer {answer}", &slacker.CommandDefinition{
		Description: "Answer a scrabble",
		Examples:    []string{"answer gridlock"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			writer := botScrabbleWriter(response)
			answer := request.Param("answer")
			if len(answer) == 0 {
				writer.Write("Attach the answer after the word \"answer")
			} else {
				se.onAnswerScrabbleCommand(answer, writer)
			}
		},
	})

	se.Command("reveal", &slacker.CommandDefinition{
		Description: "Reveal a scrabble answer",
		Examples:    []string{"reveal"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			writer := botScrabbleWriter(response)
			se.onRevealAnswerCommand(writer)
		},
	})

}

func botScrabbleWriter(response slacker.ResponseWriter) *ScrabbleWriter {
	return &ScrabbleWriter{
		func(reply string) error {
			response.Reply(reply)
			return nil
		},
	}
}
