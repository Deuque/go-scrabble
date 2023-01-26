package executor

type ScrabbleExecutor interface {
	Init()
	SessionStarted() bool
	OnNewScrabbleCommand()
	OnReScrabbleCommand()
	OnAnswerScrabbleCommand(answer string)
}
