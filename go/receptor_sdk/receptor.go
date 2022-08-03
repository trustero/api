package receptor_sdk

import (
	"encoding/json"
)

// Global variables available to the receptor
var (
	Host                 string // Trustero GRPC host name, typically "prod.api.infra.trustero.com"
	Port                 int    // Trustero GRPC host port, typically 8443
	Cert                 string // HTTPS public certificate for GRPC host
	CertServerOverride   string // Do not verify remote Trustero GRPC hostname against the host set in HTTPS certificate.
	LogLevel             string // Log level.  From least to most verbose: panic, fatal, error, warn, info, debug, trace.
	LogFile              string // Logfile path
	ModelID              string // Receptor type name.  A receptor must set this with the SetReceptorType function.
	NoSave               bool   // If true, do not contact Trustero with results from the command.
	Notify               string // Trustero will provide a string tracer ID when it's tracing a receptor execution path.
	FindEvidence         bool   // If true as part of a scan command, scan for evidence in a service provider account.
	CredentialsBase64URL string // Service provider credentials as a base64 URL encoded json string.
	ReceptorId           string // Trustero's persistent record ID of a record holding a receptor's service provider credentials.
)

// Receptor is the main interface for the Receptor implementor-facing  API.
type Receptor interface {
	// GetReceptorType returns the receptor's type.  A receptor is expected to report findings from only one service
	// provider type.  Receptor type is a stable identifier that represent the type of receptor reporting this finding.
	//  The identifier is a simple URL encoded string that includes an organization name and the service provider name.
	//  For example: "trustero_gitlab".
	GetReceptorType() (receptorType string)

	// UnmarshalCredentials converts the service provider account credential json string into a Go object used by
	// the receptor's Verify, Discover, and Report methods.
	UnmarshalCredentials(credentials string) (obj interface{}, err error)

	// Verify read-only access to a service provider account.  This method is invoked from the following ClI:
	// <receptor> verify
	Verify(credentials interface{}) (ok bool, err error)

	// Discover returns the list of services that can be used to generate evidence. This method is invoked from
	// the following CLI:
	// <receptor> scan
	Discover(credentials interface{}) (services []*Service, err error)

	// Report returns the list of discovered evidence.  This method is invoked from the following CLI:
	// <receptor> scan --find-evidence
	Report(credentials interface{}) (evidences []*Evidence, err error)
}

// Service is a discovered in-use service in a service provider account.
type Service struct {
	Name       string // Service name, for example "S3".
	InstanceId string // Service identifier, for example S3 bucket name.
}

// Evidence is a discovered evidence from an in-use service.
type Evidence struct {
	ServiceName string        // Service name where this evidence wa gathered. For example, "S3".
	Caption     string        // Caption describing the evidence.  This string must be stable this type of evidence.
	Description string        // Description of what this evidence is.
	Sources     []*Source     // Source API requests and responses used to generate the evidence.
	Rows        []interface{} // REMIND.  Need to describe how to annotate an evidence struct.
}

// Source contain the raw API request and response from a service provider used to generate evidence.
type Source struct {
	ProviderAPIRequest  string // API request
	ProviderAPIResponse string // API response
}

// CredentialsFromFlagsFunc
// A utility function to convert receptor specific CLI arguments to a credentials string.  A credentials string
// is a json object as a string.  See CredentialsFromFlags variable for more background.
type CredentialsFromFlagsFunc func() string

// CredentialsFromFlags
// A receptor CLI can add custom CLI flags as credentials instead of using the base64 URL encoded --credentials
// flag or getting the credentials stored in Trustero.  When a CLI chooses to add custom credential flags, it must
// implement the CredentialsFromFlagsFunc to return credentials as a json object as a string.
var CredentialsFromFlags CredentialsFromFlagsFunc = func() string { return "" }

// UnmarshalCredentials
// A utilities function to unmarshal a json string into the provided object type.
func UnmarshalCredentials(credentials string, credentialsObj interface{}) (obj interface{}, err error) {
	err = json.Unmarshal([]byte(credentials), credentialsObj)
	obj = credentialsObj
	return
}
