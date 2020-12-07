---
layout: "ibm"
page_title: "IBM : "
sidebar_current: "docs-ibm-datasources-is-vpn-gateways"
description: |-
  Manages IBM vpn gateways.
---

# ibm\_is_vpn_gateways

Import the details of an existing IBM VPN Gateways as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```hcl

data "ibm_is_vpn_gateways" "ds_vpn_gateways" {
  
}

```

## Argument Reference

The following arguments are supported:

* `start` - (Optional, string) A server-supplied token determining what resource to start the page on.
* `limit` - (Optional, string) The number of resources to return on a page.
* `resource_group_id` - (Optional, string) resource group identifiers.
* `mode` - (Optional, string) Filters the collection to VPN gateways with the specified mode.

## Attribute Reference

The following attributes are exported:

* `id` - ID of the VPN Gateway.
* `name` - VPN Gateway instance name.
* `created_at` - The date and time that this VPN gateway was created.
* `crn` - The VPN gateway's CRN.
* `members` - Collection of VPN gateway members.
  * `address` - The public IP address assigned to the VPN gateway member.
  * `role` - The high availability role assigned to the VPN gateway member.
  * `status` - The status of the VPN gateway member.
* `resource_group` - The resource group ID for this VPN gateway
* `resource_type` - The resource type(vpn_gateway)
* `status` - The status of the VPN gateway(available, deleting, failed, pending)
* `subnet` - VPNGateway subnet info
* `mode` - mode in VPN gateway(route/policy)
