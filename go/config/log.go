package config

import (
	"encoding/json"
	"fmt"
	"github.com/natefinch/lumberjack"
	"io"
	"os"
	"path"
	"time"

	grpcZeroLog "github.com/cheapRoc/grpc-zerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/grpclog"
)

const (
	colorBold     = 1
	colorCyan     = 36
	colorDarkGray = 90
)

var tz *time.Location
var noColor bool

// InitLog Setup server logging using zerolog
func InitLog(levelStr string, logFile string) {

	// Use current timezone when printing console log messages
	tz = time.Now().Location()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// set log level
	level, _ := zerolog.ParseLevel(levelStr)
	if level == zerolog.NoLevel {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}

	var writers []io.Writer

	// create log writer
	writers = append(writers, zerolog.ConsoleWriter{
		Out:             os.Stderr,
		FormatTimestamp: consoleFormatTimestamp,
		FormatCaller:    consoleFormatCaller,
		NoColor:         true})
	noColor = true

	if logFile != "" {
		writers = append(writers, rollingLog(logFile, 3, 1, 1))
	}

	mw := io.MultiWriter(writers...)

	// set global logger to our setup
	log.Logger = zerolog.New(mw).With().Timestamp().Logger()
	log.Logger = log.Logger.With().Caller().Logger()

	// setup grpc logging
	grpcLog := log.Logger.With().CallerWithSkipFrameCount(7).Logger()
	grpclog.SetLoggerV2(grpcZeroLog.New(grpcLog.With().Str("workstation", "grpc").Logger()))
}

func rollingLog(filePath string, maxBackups, maxSize, maxAge int) io.Writer {
	folder := path.Dir(filePath)
	if err := os.MkdirAll(folder, 0744); err != nil {
		log.Error().Err(err).Str("path", folder).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   filePath,
		MaxBackups: maxBackups, // files
		MaxSize:    maxSize,    // megabytes
		MaxAge:     maxAge,     // days
	}
}

func consoleFormatCaller(i interface{}) string {
	var c string
	if cc, ok := i.(string); ok {
		c = cc
	}
	if len(c) > 0 {
		c = trimPath(c)
		c = colorize(c, colorBold, noColor) + colorize(" >", colorCyan, noColor)
	}
	return c
}

func trimPath(s string) string {
	j := 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			j++
			if j == 2 {
				if s[i] == '/' {
					return s[i+1:]
				}
				return s[i:]
			}
		}
	}
	if s[0] == '/' {
		return s[1:]
	}
	return s
}

func consoleFormatTimestamp(i interface{}) string {
	t := "<nil>"
	switch tt := i.(type) {
	case string:
		ts, err := time.Parse(zerolog.TimeFieldFormat, tt)
		if err != nil {
			t = tt
		} else {
			t = ts.Format(time.Kitchen)
		}
	case json.Number:
		i, err := tt.Int64()
		if err != nil {
			t = tt.String()
		} else {
			var sec, nsec int64 = i, 0
			switch zerolog.TimeFieldFormat {
			case zerolog.TimeFormatUnixMs:
				nsec = int64(time.Duration(i) * time.Millisecond)
				sec = 0
			case zerolog.TimeFormatUnixMicro:
				nsec = int64(time.Duration(i) * time.Microsecond)
				sec = 0
			}
			ts := time.Unix(sec, nsec).In(tz)
			t = ts.Format(time.Kitchen)
		}
	}
	if noColor {
		return t
	}
	return colorize(t, colorDarkGray, noColor)
}

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
