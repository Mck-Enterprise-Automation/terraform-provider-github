package github

import (
	"context"
	"log"
	"net/http"

	"github.com/google/go-github/v29/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGithubMemberPrivileges() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubMemberPrivilegesUpdate,
		Read:   resourceGithubMemberPrivilegesRead,
		Update: resourceGithubMemberPrivilegesUpdate,
		Delete: resourceGithubMemberPrivilegesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"members_can_create_internal_repositories": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"members_can_create_private_repositories": {
				Type:     schema.TypeBool,
				Optional: true,
			}, "members_can_create_public_repositories": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceGithubMemberPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	err := checkOrganization(meta)
	if err != nil {
		return err
	}

	client := meta.(*Organization).v3client

	orgName := meta.(*Organization).name

	if err != nil {
		return err
	}
	ctx := context.Background()

	log.Printf("[DEBUG] Reading Member Privileges : %s", d.Id())
	org, _, err := client.Organizations.Get(ctx, orgName)
	if err != nil {
		if ghErr, ok := err.(*github.ErrorResponse); ok {
			if ghErr.Response.StatusCode == http.StatusNotModified {
				return nil
			}
			if ghErr.Response.StatusCode == http.StatusNotFound {
				log.Printf("[WARN] Removing Member Privileges %s from state because it no longer exists in GitHub",
					d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("members_can_create_internal_repositories", org.GetMembersCanCreateInternalRepos())
	d.Set("members_can_create_private_repositories", org.GetMembersCanCreatePrivateRepos())
	d.Set("members_can_create_public_repositories", org.GetMembersCanCreatePublicRepos())

	return nil
}

func resourceGithubMemberPrivilegesUpdate(d *schema.ResourceData, meta interface{}) error {
	err := checkOrganization(meta)
	if err != nil {
		return err
	}

	client := meta.(*Organization).v3client
	orgName := meta.(*Organization).name
	ctx := context.Background()

	orgRq := resourceGithubMemberPrivilegesObject(d)

	log.Printf("[DEBUG] Updating Organization: %s", orgName)
	_, _, err = client.Organizations.Edit(ctx, orgName, orgRq)
	if err != nil {
		return err
	}
	d.SetId(orgName)

	return resourceGithubMemberPrivilegesRead(d, meta)
}

func resourceGithubMemberPrivilegesObject(d *schema.ResourceData) *github.Organization {
	return &github.Organization{
		MembersCanCreatePublicRepos:   github.Bool(d.Get("members_can_create_public_repositories").(bool)),
		MembersCanCreatePrivateRepos:  github.Bool(d.Get("members_can_create_private_repositories").(bool)),
		MembersCanCreateInternalRepos: github.Bool(d.Get("members_can_create_internal_repositories").(bool)),
	}
}

func resourceGithubMemberPrivilegesDelete(d *schema.ResourceData, meta interface{}) error {
	err := checkOrganization(meta)
	if err != nil {
		return err
	}
	return err
	// TODO Explain the function
}
