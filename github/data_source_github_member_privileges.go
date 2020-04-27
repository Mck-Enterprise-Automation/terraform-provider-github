package github

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGithubMemberPrivileges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGithubMemberPrivilegesRead,

		Schema: map[string]*schema.Schema{
			"members_can_create_internal_repositories": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"members_can_create_private_repositories": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"members_can_create_public_repositories": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceGithubMemberPrivilegesRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Getting GitHub MemberPrivileges:")

	client := meta.(*Organization).v3client
	orgName := meta.(*Organization).name

	ctx := context.Background()

	org, resp, err := client.Organizations.Get(ctx, orgName)

	if err != nil {
		return err
	}

	d.SetId(orgName)

	d.Set("members_can_create_internal_repositories", org.GetMembersCanCreateInternalRepos())
	d.Set("members_can_create_private_repositories", org.GetMembersCanCreatePrivateRepos())
	d.Set("members_can_create_public_repositories", org.GetMembersCanCreatePublicRepos())
	return nil
}
