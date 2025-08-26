package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var LogLevel = map[string]zerolog.Level{
	"TRACE":    -1,
	"DEBUG":    0,
	"INFO":     1,
	"WARN":     2,
	"ERROR":    3,
	"FATAL":    4,
	"PANIC":    5,
	"DISABLED": 7,
}

func SetGlobalLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

func New(data map[string]any) {
	zerologContext := zerolog.New(os.Stdout).With()
	for k, v := range data {
		switch v := v.(type) {
		case string:
			zerologContext = zerologContext.Str(k, v)
		case int:
			zerologContext = zerologContext.Int(k, v)
		case float64:
			zerologContext = zerologContext.Float64(k, v)
		case bool:
			zerologContext = zerologContext.Bool(k, v)
		case time.Time:
			zerologContext = zerologContext.Time(k, v)
		}
	}

	log.Logger = zerologContext.Timestamp().Logger()
}

func Debug(message string) {
	log.Debug().Msg(message)
}

func Info(message string) {
	log.Info().Msg(message)
}

func Warn(message string) {
	log.Warn().Msg(message)
}

func Error(message string) {
	log.Error().Msg(message)
}

func Debugf(message string, args ...any) {
	log.Debug().Msgf(message, args...)
}

func Infof(message string, args ...any) {
	log.Info().Msgf(message, args...)
}

func Warnf(message string, args ...any) {
	log.Warn().Msgf(message, args...)
}

func Errorf(message string, args ...any) {
	log.Error().Msgf(message, args...)
}
