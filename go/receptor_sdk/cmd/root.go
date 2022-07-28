package cmd

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trustero/api/go/client"
	receptor "github.com/trustero/api/go/receptor_v1"
)

var (
	cfgFile            string
	Host               string
	Port               int
	Cert               string
	CertServerOverride string
	LogLevel           string
	LogFile            string
	ModelID            string
	NoSave             bool
	Notify             string
	Credentials        string
	serviceExcludeList string
	ReceptorId         string
	TenantID           string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "", // override to receptor type
	Short: "Run a receptor in one of 3 modes: verify, scan, and scanall.",
	Long: `
Run a receptor in one of 3 modes: verify, scan, and scanall.  The verify mode
mock the configured receptor credentials against its intended target service.
The scan mode conducts a discovery scan of service configuration against its
target service.  And the scanall mode performs a scan for each account with
the receptor enabled.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		client.InitGRPCClient(Cert, CertServerOverride)
	},
}

func SetReceptorType(modelId string) {
	ModelID = modelId
	RootCmd.Use = ModelID
}

type CredentialsFromFlagsFunc func() string

var CredentialsFromFlags CredentialsFromFlagsFunc = func() string { return "" }

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

var commands = []*cobra.Command{
	verifyCmd,
	scanCmd,
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.receptor.yaml)")
	RootCmd.PersistentFlags().StringVarP(&Host, "host", "s", "localhost", "Trustero GRPC API endpoint host name. If host is set to 'prod.api.infra.trustero.com', port, cert, and certoverride flags will be set to match.")
	RootCmd.PersistentFlags().IntVarP(&Port, "port", "p", 8888, "Trustero GRPC API endpoint port number.")
	RootCmd.PersistentFlags().StringVarP(&Cert, "cert", "c", "dev", "Server cert ca to use - dev or prod.")
	RootCmd.PersistentFlags().StringVarP(&CertServerOverride, "certoverride", "o", "dev.ntrce.co", "Server cert ca server override.")
	RootCmd.PersistentFlags().StringVarP(&LogLevel, "level", "l", "error", "Log level, one of trace, debug, info, warn, error, fatal, or panic.")
	RootCmd.PersistentFlags().StringVarP(&LogFile, "log-file", "", "", "Log file path")
	RootCmd.PersistentFlags().StringVarP(&serviceExcludeList, "exclude", "", "", "Comma-separated list of services types to exclude")
	RootCmd.PersistentFlags().StringVarP(&ReceptorId, "receptor-id", "r", "", "Unique identifier for the receptor record.")

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
	cmd.PersistentFlags().BoolVarP(&NoSave, "nosave", "n", false, "Print command results to console in yaml instead of saving them to prod.api.infra.trustero.com.")
	cmd.PersistentFlags().StringVarP(&Notify, "notify", "", "", "Notify prod.api.infra.trustero.com the result of the command.")
	cmd.PersistentFlags().StringVarP(&Credentials, "credentials", "", "", "Base64 URL encoded service credentials in receptor native format.")
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

	// Set GRPC host related flags if we see Host set to api.infra.trustero.com
	if strings.HasSuffix(Host, ".api.infra.trustero.com") {
		Port = 8443
		Cert = "infra"
		CertServerOverride = ""
	}
}

type runInReceptorContext func(ctx context.Context, rc receptor.ReceptorClient, credentials, config string) error

func runReceptorCmd(token, credentials string, run runInReceptorContext) (err error) {
	// Connect to ntrced
	if err = client.Ntrced.Dial(token, Host, Port); err != nil {
		return
	}
	defer func(ntrced *client.NtrcedClient) {
		if err = ntrced.CloseClient(); err != nil {
			log.Error().Msg(err.Error())
		}
	}(client.Ntrced)

	// Get client w/ timeout
	rc, ctx, cancel := client.Ntrced.GetReceptorClient()
	defer cancel()

	var receptorConfig string
	if len(credentials) == 0 {
		// Get credentials and config
		var receptorInfo *receptor.ReceptorConfiguration
		if receptorInfo, err = rc.GetConfiguration(ctx, &receptor.ReceptorOID{ReceptorObjectId: ReceptorId}); err != nil {
			return err
		}
		credentials = receptorInfo.GetCredential()
		receptorConfig = receptorInfo.GetConfig()
	}

	err = run(ctx, rc, credentials, receptorConfig)
	return
}

func notify(command, result string, e error) (err error) {
	up, ctx, cancel := client.Ntrced.GetReceptorClient()
	defer cancel()

	if e != nil {
		result = "error"
	}

	at := receptor.JobResult{
		TracerId:         Notify,
		ReceptorObjectId: ReceptorId,
		Command:          command,
		Result:           result}
	_, err = up.Notify(ctx, &at)
	return
}
