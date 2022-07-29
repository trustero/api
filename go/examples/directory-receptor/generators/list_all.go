package generators

import (
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"os/exec"
	"time"
)

type ListAll struct{}

func (m ListAll) Caption() string {
	return "trr-100000"
}

func (m ListAll) Report(credentials interface{}) (report []interface{}, sources []*cmd.Source, err error) {
	creds := credentials.(*Credentials)

	var output []byte
	if output, err = exec.Command("ls", "-al", creds.Path).Output(); err != nil {
		return
	}
	println(string(output))
	report = []interface{}{
		&Directory{
			Id:        "100000",
			CreatedAt: time.Now(),
			IsFolder:  true,
			Owner:     "Root",
		},
		&Directory{
			Id:        "100000",
			CreatedAt: time.Now(),
			IsFolder:  true,
			Owner:     "Root",
		},
		&Directory{
			Id:        "100000",
			CreatedAt: time.Now(),
			IsFolder:  true,
			Owner:     "Root",
		},
	}
	return
}
