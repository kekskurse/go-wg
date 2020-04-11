package wg

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

type DeviceInterface interface {
	Remove() (err error)
	Add() (err error)
	CheckExists() (status bool, err error)
	AddIPRange(iprange string) (err error)
	IPRangeExists(iprange string) (status bool, err error)
	GetName() (name string)
	Up () (err error)
}

type Device struct {
	Name string
	Shell WGShellInterface
}

func (d Device) Up() (err error) {
	cmd := d.Shell.Command("ip", "link", "set", d.Name, "up")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(out)
	}
	return
}

func (d Device) Remove() (err error) {
	log.Println("Try to remove device")
	cmd := d.Shell.Command("ip", "link", "delete", d.Name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Delete Device failed")
		if string(out) == "Cannot find device \"" + d.Name +"\"\n" {
			err = errors.New("Device not exists")
		}
	}
	log.Println("Device remove ok")
	return
}

func (d Device) Add() (err error) {
	log.Println("Try to add Device")
	cmd := d.Shell.Command("ip", "link", "add", "dev", d.Name, "type", "wireguard")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Cant create Device, stdout: "+ string(out))
	}
	log.Println("Devices Added OK")
	return
}

func (d Device) CheckExists() (status bool, err error) {
	log.Println("Check if Device exists")
	status = false
	if d.Name == "" {
		err = errors.New("No Device Name set")
		return
	}

	//ip link list | grep ": go-wg:" | wc -l
	cmd := d.Shell.Command("ip", "link", "list")
	out, err := cmd.CombinedOutput()
	lines := strings.Split(string(out),"\n")
	for _, v := range lines {
		matched, _ := regexp.MatchString(`^\d{1,5}:\s`+d.Name+`:`, v)
		if matched {
			status = true
			log.Println("Device exist")
			return
		}
	}

	log.Println("Device dont exist")

	return
}

func (d Device) AddIPRange(iprange string) (err error) {
	log.Println("Try to add IP Range >"+iprange+"< to device")
	cmd := d.Shell.Command("ip", "addr", "add", iprange, "dev", d.Name)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("Cant create Device, stdout: "+ string(out))
		return
	}

	log.Println("IP Ranged added")

	return
}

func (d Device) IPRangeExists(iprange string) (status bool, err error) {
	log.Println("Check if IP-Range already exists")
	status = false
	cmd := d.Shell.Command("ip", "addr", "show", "dev", d.Name)
	out, err := cmd.CombinedOutput()

	lines := strings.Split(string(out),"\n")
	for _, v := range lines {
		matched, _ := regexp.MatchString(`inet ` + iprange + ` scope global `+d.Name, v)
		if matched {
			status = true
			log.Println("IP Range found")
			return
		}
	}
	log.Println("IP Range not found")
	return
}

func (d Device) GetName() (name string) {
	name = d.Name
	return
}