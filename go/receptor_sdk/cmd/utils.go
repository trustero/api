package cmd

import (
	"github.com/rs/zerolog"
)

func onDebug(credentials interface{}, call func(interface{}) error) (err error) {
	NoSave = true
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return call(credentials)

}
