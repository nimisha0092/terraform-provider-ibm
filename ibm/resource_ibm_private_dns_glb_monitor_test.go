package ibm

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIBMPrivateDNSGlbMonitor_Basic(t *testing.T) {
	var resultprivatedns string
	log.Printf("NIMISHAA------BEFORE-----TestAccIBMPrivateDNSGlbMonitor_Basic")
	name := fmt.Sprintf("testpdnspn%s.com", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPrivateDNSGlbMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMPrivateDNSGlbMonitorBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPrivateDNSGlbMonitorExists("ibm_dns_glb_monitor.test-pdns-monitor", resultprivatedns),
				),
			},
		},
	})
	log.Printf("NIMISHAA------AFTER-----TestAccIBMPrivateDNSGlbMonitor_Basic")
}

func TestAccIBMPrivateDNSGlbMonitorImport(t *testing.T) {
	var resultprivatedns string
	log.Printf("NIMISHAA------BEFORE-----TestAccIBMPrivateDNSGlbMonitorImport")
	name := fmt.Sprintf("testpdnszone%s.com", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPrivateDNSGlbMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMPrivateDNSPermittedNetworkBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPrivateDNSGlbMonitorExists("ibm_dns_glb_monitor.test-pdns-monitor", resultprivatedns),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_dns_glb_monitor.test-pdns-monitor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
	log.Printf("NIMISHAA------AFTER-----TestAccIBMPrivateDNSGlbMonitorImport")
}

func testAccCheckIBMPrivateDNSGlbMonitorBasic(name string) string {
	return fmt.Sprintf(`
	data "ibm_resource_group" "rg" {
		name = "Proof of Concepts"
    }

    resource "ibm_is_vpc" "nimisha-testpdnsvpc" {
		depends_on = [data.ibm_resource_group.rg]
		name = "nimisha-test-lb-vpc"
		resource_group = data.ibm_resource_group.rg.id
    }

    resource "ibm_resource_instance" "nimisha-test-pdns-instance" {
		depends_on = [ibm_is_vpc.nimisha-testpdnsvpc]
		name = "nimisha-test-pdns"
		resource_group_id = data.ibm_resource_group.rg.id
		location = "global"
		service = "dns-svcs"
		plan = "standard-dns"
    }

    resource "ibm_dns_zone" "nimisha-test-pdns-zone" {
		depends_on = [ibm_resource_instance.nimisha-test-pdns-instance]
		name = "%s"
		instance_id = ibm_resource_instance.nimisha-test-pdns-instance.guid
		description = "testdescription"
		label = "testlabel-updated"
    }

	resource "ibm_dns_glb_monitor" "test-pdns-monitor" {
		depends_on = [ibm_dns_zone.nimisha-test-pdns-zone]
		name = "updatenimimonitor"
		instance_id = ibm_resource_instance.nimisha-test-pdns-instance.guid
		description = "new updatetestdescription"
		interval=63
		retries=3
		timeout=8
		port=8080
		type="HTTP"
		expected_codes= "200"
		path="/health"
		method="GET"
		expected_body="alive"
		headers{
			name="nimiheader"
			value=["example","abc"]
		}	
    }
	  `, name)

}

func testAccCheckIBMPrivateDNSGlbMonitorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_dns_glb_monitor" {
			continue
		}
		log.Printf("NIMISHAA------BEFORE-----testAccCheckIBMPrivateDNSGlbMonitorDestroy")
		pdnsClient, err := testAccProvider.Meta().(ClientSession).PrivateDnsClientSession()
		if err != nil {
			return err
		}

		parts := rs.Primary.ID
		partslist := strings.Split(parts, "/")
		log.Printf("partslist[0]=%s  ---- partslist[1]=%s---", partslist[0], partslist[1])
		getMonitorOptions := pdnsClient.NewGetMonitorOptions(partslist[0], partslist[1])
		_, res, err := pdnsClient.GetMonitor(getMonitorOptions)

		if err != nil &&
			res.StatusCode != 403 &&
			!strings.Contains(err.Error(), "The service instance was disabled, any access is not allowed.") {

			return fmt.Errorf("testAccCheckIBMPrivateDNSZoneDestroy: Error checking if instance (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}
	log.Printf("NIMISHAA------BEFORE-----testAccCheckIBMPrivateDNSGlbMonitorDestroy")
	return nil
}

func testAccCheckIBMPrivateDNSGlbMonitorExists(n string, result string) resource.TestCheckFunc {
	log.Printf("NIMISHAA------BEFORE-----testAccCheckIBMPrivateDNSGlbMonitorExists")
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		log.Printf("NIMISHAA------1")
		pdnsClient, err := testAccProvider.Meta().(ClientSession).PrivateDnsClientSession()
		if err != nil {
			return err
		}
		log.Printf("NIMISHAA------2")
		parts := rs.Primary.ID
		partslist := strings.Split(parts, "/")
		log.Printf("partslist[0]=%s  ---- partslist[1]=%s---", partslist[0], partslist[1])
		getMonitorOptions := pdnsClient.NewGetMonitorOptions(partslist[0], partslist[1])
		r, res, err := pdnsClient.GetMonitor(getMonitorOptions)

		if err != nil &&
			res.StatusCode != 403 &&
			!strings.Contains(err.Error(), "The service instance was disabled, any access is not allowed.") {
			return fmt.Errorf("testAccCheckIBMPrivateDNSZoneExists: Error checking if instance (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
		log.Printf("NIMISHAA------BEFORE-----testAccCheckIBMPrivateDNSGlbMonitorExists")
		result = *r.ID
		return nil
	}
}
