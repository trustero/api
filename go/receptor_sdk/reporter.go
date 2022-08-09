// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package receptor_sdk

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/trustero/api/go/receptor_v1"
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

// NewEvidence is a helper to instantiate a new Evidence struct.
func NewEvidence(serviceName, caption, description string) *Evidence {
	return &Evidence{
		ServiceName: serviceName,
		Caption:     caption,
		Description: description,
		Sources:     []*receptor_v1.Source{},
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

	ev.Sources = append(ev.Sources, &receptor_v1.Source{
		RawApiResponse: strResponse,
		RawApiRequest:  strRequest,
	})
	return ev
}

// AddRow appends the given Evidence row.
func (ev *Evidence) AddRow(row interface{}) *Evidence {
	ev.Rows = append(ev.Rows, row)
	return ev
}

type Services receptor_v1.Services

// NewServices instantiate a new service
func NewServices() *Services {
	return (*Services)(&receptor_v1.Services{})
}

// AddService adds a services to discovered services
func (s *Services) AddService(typeName, typeId, instanceName, instanceId string) *Services {
	if len(typeName) > 0 && len(typeId) > 0 && len(instanceName) > 0 && len(instanceId) > 0 {
		s.Services = append(s.Services, newService(typeName, typeId, instanceName, instanceId))
	}
	return s
}

// NewService is a helper to instantiate a new Source struct
func newService(typeId, subtypeName, instanceName, instanceId string) *receptor_v1.Service {
	return &receptor_v1.Service{
		TypeId:       typeId,
		SubtypeName:  subtypeName,
		InstanceName: instanceName,
		InstanceId:   instanceId,
	}
}
