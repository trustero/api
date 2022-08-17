// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package receptor_sdk provides a CLI platform simplifying the development of Trustero Receptor(s).  A Receptor
// is a CLI that implements a contract of CLI arguments and flags the Trustero services invokes to collect
// evidence of how a business uses a service provider's services.
package receptor_sdk

import (
	"github.com/trustero/api/go/receptor_v1"
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

	// GetKnownServices returns a list of service names this receptor will be providing evidence for.
	GetKnownServices() (serviceNames []string)

	// GetCredentialObj returns an instance of a credential struct used for service provider authentication.  The
	// credential struct contains only public string fields with Go struct field tags:
	//
	// Field tag name is 'trustero' with sub-tags separated by ';' and valid sub-tags are 'display', and 'placeholder'
	//  - display provides the human-readable name of the field
	//  - placeholder provides a default field value suggestion for the field
	//
	// For example:
	//
	//  type Credentials struct {
	//      GroupId string `trustero:"display:Group Identifier;placeholder:abcdefg123"`
	//      Token   string `trustero:"display:Access Token;placeholder:1234wxyz"`
	//  }
	GetCredentialObj() (credentialObj interface{})

	// Verify read-only access to a service provider account.  Return ok if the credentials are valid and err
	// if any error is encountered in contacting the service provider.  This method is invoked from the following
	// ClI:  <receptor_type> verify
	Verify(credentials interface{}) (ok bool, err error)

	// Discover in-use service entities in a service provider.  Return an array of [receptor_v1.ServiceEntity]
	// discovered and err if any error is encountered.  This method is invoked from the following CLI:
	// <receptor_type> scan
	Discover(credentials interface{}) (services []*receptor_v1.ServiceEntity, err error)

	// Report in-use service entity's configurations as evidence.  Return an array of [Evidence] found and an error
	// if any error is encountered in contacting the service provider.  This method is invoked from the following
	// CLI: <receptor_type> scan --find-evidence
	Report(credentials interface{}) (evidences []*Evidence, err error)
}

// Evidence is a discovered evidence from an in-use service.  All rows in the evidence are instances of the same
// Golang struct.  Fields of this evidence row struct must be public and annotated with Trustero's field annotation
// where:
//
// Field tag name is 'trustero' with sub-tags separated by ';' and valid sub-tags are 'id', 'display', and 'order'
//   - id specifies the field is unique identifier for the struct.
//   - display provides the human-readable name for the field.
//   - order is an integer number starting with 1, denoting the order in which the field should be displayed in a
//     table.
//
// For example:
//
//	type User struct {
//	    Name     string  `trustero:"display:Name;order:2"`
//	    IsAdmin  bool    `trustero:"display:Admin;order:3"`
//	    Username string  `trustero:"id;display:User Name;order:1"`
//	}
type Evidence struct {
	ServiceName string                // ServiceName where this evidence wa gathered. For example, "S3".
	EntityType  string                // EntityType of rows of evidence.  For example, "bucket".
	Caption     string                // Caption identifies the evidence.
	Description string                // Description provides additional information on origins of the evidence.
	Sources     []*receptor_v1.Source // Sources of raw API request and response used to gather the evidence.
	Rows        []interface{}         // Rows of formatted evidence represented by a Golang struct.
}
