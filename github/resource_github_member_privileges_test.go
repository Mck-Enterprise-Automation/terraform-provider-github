package github

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-github/v29/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGithubMemberPrivileges_basic(t *testing.T) {
	if testCollaborator == "" {
		t.Skip("Skipping because `GITHUB_TEST_COLLABORATOR` is not set")
	}

	var org github.Organization

	rn := "github_member_privileges.test_org_member_privileges"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubMemberPrivilegesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubMemberPrivilegesConfig(testCollaborator),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubMemberPrivilegesExists(rn, &org),
					testAccCheckGithubMemberPrivilegesRoleState(rn, &org),
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

func TestAccGithubMemberPrivileges_caseInsensitive(t *testing.T) {
	if testCollaborator == "" {
		t.Skip("Skipping because `GITHUB_TEST_COLLABORATOR` is not set")
	}

	var membership github.Membership
	var otherMembership github.Membership

	rn := "github_membership.test_org_membership"
	otherCase := flipUsernameCase(testCollaborator)

	if testCollaborator == otherCase {
		t.Skip("Skipping because `GITHUB_TEST_COLLABORATOR` has no letters to flip case")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGithubMemberPrivilegesConfig(testCollaborator),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubMemberPrivilegesExists(rn, &membership),
				),
			},
			{
				Config: testAccGithubMemberPrivilegesConfig(otherCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubMemberPrivilegesExists(rn, &otherMembership),
					testAccGithubMemberPrivilegesTheSame(&membership, &otherMembership),
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

func testAccCheckGithubMemberPrivilegesExists(n string, membership *github.Membership) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No membership ID is set")
		}

		conn := testAccProvider.Meta().(*Organization).v3client
		orgName, username, err := parseTwoPartID(rs.Primary.ID, "organization", "username")
		if err != nil {
			return err
		}

		githubMembership, _, err := conn.Organizations.GetOrgMembership(context.TODO(), username, orgName)
		if err != nil {
			return err
		}
		*membership = *githubMembership
		return nil
	}
}

func testAccCheckGithubMemberPrivilegesRoleState(n string, membership *github.Membership) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No membership ID is set")
		}

		conn := testAccProvider.Meta().(*Organization).v3client
		orgName, username, err := parseTwoPartID(rs.Primary.ID, "organization", "username")
		if err != nil {
			return err
		}

		githubMembership, _, err := conn.Organizations.GetOrgMembership(context.TODO(), username, orgName)
		if err != nil {
			return err
		}

		resourceRole := membership.GetRole()
		actualRole := githubMembership.GetRole()

		if resourceRole != actualRole {
			return fmt.Errorf("Membership role %v in resource does match actual state of %v",
				resourceRole, actualRole)
		}
		return nil
	}
}

func testAccGithubMemberPrivilegesConfig(username string) string {
	return fmt.Sprintf(`
  resource "github_membership" "test_org_membership" {
    username = "%s"
    role = "member"
  }
`, username)
}

func testAccGithubMemberPrivilegesTheSame(orig, other *github.Membership) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if orig.GetURL() != other.GetURL() {
			return errors.New("users are different")
		}

		return nil
	}
}
