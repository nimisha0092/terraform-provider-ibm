---
layout: "ibm"
page_title: "IBM : "
sidebar_current: "docs-ibm-datasources-dns-glb-monitors"
description: |-
  Manages IBM Cloud Infrastructure Private Domain Name Service GLB monitors.
---

# ibm\_dns_glb_monitors

Import the details of an existing IBM Cloud Infrastructure private domain name service GLB monitors as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```hcl

data "ibm_dns_glbs" "ds_pdns_glbs" {
  instance_id = "resource_instance_guid"
  zone_id= "resource_instance_id"
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, string) The GUID of the private DNS. 
* `zone_id` - (Required, string) The ID of the private DNS. 


## Attribute Reference

The following attributes are exported:

* `dns_glbs` - List of all private domain name service Load balancers in the IBM Cloud Infrastructure.
  * `name` -  The name of the load balancer monitor.
  * `description` -   Descriptive text of the load balancer monitor.
  * `ttl` -  Time to live in second.
  * `fallback_pool` -  The pool ID to use when all other pools are detected as unhealthy.
  * `default_pools` -  TA list of pool IDs ordered by their failover priority.
  * `az_pools` - Map availability zones to pool ID's.
    * `availability_zone` -  Availability zone..
    * `pools` -  List of load balancer pools.
  * `created_on` -  Load Balancer creation date. 
  * `modified_on` -  Load Balancer Modification date. 
  * `glb_id` - Load balancer Id.
  * `health` - Healthy state of the load balancer.Possible values: [DOWN,UP,DEGRADED]

