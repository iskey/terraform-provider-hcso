---
subcategory: "Identity and Access Management (IAM)"
---

# hcso_identity_group

Manages an IAM user group resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

```hcl
resource "hcso_identity_group" "group_1" {
  name        = "group_1"
  description = "This is a test group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the group. The length is less than or equal to 64 bytes.

* `description` - (Optional, String) Specifies the description of the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Groups can be imported using the `id`, e.g.

```
$ terraform import hcso_identity_group.group_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
