package generators

import "encoding/json"

type Credentials struct {
	Path string `json:"path"`
}

func NewCredentials(rawCredentials string) (credentials Credentials, err error) {
	if err = json.Unmarshal([]byte(rawCredentials), &credentials); err != nil {
		return credentials, err
	}
	return credentials, err
}
