package executor

type ScrabbleExecutor interface {
	Init()
	sessionStarted() bool
	onNewScrabbleCommand(*ScrabbleWriter)
	onReShuffleCommand(*ScrabbleWriter)
	onAnswerScrabbleCommand(answer string, writer *ScrabbleWriter)
	onRevealAnswerCommand(*ScrabbleWriter)
}

type ScrabbleWriter struct {
	Write func(string) error
}
