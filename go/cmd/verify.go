package cmd

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trustero/api/go/receptor_v1"
)

var Verify VerifyFunc = func(map[string]string) (bool, error) { return true, nil }

var verifyCmd = &cobra.Command{
	Use:   "verify <access_token>",
	Short: "Verify receptor configured credentials on the receptor target endpoint.",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	RunE:  verify,
}

func verify(_ *cobra.Command, args []string) (err error) {
	var credentials string
	// Get credentials from per-receptor customized way to enter credentials.
	// This is used primarily for testing.
	if credentials = CredentialsFromFlags(); len(credentials) > 0 {
		NoSave = true
		var msg string
		_, msg, err = verifyWithMessage(credentials, "")
		fmt.Printf("verify %s.  %#v\n", msg, err)
		return

	} else if len(Credentials) > 0 {
		// Get credentials from the --credentials flag
		var creds []byte
		if creds, err = base64.URLEncoding.DecodeString(Credentials); err != nil {
			return
		}
		credentials = string(creds)
	}

	err = runReceptorCmd(args[0], credentials,
		func(ctx context.Context, rc receptor_v1.ReceptorClient, credentials, config string) (err error) {
			// Call Verify
			var message string
			var ok bool
			ok, message, err = verifyWithMessage(credentials, config)

			if NoSave {
				println(message)
				return
			}

			if len(Notify) > 0 {
				_ = notify("verify", message, err)
				return
			}
			// Save verify results
			_, err = rc.Verified(ctx, &receptor_v1.Credential{ReceptorObjectId: ReceptorId, Message: message, IsCredentialValid: ok})
			return
		})
	return
}

func verifyWithMessage(credentials, config string) (ok bool, message string, err error) {
	if ok, err = Verify(credentials, config); err != nil {
		message = "error"
	} else if ok {
		message = "successful"
	} else {
		message = "failed"
	}
	return
}
