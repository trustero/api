package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
)

const (
	eviUse   = "evidenceinfo"
	eviShort = "Print evidence caption and description"
)

type evi struct {
	cmd *cobra.Command
}

func (e *evi) getCommand() *cobra.Command {
	return e.cmd
}

func (e *evi) setup() {
	e.cmd = &cobra.Command{
		Use:          eviUse,
		Short:        eviShort,
		Args:         cobra.MinimumNArgs(0),
		RunE:         printEvidenceInfo,
		SilenceUsage: true,
	}
	e.cmd.FParseErrWhitelist.UnknownFlags = true
}

type EvidenceInfo struct {
	Caption     string `json:"caption"`
	Description string `json:"description"`
}

// Cobra executes this function on evidenceinfo command.
func printEvidenceInfo(_ *cobra.Command, args []string) (err error) {
	for _, e := range receptorImpl.GetEvidenceInfo() {
		evidenceInfo := EvidenceInfo{
			Caption:     e.Caption,
			Description: e.Description,
		}
		evS, _ := json.MarshalIndent(evidenceInfo, "", "    ")
		println(evS)
	}
	return
}
