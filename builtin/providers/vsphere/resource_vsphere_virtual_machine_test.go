package vsphere

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

func TestAccVSphereVirtualMachine_basic(t *testing.T) {
	var vm virtualMachine
	var locationOpt string
	var datastoreOpt string

	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		locationOpt += fmt.Sprintf("    datacenter = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_CLUSTER"); v != "" {
		locationOpt += fmt.Sprintf("    cluster = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v != "" {
		locationOpt += fmt.Sprintf("    resource_pool = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_DATASTORE"); v != "" {
		datastoreOpt = fmt.Sprintf("        datastore = \"%s\"\n", v)
	}
	template := os.Getenv("VSPHERE_TEMPLATE")
	gateway := os.Getenv("VSPHERE_IPV4_GATEWAY")
	label := os.Getenv("VSPHERE_NETWORK_LABEL")
	ip_address := os.Getenv("VSPHERE_IPV4_ADDRESS")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVSphereVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccCheckVSphereVirtualMachineConfig_basic,
					locationOpt,
					gateway,
					label,
					ip_address,
					datastoreOpt,
					template,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVSphereVirtualMachineExists("vsphere_virtual_machine.foo", &vm),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.foo", "name", "terraform-test"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.foo", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.foo", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.foo", "disk.#", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.foo", "disk.0.template.label", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.foo", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.foo", "network_interface.0.label", label),
				),
			},
		},
	})
}

func TestAccVSphereVirtualMachine_dhcp(t *testing.T) {
	var vm virtualMachine
	var locationOpt string
	var datastoreOpt string

	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		locationOpt += fmt.Sprintf("    datacenter = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_CLUSTER"); v != "" {
		locationOpt += fmt.Sprintf("    cluster = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v != "" {
		locationOpt += fmt.Sprintf("    resource_pool = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_DATASTORE"); v != "" {
		datastoreOpt = fmt.Sprintf("        datastore = \"%s\"\n", v)
	}
	template := os.Getenv("VSPHERE_TEMPLATE")
	label := os.Getenv("VSPHERE_NETWORK_LABEL_DHCP")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVSphereVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccCheckVSphereVirtualMachineConfig_dhcp,
					locationOpt,
					label,
					datastoreOpt,
					template,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVSphereVirtualMachineExists("vsphere_virtual_machine.bar", &vm),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.bar", "name", "terraform-test"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.bar", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.bar", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.bar", "disk.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.bar", "disk.0.template.label", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.bar", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.bar", "network_interface.0.label", label),
				),
			},
		},
	})
}

func TestAccVSphereVirtualMachine_custom_configs(t *testing.T) {
	var vm virtualMachine
	var locationOpt string
	var datastoreOpt string

	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		locationOpt += fmt.Sprintf("    datacenter = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_CLUSTER"); v != "" {
		locationOpt += fmt.Sprintf("    cluster = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v != "" {
		locationOpt += fmt.Sprintf("    resource_pool = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_DATASTORE"); v != "" {
		datastoreOpt = fmt.Sprintf("        datastore = \"%s\"\n", v)
	}
	template := os.Getenv("VSPHERE_TEMPLATE")
	label := os.Getenv("VSPHERE_NETWORK_LABEL_DHCP")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVSphereVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccCheckVSphereVirtualMachineConfig_custom_configs,
					locationOpt,
					label,
					datastoreOpt,
					template,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVSphereVirtualMachineExistsHasCustomConfig("vsphere_virtual_machine.car", &vm),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "name", "terraform-test-custom"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "disk.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "disk.0.template.label", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "custom_configuration_parameters.foo", "bar"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "custom_configuration_parameters.car", "ferrari"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "custom_configuration_parameters.num", "42"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.car", "network_interface.0.label", label),
				),
			},
		},
	})
}

func TestAccVSphereVirtualMachine_createInExistingFolder(t *testing.T) {
	var vm virtualMachine
	var locationOpt string
	var datastoreOpt string
	var datacenter string

	folder := "tf_test_createInExistingFolder"

	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		locationOpt += fmt.Sprintf("    datacenter = \"%s\"\n", v)
		datacenter = v
	}
	if v := os.Getenv("VSPHERE_CLUSTER"); v != "" {
		locationOpt += fmt.Sprintf("    cluster = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v != "" {
		locationOpt += fmt.Sprintf("    resource_pool = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_DATASTORE"); v != "" {
		datastoreOpt = fmt.Sprintf("        datastore = \"%s\"\n", v)
	}
	template := os.Getenv("VSPHERE_TEMPLATE")
	label := os.Getenv("VSPHERE_NETWORK_LABEL_DHCP")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckVSphereVirtualMachineDestroy,
			removeVSphereFolder(datacenter, folder, ""),
		),
		Steps: []resource.TestStep{
			resource.TestStep{
				PreConfig: func() { createVSphereFolder(datacenter, folder) },
				Config: fmt.Sprintf(
					testAccCheckVSphereVirtualMachineConfig_createInFolder,
					folder,
					locationOpt,
					label,
					datastoreOpt,
					template,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVSphereVirtualMachineExists("vsphere_virtual_machine.folder", &vm),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "name", "terraform-test-folder"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "folder", folder),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "disk.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "disk.0.template.label", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.folder", "network_interface.0.label", label),
				),
			},
		},
	})
}

func TestAccVSphereVirtualMachine_createWithFolder(t *testing.T) {
	var vm virtualMachine
	var f folder
	var locationOpt string
	var folderLocationOpt string
	var datastoreOpt string

	folder := "tf_test_createWithFolder"

	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		folderLocationOpt = fmt.Sprintf("    datacenter = \"%s\"\n", v)
		locationOpt += folderLocationOpt
	}
	if v := os.Getenv("VSPHERE_CLUSTER"); v != "" {
		locationOpt += fmt.Sprintf("    cluster = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v != "" {
		locationOpt += fmt.Sprintf("    resource_pool = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_DATASTORE"); v != "" {
		datastoreOpt = fmt.Sprintf("        datastore = \"%s\"\n", v)
	}
	template := os.Getenv("VSPHERE_TEMPLATE")
	label := os.Getenv("VSPHERE_NETWORK_LABEL_DHCP")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckVSphereVirtualMachineDestroy,
			testAccCheckVSphereFolderDestroy,
		),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccCheckVSphereVirtualMachineConfig_createWithFolder,
					folder,
					folderLocationOpt,
					locationOpt,
					label,
					datastoreOpt,
					template,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVSphereVirtualMachineExists("vsphere_virtual_machine.with_folder", &vm),
					testAccCheckVSphereFolderExists("vsphere_folder.with_folder", &f),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.with_folder", "name", "terraform-test-with-folder"),
					// resource.TestCheckResourceAttr(
					// 	"vsphere_virtual_machine.with_folder", "folder", folder),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.with_folder", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.with_folder", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.with_folder", "disk.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.with_folder", "disk.0.template.label", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.with_folder", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.with_folder", "network_interface.0.label", label),
				),
			},
		},
	})
}

// TestAccVSphereVirtualMachine_linkedNamedSnapshot tests if it works to create a linked clone from a named snapshot.
func TestAccVSphereVirtualMachine_linkedNamedSnapshot(t *testing.T) {
	var vm virtualMachine
	var locationOpt string
	var datastoreOpt string

	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		locationOpt += fmt.Sprintf("    datacenter = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_CLUSTER"); v != "" {
		locationOpt += fmt.Sprintf("    cluster = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v != "" {
		locationOpt += fmt.Sprintf("    resource_pool = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_DATASTORE"); v != "" {
		datastoreOpt = fmt.Sprintf("        datastore = \"%s\"\n", v)
	}
	template := os.Getenv("VSPHERE_TEMPLATE")
	label := os.Getenv("VSPHERE_NETWORK_LABEL_DHCP")
	snapshot := os.Getenv("VSPHERE_SNAPSHOT")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVSphereVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccCheckVSphereVirtualMachineConfig_linkedNamedSnapshot,
					locationOpt,
					label,
					datastoreOpt,
					template,
					snapshot,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVSphereVirtualMachineExists("vsphere_virtual_machine.linked_named_snapshot", &vm),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "name", "terraform-test-linked-named-snapshot"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "disk.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "disk.0.template.label", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "disk.0.template.linked", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "disk.0.template.snapshot", snapshot),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_named_snapshot", "network_interface.0.label", label),
				),
			},
		},
	})
}

// TestAccVSphereVirtualMachine_linkedCurrentSnapshotWithExtraDisk tests if it works to create a linked clone from a current snapshot and a 1 GB additional disk.
func TestAccVSphereVirtualMachine_linkedCurrentSnapshotWithExtraDisk(t *testing.T) {
	var vm virtualMachine
	var locationOpt string
	var datastoreOpt string

	if v := os.Getenv("VSPHERE_DATACENTER"); v != "" {
		locationOpt += fmt.Sprintf("    datacenter = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_CLUSTER"); v != "" {
		locationOpt += fmt.Sprintf("    cluster = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_RESOURCE_POOL"); v != "" {
		locationOpt += fmt.Sprintf("    resource_pool = \"%s\"\n", v)
	}
	if v := os.Getenv("VSPHERE_DATASTORE"); v != "" {
		datastoreOpt = fmt.Sprintf("        datastore = \"%s\"\n", v)
	}
	template := os.Getenv("VSPHERE_TEMPLATE")
	label := os.Getenv("VSPHERE_NETWORK_LABEL_DHCP")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVSphereVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccCheckVSphereVirtualMachineConfig_linkedCurrentSnapshotWithExtraDisk,
					locationOpt,
					label,
					datastoreOpt,
					template,
					datastoreOpt,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVSphereVirtualMachineExists("vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", &vm),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "name", "terraform-test-linked-current-snapshot-with-extra-disk"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "disk.#", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "disk.0.template.label", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "disk.0.template.linked", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "disk.1.size", "30"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "disk.1.type", "thin"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.linked_current_snapshot_with_extra_disk", "network_interface.0.label", label),
=======
					testAccCheckVSphereVirtualMachineExists("vsphere_virtual_machine.ipv4ipv6", &vm),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "name", "terraform-test-ipv4-ipv6"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "vcpu", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "memory", "4096"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "disk.#", "2"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "disk.0.template", template),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "network_interface.#", "1"),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "network_interface.0.label", label),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "network_interface.0.ipv4_address", ipv4Address),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "network_interface.0.ipv4_gateway", ipv4Gateway),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "network_interface.0.ipv6_address", ipv6Address),
					resource.TestCheckResourceAttr(
						"vsphere_virtual_machine.ipv4ipv6", "network_interface.0.ipv6_gateway", ipv6Gateway),
>>>>>>> provider/vsphere: IPv6 support.
				),
			},
		},
	})
}

func testAccCheckVSphereVirtualMachineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*govmomi.Client)
	finder := find.NewFinder(client.Client, true)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vsphere_virtual_machine" {
			continue
		}

		dc, err := finder.Datacenter(context.TODO(), rs.Primary.Attributes["datacenter"])
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		dcFolders, err := dc.Folders(context.TODO())
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		folder := dcFolders.VmFolder
		if len(rs.Primary.Attributes["folder"]) > 0 {
			si := object.NewSearchIndex(client.Client)
			folderRef, err := si.FindByInventoryPath(
				context.TODO(), fmt.Sprintf("%v/vm/%v", rs.Primary.Attributes["datacenter"], rs.Primary.Attributes["folder"]))
			if err != nil {
				return err
			} else if folderRef != nil {
				folder = folderRef.(*object.Folder)
			}
		}

		_, err = object.NewSearchIndex(client.Client).FindChild(context.TODO(), folder, rs.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckVSphereVirtualMachineExistsHasCustomConfig(n string, vm *virtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*govmomi.Client)
		finder := find.NewFinder(client.Client, true)

		dc, err := finder.Datacenter(context.TODO(), rs.Primary.Attributes["datacenter"])
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		dcFolders, err := dc.Folders(context.TODO())
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		_, err = object.NewSearchIndex(client.Client).FindChild(context.TODO(), dcFolders.VmFolder, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		finder = finder.SetDatacenter(dc)
		instance, err := finder.VirtualMachine(context.TODO(), rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		var mvm mo.VirtualMachine

		collector := property.DefaultCollector(client.Client)

		if err := collector.RetrieveOne(context.TODO(), instance.Reference(), []string{"config.extraConfig"}, &mvm); err != nil {
			return fmt.Errorf("error %s", err)
		}

		var configMap = make(map[string]types.AnyType)
		if mvm.Config != nil && mvm.Config.ExtraConfig != nil && len(mvm.Config.ExtraConfig) > 0 {
			for _, v := range mvm.Config.ExtraConfig {
				value := v.GetOptionValue()
				configMap[value.Key] = value.Value
			}
		} else {
			return fmt.Errorf("error no ExtraConfig")
		}

		if configMap["foo"] == nil {
			return fmt.Errorf("error no ExtraConfig for 'foo'")
		}

		if configMap["foo"] != "bar" {
			return fmt.Errorf("error ExtraConfig 'foo' != bar")
		}

		if configMap["car"] == nil {
			return fmt.Errorf("error no ExtraConfig for 'car'")
		}

		if configMap["car"] != "ferrari" {
			return fmt.Errorf("error ExtraConfig 'car' != ferrari")
		}

		if configMap["num"] == nil {
			return fmt.Errorf("error no ExtraConfig for 'num'")
		}

		// todo this should be an int, getting back a string
		if configMap["num"] != "42" {
			return fmt.Errorf("error ExtraConfig 'num' != 42")
		}
		*vm = virtualMachine{
			name: rs.Primary.ID,
		}

		return nil
	}
}

func testAccCheckVSphereVirtualMachineExists(n string, vm *virtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*govmomi.Client)
		finder := find.NewFinder(client.Client, true)

		dc, err := finder.Datacenter(context.TODO(), rs.Primary.Attributes["datacenter"])
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		dcFolders, err := dc.Folders(context.TODO())
		if err != nil {
			return fmt.Errorf("error %s", err)
		}

		folder := dcFolders.VmFolder
		if len(rs.Primary.Attributes["folder"]) > 0 {
			si := object.NewSearchIndex(client.Client)
			folderRef, err := si.FindByInventoryPath(
				context.TODO(), fmt.Sprintf("%v/vm/%v", rs.Primary.Attributes["datacenter"], rs.Primary.Attributes["folder"]))
			if err != nil {
				return err
			} else if folderRef != nil {
				folder = folderRef.(*object.Folder)
			}
		}

		_, err = object.NewSearchIndex(client.Client).FindChild(context.TODO(), folder, rs.Primary.Attributes["name"])

		*vm = virtualMachine{
			name: rs.Primary.ID,
		}

		return nil
	}
}

const testAccCheckVSphereVirtualMachineConfig_basic = `
resource "vsphere_virtual_machine" "foo" {
    name = "terraform-test"
%s
    vcpu = 2
    memory = 4096
    gateway = "%s"
    network_interface {
        label = "%s"
        ipv4_address = "%s"
        ipv4_prefix_length = 24
    }
    disk {
%s
        template {
        	label = "%s"
        }
        iops = 500
    }
    disk {
        size = 1
        iops = 500
    }
}
`
const testAccCheckVSphereVirtualMachineConfig_dhcp = `
resource "vsphere_virtual_machine" "bar" {
    name = "terraform-test"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
    }
    disk {
%s
        template {
        	label = "%s"
        }
    }
}
`

const testAccCheckVSphereVirtualMachineConfig_custom_configs = `
resource "vsphere_virtual_machine" "car" {
    name = "terraform-test-custom"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
    }
    custom_configuration_parameters {
	"foo" = "bar"
	"car" = "ferrari"
	"num" = 42
    }
    disk {
%s
        template {
        	label = "%s"
        }
    }
}
`

const testAccCheckVSphereVirtualMachineConfig_createInFolder = `
resource "vsphere_virtual_machine" "folder" {
    name = "terraform-test-folder"
    folder = "%s"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
    }
    disk {
%s
        template {
        	label = "%s"
        }
    }
}
`

const testAccCheckVSphereVirtualMachineConfig_createWithFolder = `
resource "vsphere_folder" "with_folder" {
	path = "%s"
%s
}
resource "vsphere_virtual_machine" "with_folder" {
    name = "terraform-test-with-folder"
    folder = "${vsphere_folder.with_folder.path}"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
    }
    disk {
%s
        template {
        	label = "%s"
        }
    }
}
`

const testAccCheckVSphereVirtualMachineConfig_linkedCurrentSnapshot = `
resource "vsphere_virtual_machine" "linked_current_snapshot" {
    name = "terraform-test-linked-current-snapshot"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
    }

    disk {
%s
        template {
        	label = "%s"
        	linked = true
        }
    }
}
`

const testAccCheckVSphereVirtualMachineConfig_linkedNamedSnapshot = `
resource "vsphere_virtual_machine" "linked_named_snapshot" {
    name = "terraform-test-linked-named-snapshot"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
    }

    disk {
%s
        template {
        	label = "%s"
        	linked = true
        	snapshot = "%s"
        }
    }
}
`

const testAccCheckVSphereVirtualMachineConfig_linkedCurrentSnapshotWithExtraDisk = `
resource "vsphere_virtual_machine" "linked_current_snapshot_with_extra_disk" {
    name = "terraform-test-linked-current-snapshot-with-extra-disk"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
    }

    disk {
%s
        template {
        	label = "%s"
        	linked = true
        }
    }
    disk {
%s
		size = 30
		type = "thin"
    }
}
`

const testAccCheckVSphereVirtualMachineConfig_ipv4Andipv6 = `
resource "vsphere_virtual_machine" "ipv4ipv6" {
    name = "terraform-test-ipv4-ipv6"
%s
    vcpu = 2
    memory = 4096
    network_interface {
        label = "%s"
        ipv4_address = "%s"
        ipv4_prefix_length = 24
        ipv4_gateway = "%s"
        ipv6_address = "%s"
        ipv6_prefix_length = 64
        ipv6_gateway = "%s"
    }
    disk {
%s
        template = "%s"
        iops = 500
    }
    disk {
        size = 1
        iops = 500
    }
}
`
