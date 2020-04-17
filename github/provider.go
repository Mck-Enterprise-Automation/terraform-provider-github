package github

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_TOKEN", nil),
				Description: descriptions["token"],
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_OWNER", nil),
				Description: descriptions["owner"],
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_ORGANIZATION", nil),
				Deprecated:  "Use owner field (or GITHUB_OWNER ENV variable)",
				Description: descriptions["organization"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_BASE_URL", "https://api.github.com/"),
				Description: descriptions["base_url"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"github_actions_secret":           resourceGithubActionsSecret(),
			"github_branch_protection":        resourceGithubBranchProtection(),
			"github_issue_label":              resourceGithubIssueLabel(),
			"github_membership":               resourceGithubMembership(),
			"github_organization_block":       resourceOrganizationBlock(),
			"github_organization_project":     resourceGithubOrganizationProject(),
			"github_organization_webhook":     resourceGithubOrganizationWebhook(),
			"github_project_column":           resourceGithubProjectColumn(),
			"github_repository_collaborator":  resourceGithubRepositoryCollaborator(),
			"github_repository_deploy_key":    resourceGithubRepositoryDeployKey(),
			"github_repository_file":          resourceGithubRepositoryFile(),
			"github_repository_project":       resourceGithubRepositoryProject(),
			"github_repository_webhook":       resourceGithubRepositoryWebhook(),
			"github_repository":               resourceGithubRepository(),
			"github_team_membership":          resourceGithubTeamMembership(),
			"github_team_repository":          resourceGithubTeamRepository(),
			"github_team":                     resourceGithubTeam(),
			"github_user_gpg_key":             resourceGithubUserGpgKey(),
			"github_user_invitation_accepter": resourceGithubUserInvitationAccepter(),
			"github_user_ssh_key":             resourceGithubUserSshKey(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"github_collaborators":      dataSourceGithubCollaborators(),
			"github_ip_ranges":          dataSourceGithubIpRanges(),
			"github_membership":         dataSourceGithubMembership(),
			"github_release":            dataSourceGithubRelease(),
			"github_repositories":       dataSourceGithubRepositories(),
			"github_repository":         dataSourceGithubRepository(),
			"github_team":               dataSourceGithubTeam(),
			"github_user":               dataSourceGithubUser(),
			"github_actions_public_key": dataSourceGithubActionsPublicKey(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The OAuth token used to connect to GitHub.",

		"owner": "The GitHub owner name to manage.",

		"organization": "The GitHub owner name to manage.",

		"base_url": "The GitHub Base API URL",

		"insecure": "Whether server should be accessed " +
			"without verifying the TLS certificate.",
	}
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		owner := d.Get("organization").(string)
		if owner == "" {
			owner = d.Get("owner").(string)
		}

		config := Config{
			Token:    d.Get("token").(string),
			Owner:    owner,
			BaseURL:  d.Get("base_url").(string),
			Insecure: d.Get("insecure").(bool),
		}

		meta, err := config.Clients()
		if err != nil {
			return nil, err
		}

		meta.(*Owner).StopContext = p.StopContext()

		return meta, nil
	}
}
