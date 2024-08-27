// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"github.com/spf13/cobra"
)

const (
	logoUse   = "logo"
	logoShort = "Dump logo Trustero internal use"
)

type logor struct {
	cmd *cobra.Command
}

func (l *logor) getCommand() *cobra.Command {
	return l.cmd
}

func (l *logor) setup() {
	l.cmd = &cobra.Command{
		Use:          logoUse,
		Short:        logoShort,
		Args:         cobra.MinimumNArgs(0),
		RunE:         logo,
		SilenceUsage: true,
	}
	l.cmd.FParseErrWhitelist.UnknownFlags = true
}

// Cobra executes this function on logo command.
func logo(_ *cobra.Command, args []string) (err error) {
	if logo, err := receptorImpl.GetLogo(); err == nil {
		println(logo)
	}
	return
}
