// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package receptor_sdk

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

// NewService is a helper to instantiate a new Source struct
func NewService(name, id string) (source *Service) {
	return &Service{
		Name:       name,
		InstanceId: id,
	}
}

// NewEvidence is a helper to instantiate a new Evidence struct
func NewEvidence(serviceName, caption, description string) (evidence *Evidence) {
	return &Evidence{
		Caption:     caption,
		Description: description,
		ServiceName: serviceName,
		Sources:     []*Source{},
		Rows:        []interface{}{},
	}
}

// AddSource appends raw API request and response to an Evidence struct.  AddSource will json marshal non-string
// request and response.
func (ev *Evidence) AddSource(rawRequest, rawResponse interface{}) {
	var strRequest, strResponse string

	// marshal raw request
	if str, ok := rawRequest.(string); ok {
		strRequest = str
	} else {
		if bytes, err := json.Marshal(rawRequest); err == nil {
			strRequest = string(bytes)
		} else {
			log.Err(err).Msg("error marshalling API call request")
		}
	}

	// marshal raw response
	if str, ok := rawResponse.(string); ok {
		strResponse = str
	} else {
		if bytes, err := json.Marshal(rawResponse); err == nil {
			strResponse = string(bytes)
		} else {
			log.Err(err).Msg("error marshalling API call response")
		}
	}

	ev.Sources = append(ev.Sources, &Source{
		ProviderAPIResponse: strResponse,
		ProviderAPIRequest:  strRequest,
	})
}

// AddRow appends the given Evidence row
func (ev *Evidence) AddRow(row interface{}) {
	ev.Rows = append(ev.Rows, row)
}
