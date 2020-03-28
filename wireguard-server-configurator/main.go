package main

import (
	"fmt"
	"log"
)
import "wg"

func main() {
	fmt.Println("Start")
	setupSampleServer()
	log.Println("Done")
}

func setupSampleServer() {
	shell := wg.WGShell{} //Just wrapper for exec.Command
	device := wg.Device{}
	device.Name = "wg-test"
	device.Shell = shell
	if l, _ := device.CheckExists(); l == true {
		device.Remove()
	}
	device.Add()

	if l, _ := device.IPRangeExists("10.42.133.0/24"); l {
		device.AddIPRange("10.42.133.0/24")
	}

	wireguard := wg.Wireguard{}
	wireguard.Device = device
	wireguard.Path = "/tmp/wg"
	wireguard.Shell = shell
	wireguard.SetupFolder()

	if l, _ := wireguard.PrivateKeyExists(); l == false {
		wireguard.GeneratePrivateKey()
	}

	if l, _ := wireguard.PublicKeyExists(); l == false {
		wireguard.GeneratePublicKey()
	}


	wireguard.SetListenPort(52137)
	wireguard.AddClient("+a5d1tw7YQC//FEmhAPsb1PKgWw18GIVFW2Bixm2nio=", "10.42.133.23/32")
}