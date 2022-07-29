package cmd

import "github.com/trustero/api/go/receptor_v1"

type Service = receptor_v1.Services_Service
type Source = receptor_v1.Evidence_Source

// Receptor is the main interface for the Receptor implementor-facing  API.
type Receptor interface {
	// ReceptorModelId returns the model id of the receptor.
	ReceptorModelId() string

	// ServiceModelId returns the model id of the service.
	ServiceModelId() string

	// UnmarshallCredentials deserializes the credentials json string and returns the result as a struct pointer.
	UnmarshallCredentials(credentials string) (result interface{}, err error)

	// Verify verifies the given credentials. Called by the cli on
	// <receptor> verify-
	Verify(serviceCredentials interface{}) (ok bool, err error)

	// Discover returns the list of services that can be used to generate evidence. Called by the cli on
	// <receptor> scan
	Discover(serviceCredentials interface{}) (services []*Service, err error)

	// GetReporters returns the list of reporters that can be used to generate evidence. Each one wil be called
	// individually on
	//<receptor> scan --find-evidence
	GetReporters() []Reporter

	// CredentialsFromFlags Parses crentials provided via custom comandline flags as a struct
	CredentialsFromFlags() interface{}
}

// Reporter is the interface for an evidence reporter. Will be called by the cli on
// <receptor> scan --find-evidence.
type Reporter interface {
	// Caption returns the caption of the evidence.
	Caption() string

	// Report returns the evidence and the sources.
	Report(interface{}) (report []interface{}, sources []*Source, err error)
}
