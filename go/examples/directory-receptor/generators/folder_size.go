package generators

import (
	"os/exec"

	"github.com/trustero/monorepo/receptor-sdk/client/api/model"
)

type FolderSize struct{}

func (m FolderSize) ReceptorModelId() string {
	return "trr-100000"
}

func (m FolderSize) ServiceModelId() string {
	return "trs-100000"
}

func (m FolderSize) EvidenceModelId() string {
	return "tre-100002"
}

func (m FolderSize) GenerateEvidence(credentials string) (ev []interface{}, sources []*model.AddEvidenceRequest_Source, err error) {
	var creds Credentials
	if creds, err = NewCredentials(credentials); err != nil {
		return
	}
	var output []byte
	if output, err = exec.Command("du", "-h", "-d 1", creds.Path).Output(); err != nil {
		return
	}
	println(string(output))
	return
}
