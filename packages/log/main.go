package log

import (
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)

	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
}

type LogLevel = string

const (
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
)

type logger struct {
	LogLevel LogLevel
}

func NewLogger(level LogLevel) Logger {
	if level < DEBUG || level > INFO {
		level = DEBUG
	}

	return &logger{
		LogLevel: level,
	}
}

func (l *logger) Debug(message string) {
	if l.LogLevel <= DEBUG {
		log.Debug().Msg(message)
	}
}

func (l *logger) Info(message string) {
	if l.LogLevel <= INFO {
		log.Info().Msg(message)
	}
}

func (l *logger) Warn(message string) {
	if l.LogLevel <= WARN {
		log.Warn().Msg(message)
	}
}

func (l *logger) Error(message string) {
	if l.LogLevel <= ERROR {
		log.Error().Msg(message)
	}
}

func (l *logger) Debugf(message string, args ...interface{}) {
	if l.LogLevel <= DEBUG {
		log.Debug().Msgf(message, args...)
	}
}
func (l *logger) Infof(message string, args ...interface{}) {
	if l.LogLevel <= INFO {
		log.Info().Msgf(message, args...)
	}
}
func (l *logger) Warnf(message string, args ...interface{}) {
	if l.LogLevel <= WARN {
		log.Warn().Msgf(message, args...)
	}
}
func (l *logger) Errorf(message string, args ...interface{}) {
	if l.LogLevel <= ERROR {
		log.Error().Msgf(message, args...)
	}
}
