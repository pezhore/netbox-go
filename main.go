package main

import (
	"fmt"
	netbox "github.com/pezhore/netbox-go"
)

func main() {
	fmt.Println("vim-go")

	v := netbox.Doit(0, 2, 2048, 20, "test.range", "this is the description", "Turing", "VMware ESXi", "this is the comments", "eth0", "this is the interface desceription"
	fmt.Printf("%+v", v)
}
