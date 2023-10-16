---
subcategory: "Organizations"
---

# hcso_organizations

Manages an Organizations organization resource within HuaweiCloud.

## Example Usage

```hcl
resource "hcso_organizations_organization" "test"{
}
```

## Argument Reference

The following arguments are supported:

* `root_tags` - (Optional, Map) Specifies the key/value to attach to the root.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the organization.

* `master_account_id` - Indicates the unique ID of the organization's management account.

* `master_account_name` - Indicates the name of the organization's management account.

* `created_at` - Indicates the time when the organization was created.

* `root_id` - Indicates the ID of the root.

* `root_name` - Indicates the name of the root.

* `root_urn` - Indicates the urn of the root.

## Import

The Organizations organization can be imported using the `id`, e.g.

```bash
$ terraform import hcso_organizations_organization.test <id>
```
