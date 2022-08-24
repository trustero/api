// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package main is an example of how to use the Receptor SDK to build a Receptor CLI
package main

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/trustero/api/go/receptor_v1"
	"github.com/xanzy/go-gitlab"
)

// resources used to write this receptor:
// GitLab Golang SDK:
// https://pkg.go.dev/github.com/xanzy/go-gitlab
// GitLab API documentation:
// https://docs.gitlab.com/ee/api/api_resources.html

// This struct holds the credentials the receptor needs to authenticate with the
// service provider. A display name and placeholder tag should be provided
// which will be used in the UI when activating the receptor.
// This is what will be returned in the GetCredentialObj call

// Receptor defines the GitLab service credentials required for connecting to the GitLab
// service and gathering necessary evidence to support its use.
type Receptor struct {
	Token   string `trustero:"display:GitLab Access Token;placeholder:token"`
	GroupID string `trustero:"display:GitLab Group ID;placeholder:group id"`
}

const (
	serviceName  = "GitLab"
	groupEntity  = "Group"
	memberEntity = "Member"
)

// GetReceptorType implements the [receptor_sdk.Receptor] interface.
// Set the name of the receptor in the const declaration above
// This will let the receptor inform Trustero about itself
func (r *Receptor) GetReceptorType() string {
	return "gitlab_receptor"
}

// GetKnownServices implements the [receptor_sdk.Receptor] interface.
// Set the names of the services in the const declaration above
// This will let the receptor inform Trustero about itself
// Feel free to add or remove services as needed
func (r *Receptor) GetKnownServices() []string {
	return []string{serviceName}
}

// GetCredentialObj implements the [receptor_sdk.Receptor] interface.
// This will return Receptor struct defined above when the receptor is asked to
// identify itself
func (r *Receptor) GetCredentialObj() (credentialObj interface{}) {
	return r
}

// Verify implements the [receptor_sdk.Receptor] interface.
// This function will call into the service provider API with the provided
// credentials and confirm that the credentials are valid. Usually a simple
// API call like GET org name. If the credentials are not valid,
// return a relevant error message
func (r *Receptor) Verify(credentials interface{}) (ok bool, err error) {
	c := credentials.(*Receptor)
	ok = true
	var git *gitlab.Client
	if git, err = gitlab.NewClient(c.Token); err == nil {
		if _, _, err = git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{}); err != nil {
			log.Err(err).Msgf("could not verify, error in GetLab GetGroup for Group %s", c.GroupID)
			ok = false
			return
		}
	}
	return
}

// Discover implements the [receptor_sdk.Receptor] interface.
// The Discover function returns a list of Service Entities. This function
// makes any relevant API calls to the Service Provider to gather information
// about how many Service Entity Instances are in use. If at any point this
// function runs into an error, log that error and continue

// In this example, Discover is making a query to GET the group name.
// The group name is then added to a list of Service Entities that will be
// returned to the Trustero platform to display in the UI.
func (r *Receptor) Discover(credentials interface{}) (svcs []*receptor_v1.ServiceEntity, err error) {
	c := credentials.(*Receptor)
	var git *gitlab.Client

	services := receptor_sdk.NewServiceEntities()
	if git, err = gitlab.NewClient(c.Token); err == nil {
		// Get Group's name
		var group *gitlab.Group
		if group, _, err = git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{}); err != nil {
			services.AddService(serviceName, groupEntity, group.Name, strconv.Itoa(group.ID))
		} else {
			log.Err(err).Msgf("could not discover, error in GetLab GetGroup for Group %s", c.GroupID)
		}
	}
	return services.Entities, err
}

// Report implements the [receptor_sdk.Receptor] interface.
// Report will often make the same API calls made in the Discover call, but it
// will additionally create evidences with the data returned from the API calls
func (r *Receptor) Report(credentials interface{}) (evidences []*receptor_sdk.Evidence, err error) {
	c := credentials.(*Receptor)
	report := receptor_sdk.NewReport()
	var git *gitlab.Client
	if git, err = gitlab.NewClient(c.Token); err == nil {

		// Report GitLab group member information as evidence
		var ev *receptor_sdk.Evidence
		if ev, err = r.getMemberEvidence(c, git); err == nil {
			report.AddEvidence(ev)
		} else {
			log.Err(err).Msgf("could not generate evidence, error in getMemberEvidence")
		}
	}
	return report.Evidences, err
}

// In this example, Report makes a call to a helper function, getMemberEvidence.
// getMemberEvidence then makes a query to GET all group members in the account.
// With the use of another helper function, newGitLabUser, each member in the
// list of members is converted into a GitLabUser which is defined below.
// Each converted member is then added as a Row into an evidence, with the use of evidence.AddRow
// The Evidence also has a caption and a description that will be used in the
// Trustero UI.
// NOTE: The caption will automatically be prepended with the service name,
// so in the UI, it will read: "GitLab Group Members"
// NOTE: Any time any queries that need to be made for an evidence object should
// be recorded in the evidence object via the "AddSource" command
func (r *Receptor) getMemberEvidence(credentials interface{}, git *gitlab.Client) (evidence *receptor_sdk.Evidence, err error) {
	c := credentials.(*Receptor)
	evidence = receptor_sdk.NewEvidence(serviceName, memberEntity, serviceName+" Group Members",
		"List of GitLab group and inherited members includes whether a member has multi-factor authentication on and if they have group admin privilege.")
	var (
		user    *gitlab.User
		group   *gitlab.Group
		members []*gitlab.GroupMember
	)
	if group, _, err = git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{}); err == nil {
		evidence.AddSource("git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{})", group)
		if members, _, err = git.Groups.ListAllGroupMembers(c.GroupID, &gitlab.ListGroupMembersOptions{}); err == nil {
			evidence.AddSource("git.Groups.ListAllGroupMembers(c.GroupID, &gitlab.ListGroupMembersOptions{})", members)
			for _, member := range members {
				user, _, err = git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})
				if err != nil {
					log.Err(err).Msgf("error calling GetUser in GitLab for user %s", member.ID)
				} else {
					evidence.AddSource("git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})", user)
					evidence.AddRow(*newTrusteroGitLabUser(user, group))
				}
			}
		} else {
			log.Err(err).Msgf("error calling ListAllGroupMembers in GitLab for Group %s", c.GroupID)
		}
	} else {
		log.Err(err).Msgf("error calling GetGroup in GitLab for Group %s", c.GroupID)
	}
	return
}

// TrusteroGitLabUser represents a type of evidence to emit to Trustero as part of a finding.
// A list of users returned from the GitLab API will be converted to this type
// and added to an evidence object. Trustero will use the "display" and "order"
// tags to generate a table to display all users in the UI.
// The "display" tag will the the column heading, and "order" tag determines
// the order the columns show up in the table
type TrusteroGitLabUser struct {
	Username         string     `trustero:"id:;display:Username;order:1"`
	Name             string     `trustero:"display:Name;order:2"`
	Group            string     `trustero:"display:Group;order:3"`
	IsAdmin          bool       `trustero:"display:Admin;order:4"`
	CreatedAt        *time.Time `trustero:"display:Created On;order:5"`
	TwoFactorEnabled bool       `trustero:"display:MFA Enabled;order:6"`
	LastActivityOn   *time.Time `trustero:"display:Last Activity On;order:7"`
}

// This is a helper function that converts a user returned by the GitLab API
// into a TrusteroGitLabUser
func newTrusteroGitLabUser(user *gitlab.User, group *gitlab.Group) *TrusteroGitLabUser {
	return &TrusteroGitLabUser{
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
