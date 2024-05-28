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
	verifyUse   = "verify <trustero_access_token>|dryrun"
	verifyShort = "Verify read-only access to a service provider account"
	verifyLong  = `
Verify read-only access to a service provider account.  Verify command
decodes the base64 URL encoded credentials from the '--credentials' command
line flag and check it's validity.  If 'dryrun' is specified instead of a
Trustero access token, the verify command will not report the results to
Trustero and instead print the results to console.`
)

type verifi struct {
	cmd *cobra.Command
}

func (v *verifi) getCommand() *cobra.Command {
	return v.cmd
}

func (v *verifi) setup() {
	v.cmd = &cobra.Command{
		Use:          verifyUse,
		Short:        verifyShort,
		Long:         verifyLong,
		Args:         cobra.MinimumNArgs(1),
		PreRun:       grpcPreRun,
		RunE:         verify,
		PostRun:      grpcPostRun,
		SilenceUsage: true,
	}
	v.cmd.FParseErrWhitelist.UnknownFlags = true

	addGrpcFlags(v.cmd)
}

// Cobra executes this function on verify command.
func verify(_ *cobra.Command, args []string) (err error) {
	// Run receptor's Verify function and report results to Trustero
	err = invokeWithContext(args[0],
		func(rc receptor_v1.ReceptorClient, credentials interface{}, config interface{}) (err error) {

			// Call receptor's Verify method
			verifyResult := toVerifyResult(receptorImpl.Verify(credentials))

			// Notify behavior is different for the verify command.  When the '--notify' command line
			// flag is provided on a verify command, verify only notify Trustero of the command
			// status and does NOT invoke the Verified Trustero RPC method to save the credential
			// in the receptor record.
			if len(receptor_sdk.Notify) > 0 {
				_ = notify(rc, "verify", verifyResult.Message, err)
				return
			}

			// Let Trustero know if the service provider account credentials are valid.
			_, err = rc.Verified(context.Background(), verifyResult)

			// Send the config back to Trustero if there is additional config
			if config != nil {
				jsonBytes, err := json.Marshal(receptorImpl.GetConfigObj())
				if err != nil {
					return err
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

func toVerifyResult(ok bool, err error) *receptor_v1.Credential {
	var message string
	if err != nil {
		message = "error"
	} else if ok {
		message = "successful"
	} else {
		message = "failed"
	}

	return &receptor_v1.Credential{ReceptorObjectId: receptor_sdk.ReceptorId, Message: message, IsCredentialValid: ok}
}
