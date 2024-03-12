package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

var logger Logger

func NewLogger() Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatErrFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s: ", i)
		},
	}

	return Logger{
		zerolog.New(output).
			With().
			Caller().
			Timestamp().
			Logger(),
	}
}
