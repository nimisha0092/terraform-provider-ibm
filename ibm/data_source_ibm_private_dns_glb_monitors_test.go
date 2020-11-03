package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIBMPrivateDNSGlbMonitorsDataSource_basic(t *testing.T) {
	node := "data.ibm_dns_glb_monitors.test1"
	riname := fmt.Sprintf("tf-nimi-instance-%d", acctest.RandIntRange(100, 200))
	zonename := fmt.Sprintf("tf-nimi-dnszone-%d.com", acctest.RandIntRange(100, 200))
	vpcname := fmt.Sprintf("tf-nimi-vpcname-%d", acctest.RandIntRange(100, 200))
	moniname := fmt.Sprintf("tf-nimi-monitorname-%d", acctest.RandIntRange(100, 200))
	//moniname := "TestMonitorName"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMPrivateDNSGlbMonitordDataSConfig(riname, zonename, vpcname, moniname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "dns_glb_monitors.0.name"),
					resource.TestCheckResourceAttrSet(node, "dns_glb_monitors.0.interval"),
					resource.TestCheckResourceAttrSet(node, "dns_glb_monitors.0.type"),
					resource.TestCheckResourceAttrSet(node, "dns_glb_monitors.0.retries"),
					resource.TestCheckResourceAttrSet(node, "dns_glb_monitors.0.timeout"),
				),
			},
		},
	})
}

func testAccCheckIBMPrivateDNSGlbMonitordDataSConfig(riname, zonename, vpcname, moniname string) string {
	// status filter defaults to empty
	return fmt.Sprintf(`
	data "ibm_resource_group" "rg" {
		name = "Proof of Concepts"
	}	

	resource "ibm_resource_instance" "test-pdns-instance" {
		name = "%s"
		resource_group_id = data.ibm_resource_group.rg.id
		location = "global"
		service = "dns-svcs"
		plan = "standard-dns"
	}

	resource "ibm_dns_zone" "test-pdns-zone" {
		name        = "%s"
		instance_id = ibm_resource_instance.test-pdns-instance.guid
		description = "testdescription100"
		label       = "testlabel-updated100"
	  }

	resource "ibm_is_vpc" "test_pdns_vpc" {
		depends_on = [data.ibm_resource_group.rg]
		name = "%s"
		resource_group = data.ibm_resource_group.rg.id
	}  

	resource "ibm_dns_glb_monitor" "test-pdns-monitor" {
		depends_on = [ibm_dns_zone.test-pdns-zone]
		name        = "%s"
		instance_id = ibm_resource_instance.test-pdns-instance.guid
		description = "Test dns_glb_monitors"
		interval=60
		retries=3
		timeout=8
		port=8080
		type="HTTP"
		path="/"
		method="GET"
		expected_codes="200"
		expected_body="alive"
		
    }
	
    data "ibm_dns_glb_monitors" "test1" {
		instance_id = ibm_dns_zone.test-pdns-zone.instance_id
		depends_on = [ibm_dns_glb_monitor.test-pdns-monitor]	
	}`, riname, zonename, vpcname, moniname)

}
