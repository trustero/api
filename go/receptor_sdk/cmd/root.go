package cmd

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/client"
	receptor "github.com/trustero/api/go/receptor_v1"
)

var cfgFile string                     // Configuration file as an alternative to command line flags
var serviceProviderAccount string      // Receptor's configured service provider account
var receptorImpl receptor_sdk.Receptor // Receptor implementation

// RootCmd is the base command when called without any subcommands such as 'scan' or 'verify'.
// The RoodCmd command does nothing when invoked without a subcommand.
var RootCmd = &cobra.Command{
	Use:   "", // Set to this receptor's type
	Short: "Run a receptor in one of 2 modes: verify or scan.",
	Long: `
Run a receptor in one of 2 modes: verify or scan.  The verify mode checks the
validity of the given service provider account credential.  The scan mode either
discovers services enabled in the service provider account or reports the service
provider account's service configurations.  In the latter case, scan is invoked
with the --find-evidence flag.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// If the first argument is 'dryrun' then do not report the command
		// results to Trustero.  Instead, display the results to console.
		if len(args) > 0 && args[0] != "dryrun" {
			client.InitGRPCClient(receptor_sdk.Cert, receptor_sdk.CertServerOverride)
		} else {
			receptor_sdk.NoSave = true
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if client.ServerConn != nil {
			if err := client.ServerConn.CloseClient(); err != nil {
				log.Error().Msg(err.Error())
			}
		}
	},
}

// Execute the command
func Execute(r receptor_sdk.Receptor) {
	receptorImpl = r
	receptor_sdk.ModelID = receptorImpl.GetReceptorType()
	RootCmd.Use = receptor_sdk.ModelID
	cobra.CheckErr(RootCmd.Execute())
}

// Setup verify and scan subcommands.
var commands = []*cobra.Command{
	verifyCmd,
	scanCmd,
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file, defaults to $HOME/.receptor.yaml")
	RootCmd.PersistentFlags().StringVarP(&receptor_sdk.Host, "host", "s", "localhost", "Trustero GRPC API endpoint host name")
	RootCmd.PersistentFlags().IntVarP(&receptor_sdk.Port, "port", "p", 8888, "Trustero GRPC API endpoint port number")
	RootCmd.PersistentFlags().StringVarP(&receptor_sdk.Cert, "cert", "c", "dev", "Server cert ca to use - dev or prod")
	RootCmd.PersistentFlags().StringVarP(&receptor_sdk.CertServerOverride, "certoverride", "o", "dev.ntrce.co", "Server cert ca server override")
	RootCmd.PersistentFlags().StringVarP(&receptor_sdk.LogLevel, "level", "l", "error", "trace, debug, info, warn, error, fatal, or panic")
	RootCmd.PersistentFlags().StringVarP(&receptor_sdk.LogFile, "log-file", "", "", "Log file path")
	RootCmd.PersistentFlags().StringVarP(&receptor_sdk.ReceptorId, "receptor-id", "r", "", "Trustero receptor configuration identifier")

	if err := viper.BindPFlag("host", RootCmd.PersistentFlags().Lookup("host")); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	if err := viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port")); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	if err := viper.BindPFlag("cert", RootCmd.PersistentFlags().Lookup("cert")); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	if err := viper.BindPFlag("certoverride", RootCmd.PersistentFlags().Lookup("certoverride")); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	if err := viper.BindPFlag("level", RootCmd.PersistentFlags().Lookup("level")); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 8888)
	viper.SetDefault("cert", "dev")
	viper.SetDefault("certoverride", "localhost")
	viper.SetDefault("level", "error")

	for _, cmd := range commands {
		addReceptorFlags(cmd)
		RootCmd.AddCommand(cmd)
	}
}

func addReceptorFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&receptor_sdk.NoSave, "nosave", "n", false, "Send results to console instead of Trustero")
	cmd.PersistentFlags().StringVarP(&receptor_sdk.Notify, "notify", "", "", "Notify Trustero with Tracer ID on command completion")
	cmd.PersistentFlags().StringVarP(&receptor_sdk.CredentialsBase64URL, "credentials", "", "", "Base64 URL encoded service provider credential")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".receptor" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".receptor")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Initialize zerolog
	InitLog(receptor_sdk.LogLevel, receptor_sdk.LogFile)

	// Set GRPC host related flags if we see Host set to api.infra.trustero.com
	if strings.HasSuffix(receptor_sdk.Host, ".api.infra.trustero.com") {
		receptor_sdk.Port = 8443
		receptor_sdk.Cert = "infra"
		receptor_sdk.CertServerOverride = ""
	}
}

type commandInContext func(rc receptor.ReceptorClient, credentials interface{}) error

func invokeWithContext(token string, run commandInContext) (err error) {
	var (
		rc            receptor.ReceptorClient
		credentialStr string
		credentialObj interface{}
	)

	// Get Trustero GRPC client
	if rc, err = getReceptorClient(token); err != nil {
		return
	}

	// Get service provider account credentialStr from CLI
	credentialStr, err = getCredentialStringFromCLI()
	if err != nil {
		return
	}

	// If credentialStr not provided on CLI, get it from Trustero server
	if len(credentialStr) == 0 {
		// Get service provider account credentialStr and config from Trustero.
		var receptorInfo *receptor.ReceptorConfiguration
		if receptorInfo, err = getReceptorConfig(rc); err != nil {
			return err
		}
		credentialStr = receptorInfo.GetCredential()
		serviceProviderAccount = receptorInfo.ServiceProviderAccount
	}

	// Unmarshal credential
	credentialObj, err = receptorImpl.UnmarshalCredentials(credentialStr)

	// Invoke receptor's method
	if err == nil {
		err = run(rc, credentialObj)
	}
	return
}

func getReceptorClient(token string) (rc receptor.ReceptorClient, err error) {
	if receptor_sdk.NoSave {
		// Mock client
		rc = &MockReceptorClient{}
	} else {
		// Connect to Trustero
		if err = client.ServerConn.Dial(token, receptor_sdk.Host, receptor_sdk.Port); err != nil {
			return
		}
		// Get client
		rc = client.ServerConn.GetReceptorClient()
	}
	return
}

func getReceptorConfig(rc receptor.ReceptorClient) (config *receptor.ReceptorConfiguration, err error) {
	config, err = rc.GetConfiguration(context.Background(), &receptor.ReceptorOID{ReceptorObjectId: receptor_sdk.ReceptorId})
	return
}

func notify(rc receptor.ReceptorClient, command, result string, e error) (err error) {
	if e != nil {
		result = "error"
	}

	res := receptor.JobResult{
		TracerId:         receptor_sdk.Notify,
		ReceptorObjectId: receptor_sdk.ReceptorId,
		Command:          command,
		Result:           result}

	_, err = rc.Notify(context.Background(), &res)

	return
}

func getCredentialStringFromCLI() (credentials string, err error) {
	// Extract credentials from CLI flags
	if len(receptor_sdk.CredentialsBase64URL) > 0 {
		// Get credentials from the --credentials flag
		var creds []byte
		if creds, err = base64.URLEncoding.DecodeString(receptor_sdk.CredentialsBase64URL); err != nil {
			return
		}
		credentials = string(creds)
	} else {
		// Construct credentials from custom CLI flags.
		credentials = receptor_sdk.CredentialsFromFlags()
	}

	if len(credentials) == 0 {
		err = errors.New("no credentials provided")
	}

	return
}
