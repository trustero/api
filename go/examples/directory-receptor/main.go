package main

import (
	"github.com/trustero/api/go/receptor_sdk/cmd"
)

func init() {
	cmd.Config = receptor
}
func main() {
	setupCli()
	cmd.Execute()
}

func setupCli() {
	// Let the command line framework know the receptor type (aka the receptor model identifier).

	// A receptor's credentials are opaque to ntrced's grpc api.  However, it's best if the
	// credentials are stored as a json string in the ReceptorRecord.Credentials field.
	// The following converts credentials entered on the receptor command line to the same
	// representation stored in the ReceptorRecord.Credentials field.  Providing credentials
	// on the command line should only be used for testing.
	cmd.RootCmd.PersistentFlags().StringVarP(&Path, "path", "", "", "path to list")

}
