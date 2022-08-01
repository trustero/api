package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
)

// Set up the 'scan' CLI subcommand.
var scanCmd = &cobra.Command{
	Use:   "scan <trustero_access_token>|'dryrun'",
	Short: "Scan for services or evidence in a service provider account.",
	Long: `
Scan for services and evidences in-use in a service provider account.  Scan
command decodes the base64 URL encoded credentials from the '--credentials'
command line flag and check it's validity.  If 'dryrun' is specified instead
of a Trustero access token, the scan command will not report the results to
Trustero and instead print the results to console.`,
	Args: cobra.MaximumNArgs(1),
	RunE: scan,
}

func init() {
	scanCmd.PersistentFlags().BoolVarP(&receptor_sdk.FindEvidence, "find-evidence", "", false,
		"Also scan for evidences in a service provider account.")
	addReceptorFlags(scanCmd)
}

// Cobra executes this function on verify command.
func scan(_ *cobra.Command, args []string) (err error) {
	// Run receptor's Verify function and report results to Trustero
	err = invokeWithContext(args[0],
		func(rc receptor_v1.ReceptorClient, credentials interface{}) (err error) {
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
				return
			}

			// Let Trustero know the credentials have been verified.
			_, err = rc.Verified(context.Background(), toVerifyResult(ok, err))

			// Discover services in-use in the service provider account
			if err = discover(rc, credentials); err != nil {
				return
			}

			// Report evidence discovered in the service provider account
			if receptor_sdk.FindEvidence {
				err = report(rc, credentials)
			}

			return
		})
	return
}
