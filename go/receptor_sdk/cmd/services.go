// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	svcsUse   = "services"
	svcsShort = "Show list of service names this receptor will collect evidence for"
)

type svcs struct {
	cmd *cobra.Command
}

func (s *svcs) getCommand() *cobra.Command {
	return s.cmd
}

func (s *svcs) setup() {
	s.cmd = &cobra.Command{
		Use:   svcsUse,
		Short: svcsShort,
		Args:  cobra.MinimumNArgs(0),
		RunE:  services,
	}
}

// Cobra executes this function on services command.
func services(_ *cobra.Command, args []string) (err error) {
	serviceNames := receptorImpl.GetKnownServices()
	if len(serviceNames) > 0 {
		for _, name := range serviceNames {
			fmt.Println(name)
		}
	}
	return
}
