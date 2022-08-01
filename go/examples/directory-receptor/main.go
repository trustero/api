package main

import (
	"encoding/json"
	"github.com/trustero/monorepo/receptor-sdk/cmd"
	"github.com/trustero/monorepo/receptor-sdk/examples/directory-receptor/generators"
	"github.com/trustero/monorepo/receptor-sdk/pkg"
)

var (
	Path string
)

func main() {
	setupCli()
	cmd.Execute()
}

func setupCli() {
	// Let the command line framework know the receptor type (aka the receptor model identifier).

	cmd.EvidenceGenerators = []pkg.EvidenceGenerator{
		&generators.ListAll{},
	}

	// A receptor's credentials are opaque to ntrced's grpc api.  However, it's best if the
	// credentials are stored as a json string in the ReceptorRecord.Credentials field.
	// The following converts credentials entered on the receptor command line to the same
	// representation stored in the ReceptorRecord.Credentials field.  Providing credentials
	// on the command line should only be used for testing.
	cmd.RootCmd.PersistentFlags().StringVarP(&Path, "path", "", "", "path to list")
	cmd.CredentialsFromFlags = func() string {
		if len(Path) > 0 {
			c := generators.Credentials{
				Path: Path,
			}
			if marshalled, err := json.Marshal(c); err == nil {
				return string(marshalled)
			}
		}
		return ""
	}

}
