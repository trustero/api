// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
)

const (
	scanUse   = "scan <trustero_access_token>|dryrun"
	scanShort = "Scan for services or evidence in a service provider account"
	scanLong  = `
Scan for services and evidences in-use in a service provider account.  Scan
command decodes the base64 URL encoded credentials from the '--credentials'
command line flag and check it's validity.  If 'dryrun' is specified instead
of a Trustero access token, the scan command will not report the results to
Trustero and instead print the results to console.`
)

type scann struct {
	cmd *cobra.Command
}

func (s *scann) getCommand() *cobra.Command {
	return s.cmd
}

func (s *scann) setup() {
	s.cmd = &cobra.Command{
		Use:          scanUse,
		Short:        scanShort,
		Long:         scanLong,
		Args:         cobra.MinimumNArgs(1),
		PreRun:       grpcPreRun,
		RunE:         scan,
		PostRun:      grpcPostRun,
		SilenceUsage: true,
	}
	s.cmd.FParseErrWhitelist.UnknownFlags = true
	addGrpcFlags(s.cmd)
	addBoolFlag(s.cmd, &receptor_sdk.FindEvidence, "find-evidence", "", false,
		"Scan for evidences in a service provider account")
}

// Cobra executes this function on verify command.
func scan(_ *cobra.Command, args []string) (err error) {
	// Run receptor's Verify function and report results to Trustero
	err = invokeWithContext(args[0],
		func(rc receptor_v1.ReceptorClient, credentials interface{}, config interface{}) (err error) {
			defer func() {
				if len(receptor_sdk.Notify) == 0 {
					return
				}
				if receptor_sdk.FindEvidence {
					notify(rc, "scan", "successful", err)
				} else {
					notify(rc, "discover", "successful", err)
				}
			}()

			// Verify credentials.
			var ok bool
			if ok, err = receptorImpl.Verify(credentials); err != nil {
				log.Err(err).Msg("error verifying credentials")
				if !ok {
					_, err = rc.Verified(context.Background(), toVerifyResult(ok, err))
				}
				return
			}

			// Let Trustero know the credentials have been verified.
			_, err = rc.Verified(context.Background(), toVerifyResult(ok, err))
			if !ok {
				return
			}
			// Report evidence discovered in the service provider account
			if receptor_sdk.FindEvidence {
				err = report(rc, credentials)
			} else {
				// Discover services in-use in the service provider account only run if --find-evidence is not run since discover runs in report
				if err = discover(rc, credentials); err != nil {
					return
				}
			}

			return
		})
	return
}
