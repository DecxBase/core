package logger

import "github.com/phuslu/log"

var Main log.Logger

func init() {
	Main = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
}

func Create(pairs ...string) log.Logger {
	return CreateFrom(Main, pairs...)
}

func CreateFrom(l log.Logger, pairs ...string) log.Logger {
	vlen := len(pairs)
	loggerContext := log.NewContext(l.Context)

	for i := 0; i < len(pairs); i += 2 {
		if vlen > (i + 1) {
			loggerContext = loggerContext.Str(pairs[i], pairs[i+1])
		}
	}

	subLogger := l
	subLogger.Context = loggerContext.Value()

	return subLogger
}
