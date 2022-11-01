/* TODO: Name package */
package receptorPackage_test

import (
	"github.com/stretchr/testify/assert"
	"receptor/trr-receptorName/receptorPackage"
	"testing"
)

func TestGetReceptorTypeImpl(t *testing.T) {
	/* TODO: Write tests */
	assert.Equal(t, "trr-custom", receptorPackage.GetReceptorTypeImpl())
}

func TestGetKnownServicesImpl(t *testing.T) {
	/* TODO: Write tests */
	svcs := receptorPackage.GetKnownServicesImpl()
	assert.Len(t, svcs, 1)
	assert.Equal(t, "Custom Service", svcs[0])
}

func TestVerify(t *testing.T) {
	/* TODO: Write tests */
	ok, err := receptorPackage.VerifyImpl(1, "fake_data")
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestDiscover(t *testing.T) {
	/* TODO: Write tests */
	svcs, err := receptorPackage.DiscoverImpl(1, "fake_data")
	assert.Len(t, svcs, 0)
	assert.Nil(t, err)
}
func TestReport(t *testing.T) {
	/* TODO: Write tests */
	evs, err := receptorPackage.ReportImpl(1, "fake_data")
	assert.Len(t, evs, 0)
	assert.Nil(t, err)
}
