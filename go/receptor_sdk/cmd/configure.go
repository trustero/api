// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	// "context"

	"github.com/trustero/api/go/receptor_v1"
)

func configure(rc receptor_v1.ReceptorClient, credentials interface{}) (err error) {

	// configure receptor
	// var discovered []*receptor_v1.ServiceEntity
	// if discovered, err = receptorImpl.Discover(credentials); err != nil {
	// 	return
	// }

	// // Report discovered services to Trustero
	// var services receptor_v1.ServiceEntities
	// services.ReceptorType = GetParsedReceptorType()
	// services.ServiceProviderAccount = serviceProviderAccount
	// services.Entities = discovered

	// // Report discovered services to Trustero
	// _, err = rc.Discovered(context.Background(), &services)
	return
}
