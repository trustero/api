package receptorPackage

import (
	receptorLog "github.com/trustero/api/go/examples/bitbucket_receptor/logging"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/xanzy/go-gitlab"
	"time"
)

const (
	groupEntity  = "Group"
	memberEntity = "Member"
)

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
func getMemberEvidence(git *gitlab.Client, groupId string) (evidence *receptor_sdk.Evidence, err error) {
	evidence = receptor_sdk.NewEvidence(serviceName, memberEntity, serviceName+" Group Members",
		"List of GitLab group and inherited members includes whether a member has multi-factor authentication on and if they have group admin privilege.")
	var (
		user    *gitlab.User
		group   *gitlab.Group
		members []*gitlab.GroupMember
	)
	if group, _, err = git.Groups.GetGroup(groupId, &gitlab.GetGroupOptions{}); err == nil {
		evidence.AddSource("git.Groups.GetGroup(c.GroupID, &gitlab.GetGroupOptions{})", group)
		if members, _, err = git.Groups.ListAllGroupMembers(groupId, &gitlab.ListGroupMembersOptions{}); err == nil {
			evidence.AddSource("git.Groups.ListAllGroupMembers(c.GroupID, &gitlab.ListGroupMembersOptions{})", members)
			for _, member := range members {
				user, _, err = git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})
				if err != nil {
					receptorLog.Err(err, "error calling GetUser in GitLab for user %v", member.ID)
				} else {
					evidence.AddSource("git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})", user)
					evidence.AddRow(*newTrusteroGitLabUser(user, group))
				}
			}
		} else {
			receptorLog.Err(err, "error calling ListAllGroupMembers in GitLab for Group %s", groupId)
		}
	} else {
		receptorLog.Err(err, "error calling GetGroup in GitLab for Group %s", groupId)
	}
	return
}

// TrusteroGitLabUser represents a type of evidence to emit to Trustero as part of a finding.
// A list of users returned from the GitLab API will be converted to this type
// and added to an evidence object. Trustero will use the "display" and "order"
// tags to generate a table to display all users in the UI.
// The "display" tag will the column heading, and "order" tag determines
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

	TGLuser := &TrusteroGitLabUser{
		Username:         user.Username,
		Name:             user.Name,
		Group:            group.Name,
		IsAdmin:          user.IsAdmin,
		TwoFactorEnabled: user.TwoFactorEnabled,
	}

	if user.LastActivityOn != nil {
		TGLuser.LastActivityOn = (*time.Time)(user.LastActivityOn)
	}
	if user.CreatedAt != nil {
		TGLuser.CreatedAt = (*time.Time)(user.CreatedAt)
	}

	return TGLuser
}
