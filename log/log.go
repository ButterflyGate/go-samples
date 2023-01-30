package log

import (
	"github.com/ButterflyGate/logger"
	"github.com/ButterflyGate/logger/levels"
)

func Log() {
	o := logger.DefaultOutputOption().HideCursor().HideLevel()
	l := logger.NewLogger(
		levels.Trace,
		o,
	)
	l.Trace("hello,world")
}
