// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package receptorPackage

import (
	"strconv"

	receptorLog "github.com/trustero/api/go/examples/gitlab_receptor/logging"
	"github.com/xanzy/go-gitlab"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_v1"
)

const (
	receptorName = "trr-custom"
	serviceName  = "Custom Service"
)

func GetReceptorTypeImpl() string {
	return receptorName
}

func GetKnownServicesImpl() []string {
	return []string{serviceName}
}

func VerifyImpl(token string, groupId string) (ok bool, err error) {
	receptorLog.Info("Entering VerifyImpl")
	ok = true
	var git *gitlab.Client
	if git, err = gitlab.NewClient(token); err == nil {
		if _, _, err = git.Groups.GetGroup(groupId, &gitlab.GetGroupOptions{}); err != nil {
			receptorLog.Err(err, "could not verify, error in GetLab GetGroup for Group %s", groupId)
			ok = false
			return
		}
	}
	receptorLog.Info("Leaving VerifyImpl")
	return
}

// In this example, Discover is making a query to GET the group name.
// The group name is then added to a list of Service Entities that will be
// returned to the Trustero platform to display in the UI.
func DiscoverImpl(token string, groupId string) (svcs []*receptor_v1.ServiceEntity, err error) {
	receptorLog.Info("Entering DiscoverImpl")
	var git *gitlab.Client
	services := receptor_sdk.NewServiceEntities()
	if git, err = gitlab.NewClient(token); err == nil {
		// Get Group's name
		var group *gitlab.Group
		if group, _, err = git.Groups.GetGroup(groupId, &gitlab.GetGroupOptions{}); err == nil {
			services.AddService(serviceName, groupEntity, group.Name, strconv.Itoa(group.ID))
		} else {
			receptorLog.Err(err, "could not discover, error in GetLab GetGroup for Group %s", groupId)
		}
	}
	receptorLog.Info("Leaving DiscoverImpl")
	return services.Entities, err
}

func ReportImpl(token string, groupId string) (evidences []*receptor_sdk.Evidence, err error) {
	receptorLog.Info("Entering ReportImpl")
	report := receptor_sdk.NewReport()
	var git *gitlab.Client
	if git, err = gitlab.NewClient(token); err == nil {

		// Report GitLab group member information as evidence
		var ev *receptor_sdk.Evidence
		if ev, err = getMemberEvidence(git, groupId); err == nil {
			report.AddEvidence(ev)
		} else {
			receptorLog.Err(err, "could not generate evidence, error in getMemberEvidence")
		}
	}
	receptorLog.Info("Leaving ReportImpl")
	return report.Evidences, err
}
