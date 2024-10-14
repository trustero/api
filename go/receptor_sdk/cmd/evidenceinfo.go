// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package cmd

import (
	"encoding/json"
	"fmt"

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
	addGrpcFlags(e.cmd)

	e.cmd.FParseErrWhitelist.UnknownFlags = true
}

type EvidenceInfo struct {
	Caption     string `json:"caption"`
	Description string `json:"description"`
}

// Cobra executes this function on evidenceinfo command.
func printEvidenceInfo(cmd *cobra.Command, args []string) (err error) {
	var (
		credentialStr string
		credentialObj interface{}
	)
	var allEvs []EvidenceInfo
	credentialStr, err = getCredentialStringFromCLI()
	if err == nil && credentialStr != "" {
		credentialObj, err = unmarshalCredentials(credentialStr, receptorImpl.GetCredentialObj())
	}
	for _, e := range receptorImpl.GetEvidenceInfo(credentialObj) {
		if e != nil {
			evidenceInfo := EvidenceInfo{
				Caption:     e.Caption,
				Description: e.Description,
			}
			allEvs = append(allEvs, evidenceInfo)
		}
	}

	evS, err := json.MarshalIndent(allEvs, "", "    ")
	if err != nil {
		return err
	}

	if string(evS) == "null" {
		println("{}")
	} else {
		println(fmt.Sprintf("%s", evS))
	}
	return
}
