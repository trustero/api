package main

import (
	"github.com/trustero/api/go/examples/directory-receptor/generators"
	"github.com/trustero/api/go/receptor_sdk/cmd"
)

var (
	Path     string
	receptor = &Receptor{}
)

type Receptor struct{}

func (m Receptor) UnmarshallCredentials(credentials string) (result interface{}, err error) {
	return &generators.Credentials{}, nil
}

func (m Receptor) ReceptorModelId() string {
	return "trr-100000"
}

func (m Receptor) ServiceModelId() string {
	return "trs-100000"
}

func (m Receptor) Verify(credentials interface{}) (ok bool, err error) {
	return
}

func (m Receptor) CredentialsFromFlags() interface{} {
	if len(Path) > 0 {
		return &generators.Credentials{
			Path: Path,
		}
	}
	return nil
}

// Discover returns the list of services that can be used to generate evidence. Called by the cli on
// <receptor> scan
func (m Receptor) Discover(serviceCredentials interface{}) (services []*cmd.Service, err error) {
	return
}

// GetReporters returns the list of reporters that can be used to generate evidence. Each one wil be called
// individually on
//<receptor> scan --find-evidence
func (m Receptor) GetReporters() (reporters []cmd.Reporter) {
	reporters = append(reporters, &generators.ListAll{})
	return
}
