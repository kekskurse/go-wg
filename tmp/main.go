package main

import (
	"fmt"
	"net"
)

func main() {



	A := "10.10.10.5/24"

	ipA, netA, _ := net.ParseCIDR(A)

	fmt.Println(netA.Mask.String())
	fmt.Println(ipA.String())



}
