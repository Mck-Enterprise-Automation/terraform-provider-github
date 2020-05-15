package github

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGithubMemberPrivilegesDataSource_basic(t *testing.T) {
	members_can_create_internal_repositories := "true"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGithubMemberPrivilegesDataSourceConfig(members_can_create_internal_repositories),
				Check:  testMemberPrivilegesCheck(),
			},
		},
	})
}

func testMemberPrivilegesCheck() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.github_member_privileges.test", "members_can_create_internal_repositories", "true"),
		resource.TestCheckResourceAttr("data.github_member_privileges.test", "members_can_create_private_repositories", "true"),
		resource.TestCheckResourceAttr("data.github_member_privileges.test", "members_can_create_public_repositories", "false"),
	)
}

func testAccCheckGithubMemberPrivilegesDataSourceConfig(members_can_create_internal_repositories string) string {
	return fmt.Sprintf(`
data "github_member_privileges" "test" {
  members_can_create_internal_repositories = "%s"
}
`, members_can_create_internal_repositories)
}
