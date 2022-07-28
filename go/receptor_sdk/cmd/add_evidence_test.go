package cmd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustero/monorepo/receptor-sdk/client/api/model"
	"github.com/trustero/monorepo/receptor-sdk/cmd"
	"github.com/trustero/monorepo/receptor-sdk/pkg"
	"testing"
	"time"
)

func TestAddEvidenceAPI(t *testing.T) {
	cmd.EvidenceGenerators = []pkg.EvidenceGenerator{&myEvidenceFinder{}}
	cmd.NoSave = true
	err := cmd.AddEvidence("someCreds")
	assert.Nil(t, err)
}

/**
 * Mock evidence finder function that  a receptor writer would implement
 */

type myEvidenceFinder struct{}

func (m myEvidenceFinder) ReceptorModelId() string {
	return "trr-foo"
}

func (m myEvidenceFinder) ServiceModelId() string {
	return "trs-bar"
}

func (m myEvidenceFinder) EvidenceModelId() string {
	return "tre-baz"
}

type Foo struct {
	SomeID       string    `tr:"primary_key"`
	TheNAme      string    `tr:"DisplayName:Name Of Foo;DisplayOrder:4"`
	CreatedAt    time.Time `tr:"DisplayName:Created At;DisplayOrder:3"`
	DeletedAt    time.Time `tr:"DisplayName:Deleted At;DisplayOrder:2"`
	DaysSinceXXX int       `tr:"DisplayName:Days Since;DisplayOrder:1"`
}

func (m myEvidenceFinder) GenerateEvidence(credentials string) (ev []interface{}, sources []*model.AddEvidenceRequest_Source, err error) {

	ev = []interface{}{
		&Foo{
			SomeID:       "123",
			TheNAme:      "Foo",
			CreatedAt:    time.UnixMilli(1658575538000),
			DeletedAt:    time.UnixMilli(1658575538938),
			DaysSinceXXX: 12,
		}, &Foo{
			SomeID:       "123",
			TheNAme:      "Barr",
			CreatedAt:    time.UnixMilli(1658575538000),
			DeletedAt:    time.UnixMilli(1658575538938),
			DaysSinceXXX: 12,
		}, &Foo{
			SomeID:       "123",
			TheNAme:      "Bazz",
			CreatedAt:    time.UnixMilli(1658575538000),
			DeletedAt:    time.UnixMilli(1658575538938),
			DaysSinceXXX: 12,
		},
	}
	sources = []*model.AddEvidenceRequest_Source{
		{
			Operation: "opps 1",
			Output:    "some output",
		},
	}
	return
}
