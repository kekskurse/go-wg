package main

import (
	"log"
	"math/rand"
	"net"
)

type IPHandler struct {
	allFreeIPs []string
	Config config
}

func newIPHandler(c config) *IPHandler {
	handler := new(IPHandler)
	handler.Config = c
	return handler
}

func (ipHandler *IPHandler) Setup (repository TicketRepository) {
	var err error
	log.Println("IP Range", ipHandler.Config.IPRange)
	ipHandler.allFreeIPs, err = ipHandler.getIPs(ipHandler.Config.IPRange)
	if err != nil {
		panic(err)
	}

	tickets := repository.List()

	for _, ticket := range(tickets) {
		if ticket.InternIpv4 != "" {
			ipHandler.removeIpFromList(ticket.InternIpv4)
		}
	}
}


func (ipHandler *IPHandler) getIPs(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); ipHandler.inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}
func (ipHandler *IPHandler) inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
func (ipHandler *IPHandler) removeIpFromList(ipToRemove string) ([]string) {
	var index int
	for tindex, element := range ipHandler.allFreeIPs {
		if element == ipToRemove {
			index = tindex
			break
		}
	}

	return append(ipHandler.allFreeIPs[:index], ipHandler.allFreeIPs[index+1:]...)
}

func (ipHandler *IPHandler) GetRandomFreeIP() (ip string) {
	log.Println("Get Random IP")
	index := rand.Intn(len(ipHandler.allFreeIPs))
	ip = ipHandler.allFreeIPs[index]
	ipHandler.removeIpFromList(ip)

	log.Println("Random IP: ", ip)

	return
}