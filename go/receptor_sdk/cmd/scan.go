package cmd

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/trustero/api/go/receptor_sdk/client"
	"github.com/trustero/api/go/receptor_v1"
)

type ReporterFunc func(serviceCredentials map[string]string) (ev []interface{}, err error)

var runReporters bool

var scanCmd = &cobra.Command{
	Use:     "cmd <access_token>",
	Short:   "Scan for control evidence and discover services on a receptor target endpoint.",
	Long:    ``,
	Args:    cobra.MaximumNArgs(1),
	PreRunE: verifyConfig,
	RunE:    scan,
}

func init() {
	scanCmd.PersistentFlags().BoolVarP(&runReporters, "find-evidence", "", false, "Find and report evidence of control compliance from the reported services")
}

func scan(_ *cobra.Command, args []string) (err error) {
	// Get credentials from per-receptor customized way to enter credentials.
	// This is used primarily for testing.
	if credentials := Config.CredentialsFromFlags(); credentials != nil {
		return onDebug(credentials, func(m interface{}) (err error) {
			var evidences []*receptor_v1.Evidence
			if evidences, err = getReports(credentials); err != nil {
				return
			}
			for _, report := range evidences {
				jsoned, _ := json.Marshal(report)
				log.Debug().Msg(string(jsoned))
			}
			return
		})
	}
	server.Token = args[0]
	client.InitFactory(server)
	err = client.Factory.AuthScope(Config.ReceptorModelId(), doScan)
	return
}

func doScan(ctx context.Context, rc receptor_v1.ReceptorClient, serviceCredentials *receptor_v1.ReceptorConfiguration) (err error) {
	var credentials interface{}
	if credentials, err = Config.UnmarshallCredentials(serviceCredentials.Credential); err != nil {
		return
	}

	// First, verify credentials are valid before scanning
	var ok bool
	if ok, err = Config.Verify(credentials); err != nil {
		log.Err(err).Msg("error verifying credentials")
		return
	} else if !ok {
		// Don't continue scanning if the credentials are invalid
		log.Debug().Msg("invalid credentials - aborting scan")
		return
	}
	var reports []*receptor_v1.Evidence
	if reports, err = getReports(credentials); err != nil {
		return
	}
	if _, err = rc.Report(ctx, &receptor_v1.Finding{
		Evidences: reports,
	}); err != nil || len(Notify) == 0 {
		return
	}
	_ = notify("scan", "completed", err)
	return
}

func getReports(credentials interface{}) (evidence []*receptor_v1.Evidence, err error) {
	for _, reporter := range Config.GetReporters() {
		var reports []interface{}
		var sources []*Source
		if reports, sources, err = reporter.Report(credentials); err != nil {
			return
		}
		if len(reports) == 0 && len(sources) == 0 {
			continue
		}

		evidence = append(evidence, &receptor_v1.Evidence{
			Sources:      sources,
			Caption:      reporter.Caption(),
			ServiceName:  Config.ServiceModelId(),
			EvidenceType: &receptor_v1.Evidence_Struct{},
		})
	}
	return
}
