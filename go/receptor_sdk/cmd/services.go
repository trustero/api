// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set up the 'services' CLI subcommand.
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Show list of service names this receptor will collect evidence for.",
	Args:  cobra.MinimumNArgs(0),
	RunE:  services,
}

// Cobra executes this function on verify command.
func services(_ *cobra.Command, args []string) (err error) {
	serviceNames := receptorImpl.GetKnownServices()
	if len(serviceNames) > 0 {
		for _, name := range serviceNames {
			fmt.Println(name)
		}
	}
	return
}
