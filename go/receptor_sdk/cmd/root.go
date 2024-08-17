// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package cmd provides a CLI framework that implements the contract between Trustero service
// and a [receptor_v1.Receptor].  The framework allows the Receptor developer to focus on collecting
// evidences and avoid having to deal with RPC and CLI plumbing.
package cmd

import (
	"encoding/base64"
	"encoding/json"
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

var rootCmd command
var cfgFile string                     // Configuration file as an alternative to command line flags
var serviceProviderAccount string      // Receptor's configured service provider account
var receptorImpl receptor_sdk.Receptor // Receptor implementation

const (
	rootShortDesc = "Run a receptor in one of 2 modes: verify or scan."
	rootLongDesc  = `
Run a receptor in one of 2 modes: verify or scan.  The verify mode checks the
validity of the given service provider account credential.  The scan mode either
discovers services enabled in the service provider account or reports the service
provider account's service configurations.  In the latter case, scan is invoked
with the --find-evidence flag.`
)

// Setup sub commands
var cmds = map[string]command{
	"verify":     &verifi{},
	"scan":       &scann{},
	"services":   &svcs{},
	"descriptor": &desc{},
	"evidences":  &evi{},
}

// Execute is the entry point into the CLI framework.  Receptor author implements the [receptor_sdk.Receptor]
// interface and the CLI framework takes care of the rest.
func Execute(r receptor_sdk.Receptor) {
	cobra.OnInitialize(initConfig)

	// initialize cobra commands
	rootCmd = &root{}
	rootCmd.setup()
	for _, c := range cmds {
		c.setup()
		rootCmd.getCommand().AddCommand(c.getCommand())
	}

	receptorImpl = r
	receptor_sdk.ModelID = GetParsedReceptorType()
	rootCmd.getCommand().Use = receptor_sdk.ModelID
	_ = addCredentialFlags(r.GetCredentialObj())

	cobra.CheckErr(rootCmd.getCommand().Execute())
}

type command interface {
	getCommand() *cobra.Command
	setup()
}

type root struct {
	cmd *cobra.Command
}

func (r *root) getCommand() *cobra.Command {
	return r.cmd
}

func (r *root) setup() {
	r.cmd = &cobra.Command{
		Short:        rootShortDesc,
		Long:         rootLongDesc,
		SilenceUsage: true,
	}
	r.cmd.FParseErrWhitelist.UnknownFlags = true

	addStrFlag(r.cmd, &cfgFile, "config-file", "", "", "Config file, defaults to $HOME/.receptor.yaml")
	addStrFlag(r.cmd, &receptor_sdk.LogLevel, "level", "l", "error", "trace, debug, info, warn, error, fatal, or panic")
	addStrFlag(r.cmd, &receptor_sdk.LogFile, "log-file", "", "", "Log file path")
}

func addGrpcFlags(cmd *cobra.Command) {
	addStrFlag(cmd, &receptor_sdk.Host, "host", "s", "localhost", "Trustero GRPC API endpoint host name")
	addIntFlag(cmd, &receptor_sdk.Port, "port", "p", 8888, "Trustero GRPC API endpoint port number")
	addStrFlag(cmd, &receptor_sdk.Cert, "cert", "c", "dev", "Server cert ca to use - dev or prod")
	addStrFlag(cmd, &receptor_sdk.CertServerOverride, "certoverride", "o", "dev.ntrce.co", "Server cert ca server override")
	addStrFlag(cmd, &receptor_sdk.ReceptorId, "receptor-id", "r", "", "Trustero receptor configuration identifier")
	addBoolFlag(cmd, &receptor_sdk.NoSave, "nosave", "n", false, "Send results to console instead of Trustero")
	addStrFlag(cmd, &receptor_sdk.Notify, "notify", "", "", "Notify Trustero with Tracer ID on command completion")
	addStrFlag(cmd, &receptor_sdk.CredentialsBase64URL, "credentials", "", "", "Base64 URL encoded service provider credential")
	addStrFlag(cmd, &receptor_sdk.ConfigBase64URL, "config", "", "", "Base64 URL encoded receptor configuration")

}

func addStrFlagP(cmd string, p *string, name, shorthand, value, usage string) {
	if c, ok := cmds[cmd]; ok && c != nil {
		addStrFlag(c.getCommand(), p, name, shorthand, value, usage)
	}
}

func addStrFlag(cmd *cobra.Command, p *string, name, shorthand, value, usage string) {
	addFlag(cmd, name, value, func() {
		cmd.PersistentFlags().StringVarP(p, name, shorthand, value, usage)
	})
}

func addIntFlag(cmd *cobra.Command, p *int, name, shorthand string, value int, usage string) {
	addFlag(cmd, name, value, func() {
		cmd.PersistentFlags().IntVarP(p, name, shorthand, value, usage)
	})
}

func addBoolFlag(cmd *cobra.Command, p *bool, name, shorthand string, value bool, usage string) {
	addFlag(cmd, name, value, func() {
		cmd.PersistentFlags().BoolVarP(p, name, shorthand, value, usage)
	})
}

func addFlag(cmd *cobra.Command, name string, value interface{}, addFlag func()) {
	if f := cmd.PersistentFlags().Lookup(name); f == nil {
		addFlag()
		_ = viper.BindPFlag(name, f)
		viper.SetDefault(name, value)
	}
}

func grpcPreRun(_ *cobra.Command, args []string) {
	// If the first argument is 'dryrun' then do not report the command
	// results to Trustero.  Instead, display the results to console.
	if len(args) > 0 && args[0] != "dryrun" {
		client.InitGRPCClient(receptor_sdk.Cert, receptor_sdk.CertServerOverride)
	} else {
		receptor_sdk.NoSave = true
	}
}

func grpcPostRun(_ *cobra.Command, _ []string) {
	if client.ServerConn != nil {
		if err := client.ServerConn.CloseClient(); err != nil {
			log.Error().Msg(err.Error())
		}
	}
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
	initLog(receptor_sdk.LogLevel, receptor_sdk.LogFile)

	// Set GRPC host related flags if we see Host set to api.infra.trustero.com
	if strings.HasSuffix(receptor_sdk.Host, ".api.infra.trustero.com") {
		receptor_sdk.Port = 8443
		receptor_sdk.Cert = "infra"
		receptor_sdk.CertServerOverride = ""
	}
}

type commandInContext func(rc receptor.ReceptorClient, credentials interface{}, config interface{}) error

func invokeWithContext(token string, run commandInContext) (err error) {
	var (
		rc            receptor.ReceptorClient
		credentialStr string
		credentialObj interface{}
		configStr     string
		configObj     interface{}
	)

	// Get Trustero GRPC client
	if rc, err = getReceptorClient(token); err != nil {
		return
	}

	// Get service provider account credentialStr from --credentials CLI flag
	credentialStr, err = getCredentialStringFromCLI()
	// Get receptor configuration from --config CLI flag
	configStr, err = getConfigStringFromCLI()
	// If credentialStr not provided on CLI, get it from Trustero server
	if !receptor_sdk.NoSave {
		// Get service provider account credentialStr and config from Trustero.
		var receptorInfo *receptor.ReceptorConfiguration
		if receptorInfo, err = getReceptorConfig(rc); err != nil {
			return err
		}
		if len(credentialStr) == 0 {
			credentialStr = receptorInfo.GetCredential()
		}
		serviceProviderAccount = receptorInfo.ServiceProviderAccount

		if len(configStr) == 0 {
			configStr = receptorInfo.GetConfig()
		}
	}

	// Unmarshal json string credential
	if len(credentialStr) > 0 {
		credentialObj, err = unmarshalCredentials(credentialStr, receptorImpl.GetCredentialObj())
	} else {
		// If there is no credential json string provided, assume the credentials are set through
		// credential-specific CLI flags
		credentialObj = receptorImpl.GetCredentialObj()
	}

	// Unmarshal json string config
	if len(configStr) > 0 && configStr != "{}" {
		configObj, err = unmarshalConfig(configStr, receptorImpl.GetConfigObj())
	} else {
		configObj = receptorImpl.GetConfigObj()
	}

	// Invoke receptor's method
	if err == nil {
		err = run(rc, credentialObj, configObj)
	}

	// Log error
	if err != nil {
		log.Error().Err(err)
	}

	return
}

func getReceptorClient(token string) (rc receptor.ReceptorClient, err error) {
	if receptor_sdk.NoSave {
		// Mock client
		rc = &mockReceptorClient{}
	} else {
		// Connect to Trustero grpc server
		if err = client.ServerConn.Dial(token, receptor_sdk.Host, receptor_sdk.Port); err != nil {
			return
		}
		// Get grpc client
		rc = client.ServerConn.GetReceptorClient()
	}
	return
}

func getReceptorConfig(rc receptor.ReceptorClient) (config *receptor.ReceptorConfiguration, err error) {
	config, err = rc.GetConfiguration(context.Background(), &receptor.ReceptorOID{ReceptorObjectId: receptor_sdk.ReceptorId})
	return
}

func notify(rc receptor.ReceptorClient, command, result string, exceptions string, e error) (err error) {
	if e != nil {
		result = "error"
	}

	res := receptor.JobResult{
		TracerId:         receptor_sdk.Notify,
		ReceptorObjectId: receptor_sdk.ReceptorId,
		Command:          command,
		Result:           result,
		Exceptions:       exceptions,
	}

	_, err = rc.Notify(context.Background(), &res)

	return
}

func getCredentialStringFromCLI() (credentials string, err error) {
	// Extract credentials from --credentials CLI flag
	if len(receptor_sdk.CredentialsBase64URL) > 0 {
		// Get credentials from the --credentials flag
		var creds []byte
		if creds, err = base64.URLEncoding.DecodeString(receptor_sdk.CredentialsBase64URL); err != nil {
			return
		}
		credentials = string(creds)
	}

	return
}

func getConfigStringFromCLI() (config string, err error) {
	// Extract receptor configuration from --config CLI flag
	if len(receptor_sdk.ConfigBase64URL) > 0 {
		// Get receptor configuration from the --config flag
		var receptor_config []byte
		if receptor_config, err = base64.URLEncoding.DecodeString(receptor_sdk.ConfigBase64URL); err != nil {
			return
		}
		config = string(receptor_config)
	}
	return
}

func unmarshalCredentials(credentials string, credentialsObj interface{}) (obj interface{}, err error) {
	err = json.Unmarshal([]byte(credentials), credentialsObj)
	obj = credentialsObj
	return
}

func unmarshalConfig(config string, configObj interface{}) (obj interface{}, err error) {
	err = json.Unmarshal([]byte(config), &configObj)
	obj = configObj
	return
}
