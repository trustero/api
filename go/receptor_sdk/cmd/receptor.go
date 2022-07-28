package cmd

type Service = Services_Service
type Source = Evidence_Source

// Receptor is the main interface for the Receptor implementor-facing  API.
type Receptor interface {
	// Verify verifies the given credentials. Called by the CLI on
	// <receptor> verify-
	Verify(serviceCredentials map[string]string) (ok bool, err error)

	// Discover returns the list of services that can be used to generate evidence. Called by the CLI on
	// <receptor> scan
	Discover(serviceCredentials map[string]string) (services []*Service, err error)

	// GetReporters returns the list of reporters that can be used to generate evidence. Each one wil be called
	// individually on
	//<receptor> scan --find-evidence
	GetReporters() []Reporter
}

// Reporter is the interface for an evidence reporter. Will be called by the CLI on
// <receptor> scan --find-evidence.
type Reporter interface {
	// ReceptorModelId returns the model id of the receptor.
	ReceptorModelId() string

	// ServiceModelId returns the model id of the service.
	ServiceModelId() string

	// Caption returns the caption of the evidence.
	Caption() string

	// Report returns the evidence and the sources.
	Report(map[string]string) (report []interface{}, sources []*Source, err error)
}
