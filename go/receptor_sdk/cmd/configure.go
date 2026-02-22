// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
)

const (
	configureUse   = "configure <trustero_access_token>|dryrun"
	configureShort = "Configure service provider account information"
	configureLong  = `
Configure service provider account information.  Configure command
decodes the base64 URL encoded credentials from the '--credentials' command
line flag, gets the configuration for the service provider account and sends to Trustero.  If 'dryrun' is specified instead of a
Trustero access token, the configure command will not report the results to
Trustero and instead print the configuration results to console.`
)

type confi struct {
	cmd *cobra.Command
}

func (v *confi) getCommand() *cobra.Command {
	return v.cmd
}

func (v *confi) setup() {
	v.cmd = &cobra.Command{
		Use:          configureUse,
		Short:        configureShort,
		Long:         configureLong,
		Args:         cobra.MinimumNArgs(1),
		PreRun:       grpcPreRun,
		RunE:         configure,
		PostRun:      grpcPostRun,
		SilenceUsage: true,
	}
	v.cmd.FParseErrWhitelist.UnknownFlags = true

	addGrpcFlags(v.cmd)
}

// Cobra executes this function on verify command.
func configure(_ *cobra.Command, args []string) (err error) {
	// Run receptor's Verify function and report results to Trustero
	err = invokeWithContext(args[0],
		func(rc receptor_v1.ReceptorClient, credentials interface{}, config interface{}) (err error) {
			// Send the config back to Trustero if there is additional config
			if config != nil {
				jsonBytes, err := json.Marshal(receptorImpl.GetConfigObj(credentials))
				if err != nil {
					return err
				}
				if receptor_sdk.NoSave {
					println(string(jsonBytes))

					// print to console instead of sending to Trustero when 'dryrun' is specified

					return nil

				} else {
				}
				_, err = rc.SetConfiguration(context.Background(), &receptor_v1.ReceptorConfiguration{
					ReceptorObjectId: receptor_sdk.ReceptorId,
					Config:           string(jsonBytes),
					ModelId:          receptorImpl.GetReceptorType(),
				})
			}
			return
		})
	return
}
