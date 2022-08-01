package generators

import (
	"github.com/trustero/api/go/receptor_v1"
	"os/exec"
	"time"
)

type ListAll struct{}

func (m ListAll) Verify(credentials string) (successful bool, err error) {
	successful = true
	return
}

func (m ListAll) ReceptorModelId() string {
	return "trr-100000"
}

func (m ListAll) ServiceModelId() string {
	return "trs-100000"
}

func (m ListAll) EvidenceModelId() string {
	return "tre-100000"
}

func (m ListAll) GenerateEvidence(credentials string) (ev []interface{}, sources []*receptor_v1.Evidence_Source, err error) {
	var creds Credentials
	if creds, err = NewCredentials(credentials); err != nil {
		return
	}
	var output []byte
	if output, err = exec.Command("ls", "-al", creds.Path).Output(); err != nil {
		return
	}
	println(string(output))
	ev = append(ev, &Directory{
		Id:        "100000",
		CreatedAt: time.Now(),
		IsFolder:  true,
		Owner:     "Root",
	})
	ev = append(ev, &Directory{
		Id:        "100000",
		CreatedAt: time.Now(),
		IsFolder:  true,
		Owner:     "Root",
	})
	ev = append(ev, &Directory{
		Id:        "100000",
		CreatedAt: time.Now(),
		IsFolder:  true,
		Owner:     "Root",
	})
	ev = append(ev, &Directory{
		Id:        "100000",
		CreatedAt: time.Now(),
		IsFolder:  true,
		Owner:     "Root",
	})
	ev = append(ev, &Directory{
		Id:        "100000",
		CreatedAt: time.Now(),
		IsFolder:  true,
		Owner:     "Root",
	})
	return
}
