// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package receptor_sdk provides a CLI platform simplifying the development of Trustero Receptor(s).  A Receptor
// is a CLI that implements a contract of CLI arguments and flags the Trustero services invokes to collect
// evidence of how a business uses a service provider's services.
package receptor_sdk

import (
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	ConfigBase64URL      string // Receptor configuration as a base64 URL encoded json string.
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
	//  - method to which this field belongs when receptors support multiple auth methods
	//  - input_type is the html element input type. ex. text, password
	//
	// For example:
	//
	//  type Credentials struct {
	//      GroupId string `trustero:"display:Group Identifier;placeholder:abcdefg123"`
	//      Token   string `trustero:"display:Access Token;placeholder:1234wxyz"`
	//  }
	// The metadata can be extracted and then printed out using the following command
	// <receptor_type> descriptor
	GetCredentialObj() (credentialObj interface{})

	// GetConfigObj returns an instance of a Config object
	GetConfigObj() (configObj interface{})

	// GetConfigObjDesc returns an instance of struct that represents a json for the config object to be rendered
	// in the receptor config modal
	// To print what the config json will look like, use the following command
	// <receptor_type> config
	GetConfigObjDesc() (configObjDesc interface{})

	// GetAuthMethods returns an instance of struct representing the authentication methods supported by the
	// receptor
	GetAuthMethods() (authMethods interface{})

	//GetEvidenceInfo returns a list of Evidences that a receptor has implemented
	//The metadata is extracted and then printed out
	//<receptor_type> evidenceinfo
	GetEvidenceInfo(credentials interface{}) (evidences []*Evidence)

	// Verify read-only access to a service provider account.  Return ok if the credentials are valid and err
	// if any error is encountered in contacting the service provider.  This method is invoked from the following
	// ClI:  <receptor_type> verify
	Verify(credentials interface{}, config interface{}) (ok bool, err error)

	// Discover in-use service entities in a service provider.  Return an array of [receptor_v1.ServiceEntity]
	// discovered and err if any error is encountered.  This method is invoked from the following CLI:
	// <receptor_type> scan
	Discover(credentials interface{}, config interface{}) (services []*receptor_v1.ServiceEntity, err error)

	// Report in-use service entity's configurations as evidence.  Return an array of [Evidence] found and an error
	// if any error is encountered in contacting the service provider.  This method is invoked from the following
	// CLI: <receptor_type> scan --find-evidence
	Report(credentials interface{}, config interface{}) (evidences []*Evidence, err error)

	// ReportBatch reports in-use service entity's configurations as evidence in batches.
	// The receptor implementation should send the evidences to the evidenceChan.
	// CLI: <receptor_type> scan --find-evidence
	ReportBatch(credentials interface{}, evidenceChan chan []*Evidence)

	// Configure returns a ReceptorConfiguration object that represents the configuration of the receptor
	// Configure is used when there special configurations required for the receptor that the user can set
	Configure(credentials interface{}) (config *receptor_v1.ReceptorConfiguration, err error)

	// GetLogo returns the content of the logo in svg format for the receptor
	// This method is invoked from the following CLI: <receptor_type> logo
	GetLogo() (logo string, err error)

	// GetInstructions returns the instructions in markdown format for settings up the providers during receptor activation
	// This method is invoked from the following CLI: <receptor_type> instructions
	GetInstructions() (instructions string, err error)
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
	ServiceName        string                         // ServiceName where this evidence was gathered. For example, "S3".
	EntityType         string                         // EntityType of rows of evidence.  For example, "bucket".
	Caption            string                         // Caption identifies the evidence.
	Description        string                         // Description provides additional information on origins of the evidence.
	Sources            []*receptor_v1.Source          // Sources of raw API request and response used to gather the evidence.
	Rows               []interface{}                  // Rows of formatted evidence represented by a Golang struct.
	ServiceAccountId   string                         // AccountId of multi-account organization
	Document           *Document                      // Unstructured evidence in a Document format
	Controls           []string                       // Controls associated with the evidence
	IsManual           bool                           // If true, the evidence was manually collected
	RelevantDate       timestamppb.Timestamp          // Relevant date of the evidence
	EvidenceObjectType receptor_v1.EvidenceObjectType // Type of the evidence object - enum of receptor_v1.EvidenceObjectType

}

// Document is a unstructured byte array that can be used to store any type of data
// in Content field and Mime describes the Content type
type Document struct {
	Body           []byte // Content of the document in bytes
	Mime           string // Mime type of the document
	StreamFilePath string // Path to the file containing the evidence
}

// Config with Field struct defines the json shape of the custom configurations for receptors that the app can use
// for special condition receptors that might collect and post to an arbitrary number of controls
// that an user can map
type Config struct {
	Title       string  `json:"title"`       // Title of the modal for the config
	Description string  `json:"description"` // Description of the config
	Fields      []Field `json:"fields"`      // Fields in the config
}

type Field struct {
	Display         string      `json:"display"`           // Label for the field
	Placeholder     string      `json:"placeholder"`       // Placeholder for the field
	InputType       string      `json:"input_type"`        // Input Html Element type. For example, "Text", "Select"
	Field           string      `json:"field"`             // name of the field
	Options         interface{} `json:"options,omitempty"` // name-value pairs when InputType is a select from list
	EvidenceCaption string      `json:"evidence_caption"`  // caption of the evidence
	ServiceModelID  string      `json:"service_model_id"`  // trustero model id for the service
}

// AuthodMethod struct to list the authentication methods supported
// by the receptor

type AuthMethod struct {
	Display string `json:"display"` // Label for the authentication method
	Value   string `json:"value"`   // Id/Name for the authentication method
}

type Control struct {
	Id               string             `json:"id"`                 // Id of the control
	Name             string             `json:"name"`               // Name of the control
	Objective        string             `json:"objective"`          // Objective of the control
	TestProcedure    string             `json:"test_procedure"`     // Test procedure of the control
	Notes            string             `json:"notes"`              // Notes for the control
	RequiredEvidence string             `json:"required_evidences"` // Required evidences for the control
	ExternalId       string             `json:"external_id"`        // External id of the control
	ExternalLink     string             `json:"external_link"`      // External link to the control
	Procedures       []ControlProcedure `json:"control_procedures"` // Procedures for the control
}

type Policy struct {
	Id           string `json:"id"`            // Id of the policy
	Name         string `json:"name"`          // Name of the policy
	Description  string `json:"description"`   // Description of the policy
	Department   string `json:"department"`    // Department of the policy
	ExternalId   string `json:"external_id"`   // External id of the policy
	ExternalLink string `json:"external_link"` // External link to the policy
}

type ControlPolicyMapping struct {
	ControlId string `json:"control_id"` // Id of the control
	PolicyId  string `json:"policy_id"`  // Id of the policy
}

type ControlProcedure struct {
	Id                string `json:"id"`
	ControlName       string `json:"Control_Name"`
	Description       string `json:"Description"`
	ProcedureID       string `json:"Procedure_ID"`
	ProcedureName     string `json:"Procedure_Name"`
	TestingProcedures string `json:"Testing_Procedures"`
	ExternalId        string `json:"external_id"`
	ExternalLink      string `json:"external_link"`
}
