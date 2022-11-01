package receptorLog

import (
	"github.com/rs/zerolog/log"
)

func Trace(format string, v ...interface{}) {
	log.Trace().Msgf(format, v...)
}

func Debug(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Info(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

func Error(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func Panic(format string, v ...interface{}) {
	log.Panic().Msgf(format, v...)
}
