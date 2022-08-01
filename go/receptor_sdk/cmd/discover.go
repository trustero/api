package cmd

import (
	"context"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
)

func discover(rc receptor_v1.ReceptorClient, credentials interface{}) (err error) {

	// Discover services
	var discovered []*receptor_sdk.Service
	if discovered, err = receptorImpl.Discover(credentials); err != nil {
		return
	}

	// Report discovered services to Trustero
	var services receptor_v1.Services
	services.ReceptorType = receptorImpl.GetReceptorType()
	services.ServiceProviderAccount = serviceProviderAccount
	services.Services = []*receptor_v1.Services_Service{}

	// Append discovered services to reported services
	for _, svc := range discovered {
		services.Services = append(services.Services, &receptor_v1.Services_Service{
			Name:       svc.Name,
			InstanceId: svc.InstanceId,
		})
	}

	// Report discovered services to Trustero
	_, err = rc.Discovered(context.Background(), &services)
	return
}
