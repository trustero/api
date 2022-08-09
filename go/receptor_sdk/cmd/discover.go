// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package cmd

import (
	"context"

	"github.com/trustero/api/go/receptor_v1"
)

func discover(rc receptor_v1.ReceptorClient, credentials interface{}) (err error) {

	// Discover services
	var discovered []*receptor_v1.Service
	if discovered, err = receptorImpl.Discover(credentials); err != nil {
		return
	}

	// Report discovered services to Trustero
	var services receptor_v1.Services
	services.ReceptorType = receptorImpl.GetReceptorType()
	services.ServiceProviderAccount = serviceProviderAccount
	services.Services = discovered

	// Report discovered services to Trustero
	_, err = rc.Discovered(context.Background(), &services)
	return
}
