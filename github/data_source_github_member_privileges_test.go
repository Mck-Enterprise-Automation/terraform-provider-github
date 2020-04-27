package github

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGithubMemberPrivilegesDataSource_noMatchReturnsError(t *testing.T) {
	username := "admin"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckGithubMembershipDatasourceConfig(username),
				ExpectError: regexp.MustCompile(`Not Found`),
			},
		},
	})
}

func TestAccGithubMemberPrivilegesDataSource_existing(t *testing.T) {
	if testUser == "" {
		t.Skip("This test requires you to set the test user (set it by exporting GITHUB_TEST_USER)")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGithubMembershipDatasourceConfig(testUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.github_membership.test", "username", testUser),
					resource.TestCheckResourceAttrSet("data.github_membership.test", "role"),
					resource.TestCheckResourceAttrSet("data.github_membership.test", "etag"),
				),
			},
		},
	})
}

func testAccCheckGithubMemberPrivilegesDatasourceConfig(username string) string {
	return fmt.Sprintf(`
data "github_member_privileges" "test" {
  members_can_create_internal_repositories = "%s"
}
`)
}
