package wg

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type WireguardInterface interface {
	GeneratePublicKey() (err error)
	PublicKeyExists() (status bool, err error)
	GeneratePrivateKey() (err error)
	PrivateKeyExists() (status bool, err error)
	SetupFolder() (err error)
	SetListenPort(port int) (err error)
	GetListenPort() (port int, err error)
	AddClient(publickey string, ipRange string) (err error)
	ClientExists(publickey string) (status bool, err error)
	ListClients() (publicKeys []string, err error)
	RemoveClient(publickey string) (err error)
	ConnectToRemotePeer(peerPublicKey string, allowedIPs string, endpoint string) (err error)
	AddPrivateKey() (err error)
}

type Wireguard struct {
	Path string
	Device DeviceInterface
	Shell WGShellInterface
}

func (w Wireguard) ConnectToRemotePeer(peerPublicKey string, allowedIPs string, endpoint string) (err error) {
	//wg set wg0 peer fAoJ02w4ravlkLbcaiIl8bbQ6svlAZXJ3mUO3XR4u0g= allowed-ips 192.168.222.2/32 endpoint 188.xxx.xx.xx:36448 persistent-keepalive 25
	cmd := w.Shell.Command("wg", "set", w.Device.GetName(), "peer", peerPublicKey, "allowed-ips", allowedIPs, "endpoint", endpoint, "persistent-keepalive", "25")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
	}
	return
}

func (w Wireguard) AddPrivateKey() (err error) {
	cmd := w.Shell.Command("wg", "set", w.Device.GetName(), "private-key", w.Path+"/privatekey")
	_, err = cmd.CombinedOutput()
	return
}

func (w Wireguard) GeneratePrivateKey() (err error) {
	log.Println("Generate Private Key")
	cmd := w.Shell.Command("wg", "genkey")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Failed generate Priivatekey")
		return
	}
	err = ioutil.WriteFile(w.Path+"/privatekey", out, 0600)
	if err != nil {
		log.Println("Failed write Private Key to HDD")
		return
	}

	log.Println("Private Key generated")
	return
}

func (w Wireguard) PublicKeyExists() (status bool, err error) {
	log.Println("Check if Public Key exists")
	status = false
	if _, err = os.Stat(w.Path + "/publickey"); os.IsNotExist(err) {
		log.Println("Public Key does not exists")
		status = false
		return
	}
	log.Println("Public Key already exists")
	status = true
	return
}

func (w Wireguard) GeneratePublicKey() (err error) {
	log.Println("Generate Private Key")
	cmd := w.Shell.Command("wg", "pubkey")
	dat, err := ioutil.ReadFile(w.Path + "/privatekey")
	if err != nil {
		return
	}
	cmd.Stdin = bytes.NewBuffer(dat)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	err = ioutil.WriteFile(w.Path + "/publickey", out, 0600)
	if err != nil {
		return
	}

	return
}

func (w Wireguard) GetPublicKey() (key string, err error) {
	dat, err := ioutil.ReadFile(w.Path + "/publickey")
	key = string(dat)
	key = strings.TrimSpace(key)
	return
}

func (w Wireguard) PrivateKeyExists() (status bool, err error) {
	log.Println("Check if Private Key exists")
	status = false
	if _, err = os.Stat(w.Path + "/privatekey"); os.IsNotExist(err) {
		status = false
		log.Println("Private Key not exists")
		return
	}
	status = true
	log.Println("Private Key exists")
	return
}

func (w Wireguard) SetupFolder() (err error) {
	if _, err := os.Stat(w.Path); os.IsNotExist(err) {
		log.Println("Create folder")
		err = os.MkdirAll(w.Path, os.ModePerm)
	}
	return
}

func (w Wireguard) SetListenPort(port int) (err error) {
	deviceName := w.Device.GetName()
	cmd := w.Shell.Command("wg", "set", deviceName, "listen-port", strconv.Itoa(port), "private-key", w.Path + "/privatekey")
	_, err = cmd.CombinedOutput()
	return
}

func (w Wireguard) AddClient(publickey string, ipRange string) (err error) {
	deviceName := w.Device.GetName()
	cmd := w.Shell.Command("wg", "set", deviceName, "peer", publickey, "persistent-keepalive", "25", "allowed-ips", ipRange)
	_, err = cmd.CombinedOutput()

	return
}

func (w Wireguard) GetListenPort() (port int, err error) {
	deviceName := w.Device.GetName()
	cmd := w.Shell.Command("wg", "showconf", deviceName)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return
	}

	lines := strings.Split(string(out),"\n")
	for _, v := range lines {
		matched, _ := regexp.MatchString(`^ListenPort = `, v)
		if matched {
			charSlice := []rune(v)
			portString := string(charSlice[13:])
			port, err = strconv.Atoi(portString)
			return
		}
	}

	err = errors.New("No Listen Port found")

	return
}

func (w Wireguard) ListClients() (publicKeys []string, err error) {
	deviceName := w.Device.GetName()
	cmd := w.Shell.Command("wg", "show", deviceName)
	out, err := cmd.CombinedOutput()
	lines := strings.Split(string(out),"\n")
	for _, v := range lines {
		matched, _ := regexp.MatchString(`peer: `, v)
		if matched {
			charSlice := []rune(v)
			publicKey := string(charSlice[6:])
			publicKeys = append(publicKeys, publicKey)

		}
	}
	return
}

func (w Wireguard) ClientExists(publickey string) (status bool, err error) {
	log.Println("Check if client "+publickey+" exists")
	status = false
	publickeys, err := w.ListClients()
	if err != nil {
		return
	}
	for _, b := range publickeys {
		if b == publickey {
			status = true
			log.Println("Client found")
			return
		}
	}
	log.Println("Client not found")
	return
}

func (w Wireguard) RemoveClient(publickey string) (err error) {
	deviceName := w.Device.GetName()
	cmd := w.Shell.Command("wg", "set", deviceName, "peer", publickey, "remove")
	_, err = cmd.CombinedOutput()

	return
}