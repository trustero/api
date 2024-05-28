// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"context"

	"github.com/trustero/api/go/receptor_v1"
)

func configure(rc receptor_v1.ReceptorClient, credentials interface{}) (err error) {

	// Receptor configuration setup
	var config *receptor_v1.ReceptorConfiguration
	if config, err = receptorImpl.Configure(credentials); err != nil || config == nil {
		return
	}

	// Send receptor configuration to Trustero
	_, err = rc.SetConfiguration(context.Background(), config)
	return
}
