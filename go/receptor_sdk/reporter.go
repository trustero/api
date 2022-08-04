// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package receptor_sdk

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

// Report struct to hold collected evidence
type Report struct {
	Evidences []*Evidence
}

// NewReport is a helper to instantiate a Report struct
func NewReport() *Report {
	return &Report{}
}

// AddEvidence adds an evidence to report
func (r *Report) AddEvidence(evidence *Evidence) *Report {
	if evidence != nil {
		r.Evidences = append(r.Evidences, evidence)
	}
	return r
}

// NewEvidence is a helper to instantiate a new Evidence struct
func NewEvidence(serviceName, caption, description string) *Evidence {
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
func (ev *Evidence) AddSource(rawRequest, rawResponse interface{}) *Evidence {
	var strRequest, strResponse string

	// marshal raw request
	if str, ok := rawRequest.(string); ok {
		strRequest = str
	} else {
		if bytes, err := json.MarshalIndent(rawRequest, "", "  "); err == nil {
			strRequest = string(bytes)
		} else {
			log.Err(err).Msg("error marshalling API call request")
		}
	}

	// marshal raw response
	if str, ok := rawResponse.(string); ok {
		strResponse = str
	} else {
		if bytes, err := json.MarshalIndent(rawResponse, "", "  "); err == nil {
			strResponse = string(bytes)
		} else {
			log.Err(err).Msg("error marshalling API call response")
		}
	}

	ev.Sources = append(ev.Sources, &Source{
		ProviderAPIResponse: strResponse,
		ProviderAPIRequest:  strRequest,
	})
	return ev
}

// AddRow appends the given Evidence row
func (ev *Evidence) AddRow(row interface{}) *Evidence {
	ev.Rows = append(ev.Rows, row)
	return ev
}

// Services struct to hold collected service
type Services struct {
	Services []*Service
}

// NewServices instantiate a new service
func NewServices() *Services {
	return &Services{}
}

// AddService adds a services to discovered services
func (s *Services) AddService(name, id string) *Services {
	if len(name) > 0 && len(id) > 0 {
		s.Services = append(s.Services, NewService(name, id))
	}
	return s
}

// NewService is a helper to instantiate a new Source struct
func NewService(name, id string) *Service {
	return &Service{
		Name:       name,
		InstanceId: id,
	}
}
