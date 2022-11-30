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
func NewEvidence(serviceName, entityType, caption, description string) *Evidence {
	return &Evidence{
		ServiceName: serviceName,
		EntityType:  entityType,
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

type ServiceEntities receptor_v1.ServiceEntities

// NewServiceEntities instantiate a new service
func NewServiceEntities() *ServiceEntities {
	return (*ServiceEntities)(&receptor_v1.ServiceEntities{})
}

// AddService adds a services to discovered services
func (s *ServiceEntities) AddService(typeName, typeId, instanceName, instanceId string) *ServiceEntities {
	if len(typeName) > 0 {
		s.Entities = append(s.Entities, newService(typeName, typeId, instanceName, instanceId))
	}
	return s
}

// NewService is a helper to instantiate a new Source struct
func newService(serviceName, entityType, entityInstanceName, entityInstanceId string) *receptor_v1.ServiceEntity {
	return &receptor_v1.ServiceEntity{
		ServiceName:        serviceName,
		EntityType:         entityType,
		EntityInstanceName: entityInstanceName,
		EntityInstanceId:   entityInstanceId,
	}
}
