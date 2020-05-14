---
layout: "github"
page_title: "GitHub: github_member_privileges"
description: |-
  Provides a GitHub member privileges resource.
---

# github_member_privileges 

Provides a GitHub member privileges resource.

This resource allows you to update the member privileges for your organization.
Using the resource we allow the update of repository creation privileges for users.

## Example Usage

```hcl
# Update member privileges for the organization
resource "github_member_privileges" "member_privileges_for_org" {
	members_can_create_public_repositories = "false"
	members_can_create_private_repositories = "true"
	members_can_create_internal_repositories = "true"
}
```

## Argument Reference

The following arguments are supported:

* `members_can_create_public_repositories` - (Optional) Set to `true` to enable the members will be able to create public repositories, visible to anyone.
* `members_can_create_private_repositories` - (Optional) Set to `true` to enable the members will be able to create private repositories, visible to organization members with permission.
* `members_can_create_internal_repositories` - (Optional) Set to `true` to enable the members will be able to create internal repositories, visible to all enterprise members.

