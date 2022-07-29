package cmd

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/trustero/api/go/receptor_sdk/client"
	"github.com/trustero/api/go/receptor_v1"
)

var verifyCmd = &cobra.Command{
	Use:     "verify <access_token>",
	Short:   "Verify receptor configured credentials on the receptor target endpoint.",
	Long:    ``,
	Args:    cobra.MaximumNArgs(1),
	PreRunE: verifyConfig,
	RunE:    Verify,
}

func Verify(_ *cobra.Command, args []string) (err error) {
	// Get credentials from per-receptor customized way to enter credentials.
	// This is used primarily for testing.
	if credentials := Config.CredentialsFromFlags(); credentials != nil {
		return onDebug(credentials, func(credentials interface{}) (err error) {
			_, err = Config.Verify(credentials)
			return
		})
	}

	server.Token = args[0]
	client.InitFactory(server)
	err = client.Factory.AuthScope(Config.ReceptorModelId(), doVerify)

	return
}

func doVerify(ctx context.Context, rc receptor_v1.ReceptorClient, receptorAccountConfiguration *receptor_v1.ReceptorConfiguration) (err error) {
	var ok bool
	var message string
	var credentials interface{}
	if credentials, err = Config.UnmarshallCredentials(receptorAccountConfiguration.Credential); err != nil {
		return
	}
	ok, err = Config.Verify(credentials)

	if ok, message, err = verifyWithMessage(receptorAccountConfiguration.Credential); err != nil {
		return
	}
	if NoSave {
		log.Debug().Msg(message)
		return
	}
	// Save verify results
	if _, err = rc.Verified(ctx, &receptor_v1.Credential{ReceptorObjectId: ReceptorId, Message: message, IsCredentialValid: ok}); err != nil || len(Notify) == 0 {
		return
	}
	_ = notify("verify", message, err)
	return
}

func verify(jsonCredentials string) (ok bool, err error) {
	var credentials interface{}
	if credentials, err = Config.UnmarshallCredentials(jsonCredentials); err != nil {
		return
	}
	ok, err = Config.Verify(credentials)
	return
}

func verifyWithMessage(jsonCredentials string) (ok bool, message string, err error) {
	if ok, err = verify(jsonCredentials); err != nil {
		message = "error"
	} else if ok {
		message = "successful"
	} else {
		message = "failed"
	}
	log.Debug().Msgf("verify %s.  %#v\n", message, err)
	return
}
