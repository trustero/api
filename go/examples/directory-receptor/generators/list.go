package generators

import (
	"github.com/trustero/monorepo/receptor-sdk/client/api/model"
	"os/exec"
)

type List struct{}

func (m List) ReceptorModelId() string {
	return "trr-100000"
}

func (m List) ServiceModelId() string {
	return "trs-100000"
}

func (m List) EvidenceModelId() string {
	return "tre-100001"
}

func (m List) GenerateEvidence(credentials string) (ev []interface{}, sources []*model.AddEvidenceRequest_Source, err error) {
	var creds Credentials
	if creds, err = NewCredentials(credentials); err != nil {
		return
	}
	var output []byte
	if output, err = exec.Command("ls", "-l", creds.Path).Output(); err != nil {
		return
	}
	println(string(output))
	return
}
