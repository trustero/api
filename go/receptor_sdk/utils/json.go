package utils

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func Json2Map(jsonStr string) (strMap map[string]string, err error) {
	if err = json.Unmarshal([]byte(jsonStr), &strMap); err != nil {
		log.Error().Err(err).Msg("Failed to parse json")
		return
	}
	return
}
