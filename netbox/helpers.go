package netbox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

func Doit(prefixId, vcpus, memory, disk int, dnsName, description, cluster, platform, comments, interfaceName, interfaceDescription string) (*VirtVm, error) {
	prefix, err := getNextIpFromPrefix(prefixId)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	ipId, err := reserveIp(prefix.Address, dnsName, description)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	clusterId, err := getClusterByName(cluster)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	platformId, err := getPlatformByName(platform)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	vmId, err := createVirtVm(dnsName, comments, *platformId, *clusterId, vcpus, memory, disk)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	interfaceId, err := addInterfaceToVm(interfaceName, interfaceDescription, *vmId)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	err = assignIpVmInterface(*interfaceId, *ipId)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	vm, err := updateVmWithPrimaryIp(*vmId, *ipId)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}

	return vm, nil
}

func buildClient(skipVerify bool) (*http.Client, error) {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	return client, nil
}

func getNextIpFromPrefix(prefix int) (*IpamPrefix, error) {

	client, err := buildClient(true)
	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}
	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://netbox.cyber.range/api/ipam/prefixes/%d/available-ips/?limit=1", prefix), nil)

	// Headers
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failure : ", err)
		return nil, err
	}

	defer resp.Body.Close()
	prefixRes := IpamPrefix{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&prefixRes)

	return &prefixRes, nil
}

func reserveIp(address, dnsName, description string) (*int, error) {

	jsonStr := fmt.Sprintf("{\"address\": \"%s\",\"dns_name\": \"%s\",\"description\": \"%s\"}", address, dnsName, description)
	jsonReq := []byte(jsonStr)
	body := bytes.NewBuffer(jsonReq)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://netbox.cyber.range/api/ipam/ip-addresses/", body)

	// Headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Failure : %+v", err)
		return nil, err
	}

	defer resp.Body.Close()
	ipReservation := IpamIpAddress{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&ipReservation)

	return &ipReservation.ID, nil

}

func getClusterByName(cluster string) (*int, error) {

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://netbox.cyber.range/api/virtualization/clusters/?name=%s", cluster), nil)

	// Headers
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	defer resp.Body.Close()
	res := ClusterQueryResults{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&res)

	return &res.Results[0].ID, nil
}

func getPlatformByName(platform string) (*int, error) {

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://netbox.cyber.range/api/dcim/platforms/?name=%s", platform), nil)

	// Headers
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
		return nil, err
	}

	defer resp.Body.Close()
	res := DcimPlatformQueryResults{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&res)

	return &res.Results[0].ID, nil

}

func createVirtVm(name, comments string, platform, cluster, vcpus, memory, disk int) (*int, error) {

	vmReq := VirtVmReq{
		Name:     name,
		Platform: platform,
		Cluster:  cluster,
		Vcpus:    vcpus,
		Memory:   memory,
		Disk:     disk,
		Comments: comments,
	}

	vmReqByte, err := json.Marshal(vmReq)
	if err != nil {
		fmt.Printf("Failure : %s", err)
		return nil, err
	}

	// Create client
	client := &http.Client{}

	body := bytes.NewBuffer(vmReqByte)
	// Create request
	req, err := http.NewRequest("POST", "https://netbox.cyber.range/api/virtualization/virtual-machines/", body)

	// Headers
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}
	defer resp.Body.Close()
	res := VirtVm{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&res)

	return &res.ID, nil

}

func addInterfaceToVm(name, description string, vmid int) (*int, error) {

	jsonStr := fmt.Sprintf("{\"name\": \"%s\",\"description\": \"%s\",\"virtual_machine\": \"%d\"}", name, description, vmid)
	body := bytes.NewBuffer([]byte(jsonStr))

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://netbox.cyber.range/api/virtualization/interfaces/", body)

	// Headers
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
		return nil, err
	}

	defer resp.Body.Close()
	res := VirtInterface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&res)

	return &res.ID, nil

}

func assignIpVmInterface(intId, ipId int) error {

	reqStr := fmt.Sprintf("{\"interface\": \"%d\"}", intId)
	body := bytes.NewBuffer([]byte(reqStr))

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://netbox.cyber.range/api/ipam/ip-addresses/%d/", ipId), body)

	// Headers
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
		return err
	}
	defer resp.Body.Close()
	res := IpamIpAddress{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&res)

	return nil

}

func updateVmWithPrimaryIp(vmId, ipId int) (*VirtVm, error) {

	reqStr := fmt.Sprintf("{\"primary_ip\": \"%d\",\"primary_ip4\": \"%d\"}", ipId, ipId)
	body := bytes.NewBuffer([]byte(reqStr))

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://netbox.cyber.range/api/virtualization/virtual-machines/%d/", vmId), body)

	// Headers
	req.Header.Add("Authorization", "Token 87196da72a4b38b23e6ac3273291d089e73ee59e")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
		return nil, err
	}
	defer resp.Body.Close()
	res := VirtVm{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&res)

	return &res, nil
}
