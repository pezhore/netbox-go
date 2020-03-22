package netbox

import (
	"time"
)

type Config struct {
	APIKey        string
	Endpoint      string
	SkipTLSVerify bool
}

// IpamPrefix contains the details for each Netbox Prefix
type IpamPrefix struct {
	Family  int         `json:"family"`
	Address string      `json:"address"`
	Vrf     interface{} `json:"vrf"`
}

// IpamIpAddress contains the details for an IP in Netbox
type IpamIpAddress struct {
	ID     int `json:"id"`
	Family struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	} `json:"family"`
	Address string      `json:"address"`
	Vrf     interface{} `json:"vrf"`
	Tenant  interface{} `json:"tenant"`
	Status  struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	} `json:"status"`
	Role      interface{} `json:"role"`
	Interface struct {
		ID             int         `json:"id"`
		URL            string      `json:"url"`
		Device         interface{} `json:"device"`
		VirtualMachine struct {
			ID   int    `json:"id"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"virtual_machine"`
		Name string `json:"name"`
	} `json:"interface"`
	NatInside    interface{}   `json:"nat_inside"`
	NatOutside   interface{}   `json:"nat_outside"`
	DNSName      string        `json:"dns_name"`
	Description  string        `json:"description"`
	Tags         []interface{} `json:"tags"`
	CustomFields struct {
	} `json:"custom_fields"`
	Created     string    `json:"created"`
	LastUpdated time.Time `json:"last_updated"`
}

type ClusterQueryResults struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type struct {
			ID   int    `json:"id"`
			URL  string `json:"url"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"type"`
		Group struct {
			ID   int    `json:"id"`
			URL  string `json:"url"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"group"`
		Site struct {
			ID   int    `json:"id"`
			URL  string `json:"url"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"site"`
		Comments     string        `json:"comments"`
		Tags         []interface{} `json:"tags"`
		CustomFields struct {
		} `json:"custom_fields"`
		Created             string      `json:"created"`
		LastUpdated         time.Time   `json:"last_updated"`
		DeviceCount         interface{} `json:"device_count"`
		VirtualmachineCount int         `json:"virtualmachine_count"`
	} `json:"results"`
}

// DcimPlatform contains the various platform definitions (e.g. VMware ESXi, Juniper device platform, etc)
// not to be confused with Manufacturer
type DcimPlatformQueryResults struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		ID                  int         `json:"id"`
		Name                string      `json:"name"`
		Slug                string      `json:"slug"`
		Manufacturer        interface{} `json:"manufacturer"`
		NapalmDriver        string      `json:"napalm_driver"`
		NapalmArgs          interface{} `json:"napalm_args"`
		DeviceCount         int         `json:"device_count"`
		VirtualmachineCount interface{} `json:"virtualmachine_count"`
	} `json:"results"`
}

//VirtVMReq contains the data for a Virtual Machine creation request
type VirtVmReq struct {
	Name     string `json:"name"`
	Platform int    `json:"platform"`
	Cluster  int    `json:"cluster"`
	Vcpus    int    `json:"vcpus"`
	Memory   int    `json:"memory"`
	Disk     int    `json:"disk"`
	Comments string `json:"comments"`
}

// VirtVM is the datastructure of a virtual machine
type VirtVm struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	} `json:"status"`
	Site struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"site"`
	Cluster struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"cluster"`
	Role     interface{} `json:"role"`
	Tenant   interface{} `json:"tenant"`
	Platform struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"platform"`
	PrimaryIP        interface{}   `json:"primary_ip"`
	PrimaryIP4       interface{}   `json:"primary_ip4"`
	PrimaryIP6       interface{}   `json:"primary_ip6"`
	Vcpus            int           `json:"vcpus"`
	Memory           int           `json:"memory"`
	Disk             int           `json:"disk"`
	Comments         string        `json:"comments"`
	LocalContextData interface{}   `json:"local_context_data"`
	Tags             []interface{} `json:"tags"`
	CustomFields     struct {
	} `json:"custom_fields"`
	ConfigContext struct {
	} `json:"config_context"`
	Created     string    `json:"created"`
	LastUpdated time.Time `json:"last_updated"`
}

type VirtInterface struct {
	ID             int `json:"id"`
	VirtualMachine struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"virtual_machine"`
	Name string `json:"name"`
	Type struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	} `json:"type"`
	Enabled      bool          `json:"enabled"`
	Mtu          interface{}   `json:"mtu"`
	MacAddress   interface{}   `json:"mac_address"`
	Description  string        `json:"description"`
	Mode         interface{}   `json:"mode"`
	UntaggedVlan interface{}   `json:"untagged_vlan"`
	TaggedVlans  []interface{} `json:"tagged_vlans"`
	Tags         []interface{} `json:"tags"`
}
