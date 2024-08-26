// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"github.com/spf13/cobra"
)

const (
	instructionsUse   = "instructions"
	instructionsShort = "Dump instructions Trustero internal use"
)

type instruct struct {
	cmd *cobra.Command
}

func (l *instruct) getCommand() *cobra.Command {
	return l.cmd
}

func (l *instruct) setup() {
	l.cmd = &cobra.Command{
		Use:          logoUse,
		Short:        logoShort,
		Args:         cobra.MinimumNArgs(0),
		RunE:         instructions,
		SilenceUsage: true,
	}
	l.cmd.FParseErrWhitelist.UnknownFlags = true
}

// Cobra executes this function on logo command.
func instructions(_ *cobra.Command, args []string) (err error) {
	if instructions, err := receptorImpl.GetInstructions(); err == nil {
		println(instructions)
	}
	return
}
