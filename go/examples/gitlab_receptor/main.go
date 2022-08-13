// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package main

import (
	"strconv"
	"time"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/trustero/api/go/receptor_v1"
	"github.com/xanzy/go-gitlab"
)

type GitLabUser struct {
	Username         string     `trustero:"id:;display:Username;order:1"`
	Name             string     `trustero:"display:Name;order:2"`
	Group            string     `trustero:"display:Group;order:3"`
	IsAdmin          bool       `trustero:"display:Admin;order:4"`
	CreatedAt        *time.Time `trustero:"display:Created On;order:5"`
	TwoFactorEnabled bool       `trustero:"display:MFA Enabled;order:6"`
	LastActivityOn   *time.Time `trustero:"display:Last Activity On;order:7"`
}

const (
	serviceName  = "GitLab"
	groupEntity  = "Group"
	memberEntity = "Member"
)

// Credential object
type Receptor struct {
	Token   string `trustero:"display:GitLab Access Token;placeholder:token"`
	GroupID string `trustero:"display:GitLab Group ID;placeholder:group id"`
}

func (r *Receptor) GetReceptorType() string {
	return "gitlab_receptor"
}

func (r *Receptor) GetKnownServices() []string {
	return []string{serviceName}
}

func (r *Receptor) GetCredentialObj() (credentialObj interface{}) {
	return r
}

func (r *Receptor) Verify(credentials interface{}) (ok bool, err error) {
	c := credentials.(*Receptor)
	ok = true
	var git *gitlab.Client
	if git, err = gitlab.NewClient(c.Token); err == nil {
		if _, _, err = git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{}); err != nil {
			ok = false
			return
		}
	}
	return
}

func (r *Receptor) Discover(credentials interface{}) (svcs []*receptor_v1.ServiceEntity, err error) {
	c := credentials.(*Receptor)
	var git *gitlab.Client

	services := receptor_sdk.NewServiceEntities()
	if git, err = gitlab.NewClient(c.Token); err == nil {
		// Get Group's name
		var group *gitlab.Group
		group, _, err = git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{})
		services.AddService(serviceName, groupEntity, group.Name, strconv.Itoa(group.ID))
	}
	return services.Entities, err
}

func (r *Receptor) Report(credentials interface{}) (evidences []*receptor_sdk.Evidence, err error) {
	c := credentials.(*Receptor)
	report := receptor_sdk.NewReport()
	var git *gitlab.Client
	if git, err = gitlab.NewClient(c.Token); err == nil {

		// Report GitLab group member information as evidence
		var ev *receptor_sdk.Evidence
		if ev, err = r.getMemberEvidence(c, git); err == nil {
			report.AddEvidence(ev)
		}
	}
	return report.Evidences, err
}

func (r *Receptor) getMemberEvidence(credentials interface{}, git *gitlab.Client) (evidence *receptor_sdk.Evidence, err error) {
	c := credentials.(*Receptor)
	evidence = receptor_sdk.NewEvidence(serviceName, memberEntity, "Group Members",
		"List of GitLab group and inherited members includes whether a member has multi-factor authentication on and if they have group admin privilege.")
	var (
		user    *gitlab.User
		group   *gitlab.Group
		members []*gitlab.GroupMember
	)
	if group, _, err = git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{}); err == nil {
		if members, _, err = git.Groups.ListAllGroupMembers(c.GroupID, &gitlab.ListGroupMembersOptions{}); err == nil {
			for _, member := range members {
				user, _, err = git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})
				evidence.AddSource("git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})", user)
				evidence.AddRow(*newGitLabUser(user, group))
			}
		}
	}
	return
}

func newGitLabUser(user *gitlab.User, group *gitlab.Group) *GitLabUser {
	return &GitLabUser{
		Username:         user.Username,
		Name:             user.Name,
		Group:            group.Name,
		IsAdmin:          user.IsAdmin,
		CreatedAt:        user.CreatedAt,
		TwoFactorEnabled: user.TwoFactorEnabled,
		LastActivityOn:   (*time.Time)(user.LastActivityOn),
	}
}

func main() {
	cmd.Execute(&Receptor{})
}
