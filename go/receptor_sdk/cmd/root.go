package cmd

import (
	"context"
	"fmt"
	"github.com/trustero/api/go/receptor_sdk/client"
	"github.com/trustero/api/go/receptor_sdk/config"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	receptor "github.com/trustero/api/go/receptor_v1"
)

var (
	cfgFile            string
	LogLevel           string
	LogFile            string
	ModelID            string
	NoSave             bool
	Notify             string
	Credentials        string
	serviceExcludeList string
	ReceptorId         string
	TenantID           string
	server             = &client.Server{}
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
		config.InitLog(LogLevel, LogFile)
	},
}

var Config Receptor

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
	RootCmd.PersistentFlags().StringVarP(&server.Host, "host", "s", "localhost", "Trustero GRPC API endpoint host name. If host is set to 'prod.api.infra.trustero.com', port, cert, and certoverride flags will be set to match.")
	RootCmd.PersistentFlags().IntVarP(&server.Port, "port", "p", 8888, "Trustero GRPC API endpoint port number.")
	RootCmd.PersistentFlags().StringVarP(&server.Cert, "cert", "c", "dev", "Server cert ca to use - dev or prod.")
	RootCmd.PersistentFlags().StringVarP(&server.CertServerOverride, "certoverride", "o", "dev.ntrce.co", "Server cert ca server override.")
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
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set GRPC host related flags if we see Host set to api.infra.trustero.com
	if strings.HasSuffix(server.Host, ".api.infra.trustero.com") {
		server.Port = 8443
		server.Cert = "infra"
		server.CertServerOverride = ""
	}
}

func notify(command, result string, e error) (err error) {
	at := receptor.JobResult{
		TracerId:         Notify,
		ReceptorObjectId: ReceptorId,
		Command:          command,
		Result:           result}

	err = client.Factory.AuthScope(Config.ReceptorModelId(),
		func(ctx context.Context, client receptor.ReceptorClient, _ *receptor.ReceptorConfiguration) (err error) {
			if e != nil {
				result = "error"
			}
			_, err = client.Notify(ctx, &at)
			return
		})

	return
}

func verifyConfig(_ *cobra.Command, args []string) (err error) {
	if Config == nil {
		err = fmt.Errorf("receptor not configured")
	}
	if Config.CredentialsFromFlags() == nil && len(args) == 0 {
		err = fmt.Errorf("credentials not provided")
	}
	return
}
