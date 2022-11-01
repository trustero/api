/* TODO: Name package */
package receptorPackage

import (
	receptorLog "receptor/trr-receptorName/logging"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
)

const (
	receptorName = "trr-custom"
	serviceName  = "Custom Service"
)

func GetReceptorTypeImpl() string {
	return receptorName
}

func GetKnownServicesImpl() []string {
	return []string{serviceName}
}

func VerifyImpl( /* TODO: Put needed receptor creds here */ field1 int, field2 string) (ok bool, err error) {
	receptorLog.Info("Entering VerifyImpl")
	/* TODO: Implement Verify logic here */
	receptorLog.Info("Leaving VerifyImpl")
	return
}

func DiscoverImpl( /* TODO: Put needed receptor creds here */ field1 int, field2 string) (svcs []*receptor_v1.ServiceEntity, err error) {
	receptorLog.Info("Entering DiscoverImpl")
	/* TODO: Implement Discover logic here */
	receptorLog.Info("Leaving DiscoverImpl")
	return
}

func ReportImpl( /* TODO: Put needed receptor creds here */ field1 int, field2 string) (evidences []*receptor_sdk.Evidence, err error) {
	receptorLog.Info("Entering ReportImpl")
	/* TODO: Implement Report logic here */
	receptorLog.Info("Leaving ReportImpl")
	return
}
