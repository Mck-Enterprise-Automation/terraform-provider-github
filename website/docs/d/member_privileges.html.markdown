---
layout: "github"
page_title: "GitHub: github_member_privileges"
description: |-
  Get information on organization's member privileges.
-------------------------------------------

# github_member_privileges

Use this data source to retrieve the organization's member's privileges to create repositories in the organization.

## Example Usage

```hcl
data "github_member_privileges" "member_privileges_for_org" {
}
```