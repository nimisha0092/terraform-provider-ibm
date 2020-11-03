package ibm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/dns-svcs-go-sdk/dnssvcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsGlbMonitorName          = "name"
	pdnsGlbMonitorID            = "monitor_id"
	pdnsGlbMonitorDescription   = "description"
	pdnsGlbMonitorType          = "type"
	pdnsGlbMonitorPort          = "port"
	pdnsGlbMonitorInterval      = "interval"
	pdnsGlbMonitorRetries       = "retries"
	pdnsGlbMonitorTimeout       = "timeout"
	pdnsGlbMonitorMethod        = "method"
	pdnsGlbMonitorPath          = "path"
	pdnsGlbMonitorAllowInsecure = "allow_insecure"
	pdnsGlbMonitorExpectedCodes = "expected_codes"
	pdnsGlbMonitorExpectedBody  = "expected_body"
	pdnsGlbMonitorHeaders       = "headers"
	pdnsGlbMonitorHeadersName   = "name"
	pdnsGlbMonitorHeadersValue  = "value"
	pdnsGlbMonitorCreatedOn     = "created_on"
	pdnsGlbMonitorModifiedOn    = "modified_on"
)

func resourceIBMPrivateDNSGLBMonitor() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPrivateDNSGLBMonitorCreate,
		Read:     resourceIBMPrivateDNSGLBMonitorRead,
		Update:   resourceIBMPrivateDNSGLBMonitorUpdate,
		Delete:   resourceIBMPrivateDNSGLBMonitorDelete,
		Exists:   resourceIBMPrivateDNSGLBMonitorExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsGlbMonitorID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Monitor Id",
			},

			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance Id",
			},

			pdnsGlbMonitorName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a service instance.",
			},

			pdnsGlbMonitorDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the load balancer monitor",
			},

			pdnsGlbMonitorType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTP",
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS", "TCP"}),
				Description:  "The protocol to use for the health check",
			},

			pdnsGlbMonitorPort: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Port number to connect to for the health check",
			},

			pdnsGlbMonitorInterval: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "The interval between each health check",
			},

			pdnsGlbMonitorRetries: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The number of retries to attempt in case of a timeout before marking the origin as unhealthy",
			},

			pdnsGlbMonitorTimeout: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The timeout (in seconds) before marking the health check as failed",
			},

			pdnsGlbMonitorMethod: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"GET", "HEAD"}),
				Description:  "The method to use for the health check",
			},

			pdnsGlbMonitorPath: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The endpoint path to health check against",
			},

			pdnsGlbMonitorHeaders: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGlbMonitorHeadersName: {
							Type:     schema.TypeString,
							Required: true,
						},

						pdnsGlbMonitorHeadersValue: {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			pdnsGlbMonitorAllowInsecure: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Do not validate the certificate when monitor use HTTPS. This parameter is currently only valid for HTTPS monitors.",
			},

			pdnsGlbMonitorExpectedCodes: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"200", "201", "202", "203", "204", "205", "206", "207", "208", "226", "2xx"}),
				Description:  "The expected HTTP response code or code range of the health check. This parameter is only valid for HTTP and HTTPS",
			},

			pdnsGlbMonitorExpectedBody: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A case-insensitive sub-string to look for in the response body",
			},

			pdnsGlbMonitorCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Monitor creation date",
			},

			pdnsGlbMonitorModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Monitor Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSGLBMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDnsClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	CreateMonitorOptions := sess.NewCreateMonitorOptions(instanceID)

	monitorname := d.Get(pdnsGlbMonitorName).(string)
	monitorinterval := int64(d.Get(pdnsGlbMonitorInterval).(int))
	monitorretries := int64(d.Get(pdnsGlbMonitorRetries).(int))
	monitortimeout := int64(d.Get(pdnsGlbMonitorTimeout).(int))
	var monitortype string
	var monitorport int64

	if monitordescription, ok := d.GetOk(pdnsGlbMonitorDescription); ok {
		CreateMonitorOptions.SetDescription(monitordescription.(string))
	}
	if Mtype, ok := d.GetOk(pdnsGlbMonitorType); ok {
		monitortype = Mtype.(string)
	} else {
		monitortype = "HTTP"
	}

	if monitortype == "HTTP" {
		if Mport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
			monitorport = int64(Mport.(int))
		} else {
			monitorport = 80
		}
	} else if monitortype == "HTTPS" {
		if Mport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
			monitorport = int64(Mport.(int))
		} else {
			monitorport = 443
		}
	} else if monitortype == "TCP" {
		if Mport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
			monitorport = int64(Mport.(int))
		} else {
			return fmt.Errorf("Error Monitor port should be provided for TCP")
		}
	}
	CreateMonitorOptions.SetName(monitorname)
	CreateMonitorOptions.SetType(monitortype)
	CreateMonitorOptions.SetPort(monitorport)
	CreateMonitorOptions.SetInterval(monitorinterval)
	CreateMonitorOptions.SetRetries(monitorretries)
	CreateMonitorOptions.SetTimeout(monitortimeout)

	//Setting HTTP,HTTPS,TCP specific parameters
	monitorpath, pathok := d.GetOk(pdnsGlbMonitorPath)
	monitorexpectedcodes, expectedcodeok := d.GetOk(pdnsGlbMonitorExpectedCodes)
	monitorallowinsecure, allowinsecureok := d.GetOk(pdnsGlbMonitorAllowInsecure)
	monitormethod, methodok := d.GetOk(pdnsGlbMonitorMethod)
	monitorexpectedbody, expectedbodyok := d.GetOk(pdnsGlbMonitorExpectedBody)
	monitorheaders, headerok := d.GetOk(pdnsGlbMonitorHeaders)

	if (monitortype == "HTTP") || (monitortype == "HTTPS") {
		if pathok {
			CreateMonitorOptions.SetPath((monitorpath).(string))
		} else {
			CreateMonitorOptions.SetPath("/")
		}

		if expectedcodeok {
			CreateMonitorOptions.SetExpectedCodes((monitorexpectedcodes).(string))
		} else {
			CreateMonitorOptions.SetExpectedCodes("200")
		}
		if methodok {
			CreateMonitorOptions.SetMethod((monitormethod).(string))
		} else {
			CreateMonitorOptions.SetMethod("GET")
		}
		if expectedbodyok {
			CreateMonitorOptions.SetExpectedBody((monitorexpectedbody).(string))
		}
		if allowinsecureok {
			if monitortype == "HTTPS" {
				CreateMonitorOptions.SetAllowInsecure((monitorallowinsecure).(bool))
			} else {
				return fmt.Errorf("Monitor allow_insecure is not supported in type HTTP")
			}
		} else {
			CreateMonitorOptions.SetAllowInsecure(false)
		}
		if headerok {
			expandedmonitorheaders, err := expandGLBMonitors(monitorheaders, pdnsGlbMonitorHeadersName)
			if err != nil {
				return err
			}
			CreateMonitorOptions.SetHeadersVar(expandedmonitorheaders)
		}

	} else {
		if pathok || expectedcodeok || allowinsecureok || methodok || expectedbodyok || headerok {
			return fmt.Errorf("Monitor path/expected_codes/expected_body/allow_insecure/method/headers is not supported in type TCP")
		}
	}
	response, detail, err := sess.CreateMonitor(CreateMonitorOptions)
	log.Println("createprint")
	log.Println(detail)
	if err != nil {
		return fmt.Errorf("Error creating pdns GLB monitor:%s\n%s", err, detail)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, *response.ID))
	d.Set(pdnsGlbMonitorID, *response.ID)

	return resourceIBMPrivateDNSGLBMonitorRead(d, meta)
}
func expandGLBMonitors(header interface{}, geoType string) ([]dnssvcsv1.HealthcheckHeader, error) {
	headers := header.(*schema.Set).List()
	expandheaders := make([]dnssvcsv1.HealthcheckHeader, 0)
	for _, v := range headers {
		locationConfig := v.(map[string]interface{})
		hname := locationConfig[pdnsGlbMonitorHeadersName].(string)
		headers := expandStringList(locationConfig[pdnsGlbMonitorHeadersValue].([]interface{}))
		headerItem := dnssvcsv1.HealthcheckHeader{
			Name:  &hname,
			Value: headers,
		}
		expandheaders = append(expandheaders, headerItem)
	}
	return expandheaders, nil
}

func resourceIBMPrivateDNSGLBMonitorRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDnsClientSession()
	if err != nil {
		return err
	}
	idset := strings.Split(d.Id(), "/")

	getMonitorOptions := sess.NewGetMonitorOptions(idset[0], idset[1])
	response, detail, err := sess.GetMonitor(getMonitorOptions)
	if err != nil {
		return fmt.Errorf("Error fetching pdns GLB Monitor:%s\n%s", err, detail)
	}

	d.Set("id", response.ID)
	d.Set(pdnsInstanceID, idset[0])
	d.Set(pdnsGlbMonitorID, response.ID)
	d.Set(pdnsGlbMonitorName, response.Name)
	d.Set(pdnsGlbMonitorCreatedOn, response.CreatedOn)
	d.Set(pdnsGlbMonitorModifiedOn, response.ModifiedOn)
	d.Set(pdnsGlbMonitorType, response.Type)
	d.Set(pdnsGlbMonitorPort, response.Port)
	d.Set(pdnsGlbMonitorInterval, response.Interval)
	d.Set(pdnsGlbMonitorRetries, response.Retries)
	d.Set(pdnsGlbMonitorTimeout, response.Timeout)
	if response.Description != nil {
		d.Set(pdnsGlbMonitorDescription, response.Description)
	}

	if *response.Type == "HTTP" || *response.Type == "HTTPS" {
		d.Set(pdnsGlbMonitorMethod, response.Method)
		d.Set(pdnsGlbMonitorPath, response.Path)
		d.Set(pdnsGlbMonitorExpectedCodes, response.ExpectedCodes)
		d.Set(pdnsGlbMonitorHeaders, flattenDataSourceLoadBalancerHeader(response.HeadersVar))
		if response.ExpectedBody != nil {
			d.Set(pdnsGlbMonitorExpectedBody, response.ExpectedBody)
		}
	}
	if *response.Type == "HTTPS" {
		if response.AllowInsecure != nil {
			d.Set(pdnsGlbMonitorAllowInsecure, response.AllowInsecure)
		}
	}

	return nil
}

func flattenDataSourceLoadBalancerHeader(header []dnssvcsv1.HealthcheckHeader) interface{} {
	flattened := make([]interface{}, 0)
	for k, v := range header {
		cfg := map[string]interface{}{
			pdnsGlbMonitorHeadersName:  k,
			pdnsGlbMonitorHeadersValue: v,
		}
		flattened = append(flattened, cfg)
	}
	return flattened
}

func resourceIBMPrivateDNSGLBMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDnsClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")

	// Update PDNS GLB Monitor if attributes has any change

	if d.HasChange(pdnsGlbMonitorName) ||
		d.HasChange(pdnsGlbMonitorDescription) ||
		d.HasChange(pdnsGlbMonitorInterval) ||
		d.HasChange(pdnsGlbMonitorRetries) ||
		d.HasChange(pdnsGlbMonitorTimeout) ||
		d.HasChange(pdnsGlbMonitorExpectedBody) {
		updateMonitorOptions := sess.NewUpdateMonitorOptions(idset[0], idset[1])
		uname := d.Get(pdnsGlbMonitorName).(string)
		udescription := d.Get(pdnsGlbMonitorDescription).(string)
		uinterval := int64(d.Get(pdnsGlbMonitorInterval).(int))
		uretries := int64(d.Get(pdnsGlbMonitorRetries).(int))
		utimeout := int64(d.Get(pdnsGlbMonitorTimeout).(int))
		updateMonitorOptions.SetName(uname)
		updateMonitorOptions.SetDescription(udescription)
		updateMonitorOptions.SetInterval(uinterval)
		updateMonitorOptions.SetRetries(uretries)
		updateMonitorOptions.SetTimeout(utimeout)

		_, detail, err := sess.UpdateMonitor(updateMonitorOptions)

		if err != nil {
			return fmt.Errorf("Error updating pdns GLB Monitor:%s\n%s", err, detail)
		}
	}

	if d.HasChange(pdnsGlbMonitorType) ||
		d.HasChange(pdnsGlbMonitorPort) ||
		d.HasChange(pdnsGlbMonitorPath) ||
		d.HasChange(pdnsGlbMonitorAllowInsecure) ||
		d.HasChange(pdnsGlbMonitorExpectedCodes) ||
		d.HasChange(pdnsGlbMonitorExpectedBody) ||
		d.HasChange(pdnsGlbMonitorExpectedBody) ||
		d.HasChange(pdnsGlbMonitorHeaders) {

		updateMonitorOptions := sess.NewUpdateMonitorOptions(idset[0], idset[1])

		var monitortype string
		var monitorport int64
		if umonitortype, ok := d.GetOk(pdnsGlbMonitorType); ok {
			monitortype = (umonitortype).(string)
		} else {
			monitortype = "HTTP"
		}

		if monitortype == "HTTP" {
			if umonitorport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
				monitorport = int64((umonitorport).(int))
			} else {
				monitorport = 80
			}
		} else if monitortype == "HTTPS" {
			if umonitorport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
				monitorport = int64((umonitorport).(int))
			} else {
				monitorport = 443
			}
		} else if monitortype == "TCP" {
			if umonitorport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
				monitorport = int64(umonitorport.(int))
			} else {
				return fmt.Errorf("Error Monitor port should be provided for TCP")
			}
		}
		updateMonitorOptions.SetType(monitortype)
		updateMonitorOptions.SetPort(monitorport)

		monitorpath, pathok := d.GetOk(pdnsGlbMonitorPath)
		monitorexpectedcodes, expectedcodeok := d.GetOk(pdnsGlbMonitorExpectedCodes)
		monitorallowinsecure, allowinsecureok := d.GetOk(pdnsGlbMonitorAllowInsecure)
		monitormethod, methodok := d.GetOk(pdnsGlbMonitorMethod)
		monitorexpectedbody, expectedbodyok := d.GetOk(pdnsGlbMonitorExpectedBody)
		monitorheaders, headerok := d.GetOk(pdnsGlbMonitorHeaders)

		if (monitortype == "HTTP") || (monitortype == "HTTPS") {
			if pathok {
				updateMonitorOptions.SetPath((monitorpath).(string))
			} else {
				updateMonitorOptions.SetPath("/")
			}
			if expectedcodeok {
				updateMonitorOptions.SetExpectedCodes((monitorexpectedcodes).(string))
			} else {
				updateMonitorOptions.SetExpectedCodes("200")
			}
			if methodok {
				updateMonitorOptions.SetMethod((monitormethod).(string))
			} else {
				updateMonitorOptions.SetMethod("GET")
			}
			if expectedbodyok {
				updateMonitorOptions.SetExpectedBody((monitorexpectedbody).(string))
			}
			if allowinsecureok {
				if monitortype == "HTTPS" {
					updateMonitorOptions.SetAllowInsecure((monitorallowinsecure).(bool))
				} else {
					return fmt.Errorf("Monitor allow_insecure is not supported in type HTTP")
				}
			} else {
				updateMonitorOptions.SetAllowInsecure(false)
			}
			if headerok {
				expandedmonitorheaders, err := expandGLBMonitors(monitorheaders, pdnsGlbMonitorHeadersName)
				if err != nil {
					return err
				}
				updateMonitorOptions.SetHeadersVar(expandedmonitorheaders)
			}

		} else {
			if pathok || expectedcodeok || allowinsecureok || methodok || expectedbodyok || headerok {
				return fmt.Errorf("Monitor path/headers/expected_codes/expected_body/allow_insecure/method/headers is not supported in type TCP")
			}
		}
		_, detail, err := sess.UpdateMonitor(updateMonitorOptions)

		if err != nil {
			return fmt.Errorf("Error updating pdns GLB Monitor:%s\n%s", err, detail)
		}
	}

	return resourceIBMPrivateDNSGLBMonitorRead(d, meta)
}

func resourceIBMPrivateDNSGLBMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDnsClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")

	DeleteMonitorOptions := sess.NewDeleteMonitorOptions(idset[0], idset[1])
	response, err := sess.DeleteMonitor(DeleteMonitorOptions)

	if err != nil {
		return fmt.Errorf("Error deleting pdns GLB Monitor:%s\n%s", err, response)
	}

	d.SetId("")
	return nil
}

func resourceIBMPrivateDNSGLBMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(ClientSession).PrivateDnsClientSession()
	if err != nil {
		return false, err
	}

	idset := strings.Split(d.Id(), "/")

	getMonitorOptions := sess.NewGetMonitorOptions(idset[0], idset[1])
	response, detail, err := sess.GetMonitor(getMonitorOptions)
	if err != nil {
		if response != nil && detail != nil && detail.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
