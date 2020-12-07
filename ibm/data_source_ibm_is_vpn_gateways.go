package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isvpnGateways               = "vpn_gateways"
	isVPNGatewayStart           = "start"
	isVPNGatewayLimit           = "limit"
	isVPNGatewayResourceGroupID = "resource_group_id"
	isVPNGatewayResourceType    = "resource_type"
	isVPNGatewayMOde            = "mode"
	isVPNGatewayCrn             = "crn"
)

func dataSourceIBMISVPNGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMVPNGatewaysRead,

		Schema: map[string]*schema.Schema{

			isVPNGatewayStart: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A server-supplied token determining what resource to start the page on ",
			},
			isVPNGatewayLimit: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of resources to return on a page. ",
			},
			isVPNGatewayResourceGroupID: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "resource group identifiers ",
			},
			isVPNGatewayMOde: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to VPN gateways with the specified mode ",
			},

			isvpnGateways: {
				Type:        schema.TypeList,
				Description: "Collection of VPN Gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPNGatewayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN Gateway instance name",
						},
						isVPNGatewayCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this VPN gateway was created",
						},
						isVPNGatewayCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's CRN",
						},
						isVPNGatewayMembers: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of VPN gateway members",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The public IP address assigned to the VPN gateway member",
									},

									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The high availability role assigned to the VPN gateway member",
									},

									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the VPN gateway member",
									},
								},
							},
						},

						isVPNGatewayResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group ID for this VPN gateway",
						},

						isVPNGatewayResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},

						isVPNGatewayStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN gateway",
						},

						isVPNGatewaySubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPNGateway subnet info",
						},

						isVPNGatewayMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "mode in VPN gateway(route/policy)",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMVPNGatewaysRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listvpnGWOptions := sess.NewListVPNGatewaysOptions()

	if start, ok := d.GetOk(isVPNGatewayStart); ok {
		listvpnGWOptions.SetStart(start.(string))
	}
	if limit, ok := d.GetOk(isVPNGatewayLimit); ok {
		listvpnGWOptions.SetLimit(int64(limit.(int)))
	}
	if resourcegroupid, ok := d.GetOk(isVPNGatewayResourceGroupID); ok {
		listvpnGWOptions.SetResourceGroupID(resourcegroupid.(string))
	}
	if gwmode, ok := d.GetOk(isVPNGatewayMOde); ok {
		listvpnGWOptions.SetMode(gwmode.(string))
	}

	availableVPNGateways, detail, err := sess.ListVPNGateways(listvpnGWOptions)
	if err != nil {
		return fmt.Errorf("Error reading list of VPN Gateways:%s\n%s", err, detail)
	}
	log.Println("NIMISHA======", detail)
	vpngateways := make([]map[string]interface{}, 0)
	for _, instance := range availableVPNGateways.VPNGateways {
		gateway := map[string]interface{}{}
		data := instance.(*vpcv1.VPNGateway)
		gateway[isVPNGatewayName] = *data.Name
		gateway[isVPNGatewayCreatedAt] = data.CreatedAt.String()
		gateway[isVPNGatewayResourceType] = *data.ResourceType
		gateway[isVPNGatewayStatus] = *data.Status
		gateway[isVPNGatewayMode] = *data.Mode
		gateway[isVPNGatewayResourceGroup] = *data.ResourceGroup.ID
		gateway[isVPNGatewaySubnet] = *data.Subnet.ID
		gateway[isVPNGatewayCrn] = *data.CRN

		if data.Members != nil {
			vpcMembersIpsList := make([]map[string]interface{}, 0)
			for _, memberIP := range data.Members {
				currentMemberIP := map[string]interface{}{}
				if memberIP.PublicIP != nil {
					currentMemberIP["address"] = *memberIP.PublicIP.Address
					currentMemberIP["role"] = *memberIP.Role
					currentMemberIP["status"] = *memberIP.Status
					vpcMembersIpsList = append(vpcMembersIpsList, currentMemberIP)
				}
			}
			gateway[isVPNGatewayMembers] = vpcMembersIpsList
		}

		vpngateways = append(vpngateways, gateway)
	}

	d.SetId(dataSourceIBMVPNGatewaysID(d))
	d.Set(isvpnGateways, vpngateways)
	return nil
}

// dataSourceIBMVPNGatewaysID returns a reasonable ID  list.
func dataSourceIBMVPNGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
