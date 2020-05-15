package github

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/v29/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGithubMemberPrivileges_basic(t *testing.T) {
	var organization github.Organization

	rn := "github_member_privileges.test_member_privileges"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubMemberPrivilegesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubMemberPrivilegesConfig("false", "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubMemberPrivilegesExists(rn, &organization),
					testAccCheckGithubMemberPrivilegesAttributes(&organization, &testAccGithubMemberPrivilegesExpectedAttributes{
						MembersCanCreatePublicRepos:   false,
						MembersCanCreatePrivateRepos:  true,
						MembersCanCreateInternalRepos: true,
					}),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testAccGithubMemberPrivilegesExpectedAttributes struct {
	MembersCanCreatePublicRepos   bool
	MembersCanCreatePrivateRepos  bool
	MembersCanCreateInternalRepos bool
}

func testAccCheckGithubMemberPrivilegesExists(n string, org *github.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		memberPrivilegesName := rs.Primary.ID
		if memberPrivilegesName == "" {
			return fmt.Errorf("No Member Privileges is set")
		}

		defaultOrg := testAccProvider.Meta().(*Organization)
		conn := defaultOrg.v3client
		gotOrg, _, err := conn.Organizations.Get(context.TODO(), defaultOrg.name)
		if err != nil {
			return err
		}
		*org = *gotOrg
		return nil
	}
}

func testAccCheckGithubMemberPrivilegesAttributes(org *github.Organization, want *testAccGithubMemberPrivilegesExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if members_can_create_public_repositories := org.GetMembersCanCreatePublicRepos(); members_can_create_public_repositories != want.MembersCanCreatePublicRepos {
			return fmt.Errorf("got members_can_create_public_repositories for org %#v; want %#v", members_can_create_public_repositories, want.MembersCanCreatePublicRepos)
		}

		if members_can_create_private_repositories := org.GetMembersCanCreatePrivateRepos(); members_can_create_private_repositories != want.MembersCanCreatePrivateRepos {
			return fmt.Errorf("got members_can_create_private_repositories for org %#v; want %#v", members_can_create_private_repositories, want.MembersCanCreatePrivateRepos)
		}

		if members_can_create_internal_repositories := org.GetMembersCanCreateInternalRepos(); members_can_create_internal_repositories != want.MembersCanCreateInternalRepos {
			return fmt.Errorf("got members_can_create_internal_repositories for org %#v; want %#v", members_can_create_internal_repositories, want.MembersCanCreateInternalRepos)
		}

		return nil
	}
}

func testAccCheckGithubMemberPrivilegesDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Organization).v3client
	orgName := testAccProvider.Meta().(*Organization).name

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "github_member_privileges" {
			continue
		}

		gotOrg, resp, err := conn.Organizations.Get(context.TODO(), orgName)
		if err == nil {
			if name := gotOrg.GetName(); gotOrg != nil && name == rs.Primary.ID {
				return fmt.Errorf("Member Privileges on organization %s still exists", name)
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGithubMemberPrivilegesConfig(members_can_create_public_repositories string, members_can_create_private_repositories string, members_can_create_internal_repositories string) string {
	return fmt.Sprintf(`
resource "github_member_privileges" "test_member_privileges" {
	members_can_create_public_repositories = "%s"
	members_can_create_private_repositories = "%s"
	members_can_create_internal_repositories = "%s"
}
`, members_can_create_public_repositories, members_can_create_private_repositories, members_can_create_internal_repositories)
}
