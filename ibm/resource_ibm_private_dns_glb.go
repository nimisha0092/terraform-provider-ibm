package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsGLBName             = "name"
	pdnsGLBID               = "glb_id"
	pdnsGLBDescription      = "description"
	pdnsGLBEnabled          = "enabled"
	pdnsGLBTTL              = "ttl"
	pdnsGLBHealth           = "health"
	pdnsGLBFallbackPool     = "fallback_pool"
	pdnsGLBDefaultPool      = "default_pools"
	pdnsGLBAZPools          = "az_pools"
	pdnsGLBAvailabilityZone = "availability_zone"
	pdnsGLBAZPoolsPools     = "pools"
	pdnsGLBCreatedOn        = "created_on"
	pdnsGLBModifiedOn       = "modified_on"
)

func resourceIBMPrivateDNSGLB() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPrivateDNSGLBCreate,
		Read:     resourceIBMPrivateDNSGLBRead,
		Update:   resourceIBMPrivateDNSGLBUpdate,
		Delete:   resourceIBMPrivateDNSGLBDelete,
		Exists:   resourceIBMPrivateDNSGLBExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsGLBID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer Id",
			},

			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GUID of the private DNS.",
			},

			pdnsZoneID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone Id",
			},

			pdnsGLBName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the load balancer",
			},

			pdnsGLBDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the load balancer",
			},

			pdnsGLBEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the load balancer is enabled",
			},

			pdnsGLBTTL: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time to live in second",
			},

			pdnsGLBHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer Id",
			},

			pdnsGLBFallbackPool: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pool ID to use when all other pools are detected as unhealthy",
			},

			pdnsGLBDefaultPool: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of pool IDs ordered by their failover priority",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			pdnsGLBAZPools: {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Map availability zones to pool ID's.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGLBAvailabilityZone: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Availability zone.",
						},

						pdnsGLBAZPoolsPools: {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of load balancer pools",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			pdnsGLBCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Load Balancer creation date",
			},

			pdnsGLBModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Load Balancer Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSGLBCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	zoneID := d.Get(pdnsZoneID).(string)
	createlbOptions := sess.NewCreateLoadBalancerOptions(instanceID, zoneID)

	lbname := d.Get(pdnsGLBName).(string)
	createlbOptions.SetName(lbname)

	if description, ok := d.GetOk(pdnsGLBDescription); ok {
		createlbOptions.SetDescription(description.(string))
	}
	if enable, ok := d.GetOkExists(pdnsGLBEnabled); ok {
		createlbOptions.SetEnabled(enable.(bool))
	}
	if ttl, ok := d.GetOk(pdnsGLBTTL); ok {
		createlbOptions.SetTTL(int64(ttl.(int)))
	}
	if flbpool, ok := d.GetOk(pdnsGLBFallbackPool); ok {
		createlbOptions.SetFallbackPool(flbpool.(string))
	}

	createlbOptions.SetDefaultPools(expandStringList(d.Get(pdnsGLBDefaultPool).([]interface{})))

	if AZpools, ok := d.GetOk(pdnsGLBAZPools); ok {
		expandedAzpools, err := expandGLBAZPools(AZpools)
		if err != nil {
			return err
		}
		createlbOptions.SetAzPools(expandedAzpools)
	}

	result, resp, err := sess.CreateLoadBalancer(createlbOptions)
	if err != nil {
		log.Printf("create global load balancer failed %s", resp)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, zoneID, *result.ID))
	return resourceIBMPrivateDNSGLBRead(d, meta)
}

func expandGLBAZPools(azpool interface{}) ([]dnssvcsv1.LoadBalancerAzPoolsItem, error) {
	azpools := azpool.(*schema.Set).List()
	expandAZpools := make([]dnssvcsv1.LoadBalancerAzPoolsItem, 0)
	for _, v := range azpools {
		locationConfig := v.(map[string]interface{})
		avzone := locationConfig[pdnsGLBAvailabilityZone].(string)
		pools := expandStringList(locationConfig[pdnsGLBPools].([]interface{}))
		aZItem := dnssvcsv1.LoadBalancerAzPoolsItem{
			AvailabilityZone: &avzone,
			Pools:            pools,
		}
		expandAZpools = append(expandAZpools, aZItem)
	}
	return expandAZpools, nil
}

func resourceIBMPrivateDNSGLBRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	idset := strings.Split(d.Id(), "/")

	getlbOptions := sess.NewGetLoadBalancerOptions(idset[0], idset[1], idset[2])
	presponse, resp, err := sess.GetLoadBalancer(getlbOptions)
	if err != nil {
		return fmt.Errorf("Error fetching pdns GLB :%s\n%s", err, resp)
	}

	response := *presponse
	d.Set(pdnsGLBName, response.Name)
	d.Set(pdnsGLBID, response.ID)
	d.Set(pdnsGLBDescription, response.Description)
	d.Set(pdnsGLBEnabled, response.Enabled)
	d.Set(pdnsGLBTTL, response.TTL)
	d.Set(pdnsGLBHealth, response.Health)
	d.Set(pdnsGLBFallbackPool, response.FallbackPool)
	d.Set(pdnsGLBDefaultPool, response.DefaultPools)
	d.Set(pdnsGLBCreatedOn, response.CreatedOn)
	d.Set(pdnsGLBModifiedOn, response.ModifiedOn)
	d.Set(pdnsGLBAZPools, flattenDataSourceLoadBalancerAZpool(response.AzPools))
	return nil
}

func flattenDataSourceLoadBalancerAZpool(azpool []dnssvcsv1.LoadBalancerAzPoolsItem) interface{} {
	flattened := make([]interface{}, 0)
	for _, v := range azpool {
		cfg := map[string]interface{}{
			pdnsGLBAvailabilityZone: v.AvailabilityZone,
			pdnsGLBPools:            flattenStringList(v.Pools),
		}
		flattened = append(flattened, cfg)
	}
	return flattened
}

func resourceIBMPrivateDNSGLBUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")

	updatelbOptions := sess.NewUpdateLoadBalancerOptions(idset[0], idset[1], idset[2])

	if d.HasChange(pdnsGLBName) ||
		d.HasChange(pdnsGLBDescription) ||
		d.HasChange(pdnsGLBEnabled) ||
		d.HasChange(pdnsGLBTTL) ||
		d.HasChange(pdnsGLBFallbackPool) ||
		d.HasChange(pdnsGLBDefaultPool) ||
		d.HasChange(pdnsGLBAZPools) {

		if name, ok := d.GetOk(pdnsGLBName); ok {
			updatelbOptions.SetName(name.(string))
		}
		if description, ok := d.GetOk(pdnsGLBDescription); ok {
			updatelbOptions.SetDescription(description.(string))
		}
		if enable, ok := d.GetOkExists(pdnsGLBEnabled); ok {
			updatelbOptions.SetEnabled(enable.(bool))
		}
		if ttl, ok := d.GetOk(pdnsGLBTTL); ok {
			updatelbOptions.SetTTL(int64(ttl.(int)))
		}
		if flbpool, ok := d.GetOk(pdnsGLBFallbackPool); ok {
			updatelbOptions.SetFallbackPool(flbpool.(string))
		}

		if _, ok := d.GetOk(pdnsGLBDefaultPool); ok {
			updatelbOptions.SetDefaultPools(expandStringList(d.Get(pdnsGLBDefaultPool).([]interface{})))
		}

		if AZpools, ok := d.GetOk(pdnsGLBAZPools); ok {
			expandedAzpools, err := expandGLBAZPools(AZpools)
			if err != nil {
				return err
			}
			updatelbOptions.SetAzPools(expandedAzpools)
		}

		result, detail, err := sess.UpdateLoadBalancer(updatelbOptions)
		if err != nil {
			return fmt.Errorf("Error updating pdns GLB :%s\n%s", err, detail)
		}
		log.Printf("Load Balancer update succesful : %s", *result.ID)
	}

	return resourceIBMPrivateDNSGLBRead(d, meta)
}

func resourceIBMPrivateDNSGLBDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")
	deletelbOptions := sess.NewDeleteLoadBalancerOptions(idset[0], idset[1], idset[2])
	response, err := sess.DeleteLoadBalancer(deletelbOptions)
	if err != nil {
		return fmt.Errorf("Error deleting pdns GLB :%s\n%s", err, response)
	}
	return nil
}

func resourceIBMPrivateDNSGLBExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}
	idset := strings.Split(d.Id(), "/")
	getlbOptions := sess.NewGetLoadBalancerOptions(idset[0], idset[1], idset[2])
	_, detail, err := sess.GetLoadBalancer(getlbOptions)
	if err != nil {
		if detail != nil && detail.StatusCode == 404 {
			log.Printf("Get GLB failed with status code 404: %v", detail)
			return false, nil
		}
		log.Printf("Get GLB failed: %v", detail)
		return false, err
	}
	return true, nil
}
