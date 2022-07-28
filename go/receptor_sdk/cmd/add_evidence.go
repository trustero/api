package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/trustero/api/go/client"
	"github.com/trustero/api/go/receptor_v1"
	"time"

	yaml2 "github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type ReporterFunc func(serviceCredentials map[string]string) (ev []interface{}, err error)

var Reporters []ReporterFunc

var runReporters bool

var scanCmd = &cobra.Command{
	Use:   "cmd <access_token>",
	Short: "Scan for control evidence and discover services on a receptor target endpoint.",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	RunE:  scan,
}

func init() {
	scanCmd.PersistentFlags().BoolVarP(&runReporters, "find-evidence", "", false, "Find and report evidence of control compliance from the reported services")
}

func scan(_ *cobra.Command, args []string) (err error) {
	var credentials string
	// Get credentials from per-receptor customized way to enter credentials.
	// This is used primarily for testing.
	if credentials = CredentialsFromFlags(); len(credentials) > 0 {
		NoSave = true
		err = AddEvidence(credentials)
		return
	}

	if len(Credentials) > 0 {
		// Get credentials from the --credentials flag
		var creds []byte
		if creds, err = base64.URLEncoding.DecodeString(Credentials); err != nil {
			return
		}
		credentials = string(creds)
	}

	// Get credentials from ntrced
	err = runReceptorCmd(args[0], credentials, func(ctx context.Context, rc receptor_v1.ReceptorClient, credentials, config string) (err error) {
		err = AddEvidence(credentials)
		if len(Notify) == 0 {
			return
		}
		err = notify("scan", "successful", err)
		return
	})
	return
}

func AddEvidence(credentials string) (err error) {
	// First, verify credentials are valid before scanning
	var credMap map[string]string
	credMap, err = unmarshalCreds(credentials)
	isCredValid, verifyError := Verify(credMap)
	if verifyError != nil {
		log.Err(verifyError).Msg("error verifying credentials")
		err = verifyError
		return
	}

	// Update the database with the results of the Verify call

	//err = callReceptorService(func(ctx context.Context, receptorClient agent.ReceptorClient) error {
	//	rid := agent.ReceptorID{OID: ReceptorId}
	//	_, e := receptorClient.SetIsCredValid(ctx, &agent.VerifyResult{ID: &rid, IsCredValid: isCredValid})
	//	return e
	//})
	if err != nil {
		return
	}

	// Don't continue scanning if the credentials are invalid
	if !isCredValid {
		log.Error().Msg("invalid credentials")
		return errors.New("invalid credentials")
	}

	var controlClient model.ControlServiceClient
	var cancel context.CancelFunc = func() {}
	var ctx context.Context
	if !NoSave {
		controlClient = model.NewControlServiceClient(client.Ntrced.Client)
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	}
	defer cancel()

	for _, evidenceGenerator := range Reporters {
		var genericEvidence []interface{}
		var sources []*model.AddEvidenceRequest_Source
		if genericEvidence, sources, err = evidenceGenerator.GenerateEvidence(credentials); err != nil {
			return
		}
		if len(genericEvidence) == 0 && len(sources) == 0 {
			continue
		}
		var finding *evidence.Finding
		if finding, err = pkg.NewFinding(genericEvidence); err != nil {
			return
		}
		if NoSave {
			jsoned, _ := json.Marshal(finding)
			yamled, _ := yaml2.JSONToYAML(jsoned)
			log.Debug().Msg(string(yamled))
			continue
		}
		if _, err = controlClient.AddEvidence(ctx, &model.AddEvidenceRequest{
			ReceptorModelId: evidenceGenerator.ReceptorModelId(),
			EvidenceModelId: evidenceGenerator.EvidenceModelId(),
			ServiceModelId:  evidenceGenerator.ServiceModelId(),
			Evidence:        &evidence.Evidence{Findings: []*evidence.Finding{finding}},
			Sources:         sources,
		}); err != nil {
			log.Error().Err(err).Msg("error adding evidence")
		}
	}
	return
}

func unmarshalCreds(credentials string) (map[string]string, error) {
	panic("ANKIT!!! HELP!")
	return map[string]string{}, nil
}
