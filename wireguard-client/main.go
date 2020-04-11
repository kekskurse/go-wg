package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/sparrc/go-ping"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"wg"
)

type config struct {
	ServerURL string `yaml:"serverURL"`
	ConfigPath  string `yaml:"clientCertificatePath"`
}

var c config

func readConfig(configFile string) {
	log.Println("Config file: ", configFile)
	yamlFile, err := ioutil.ReadFile(configFile)

	log.Println(string(yamlFile))
	log.Println(c)

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		panic(err)
	}

	log.Println(c)
}

type Ticket struct {
	TicketID string `json:"ticketID"`
	Status string `json:"status"`
	InternIpv4 string `json:"internIpv4"`
}

type Server struct {
	Hostname string
	AllowedIP string
	PublicKey string
	InternalServerIP string `json:"internalServerIP"`
}

type Response struct {
	Ticket Ticket
	Servers []Server
}

func createTicket () {
	publicKey, _ := wireguard.GetPublicKey()
	hostname, _ := os.Hostname()
	log.Println("My Public Key: "+publicKey)
	values := map[string]string{"publicKey": publicKey, "hostname": hostname}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(c.ServerURL + "public/v1/ticket", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Cant create Ticket")
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	ticket := Ticket{}
	err = json.Unmarshal(body, &ticket)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(c.ConfigPath + "/ticket", []byte(ticket.TicketID), os.ModeExclusive)
	if err != nil {
		panic(err)
	}
}

func getTicketInfo() Response {
	ticketID, err := ioutil.ReadFile(c.ConfigPath + "/ticket")
	publicKey, err := wireguard.GetPublicKey()
	publicKey = url.QueryEscape(publicKey)
	resp, err := http.Get(c.ServerURL + "public/v1/ticket/" + string(ticketID)+"?publicKey="+publicKey)
	if err != nil {
		log.Println("Cant get Ticket")
		panic(err)
	}


	body, err := ioutil.ReadAll(resp.Body)

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	return response
}

func checkTicketExists() bool {
	if _, err := os.Stat(c.ConfigPath + "/ticket"); os.IsNotExist(err) {
		return false
	}
	return true
}

var wireguard wg.Wireguard
func NewWireguard() wg.Wireguard {
	shell := wg.WGShell{} //Just wrapper for exec.Command
	device := wg.Device{}
	device.Name = "wg-client"
	device.Shell = shell

	wireguard = wg.Wireguard{}
	wireguard.Device = device
	wireguard.Path = c.ConfigPath
	wireguard.Shell = shell
	wireguard.SetupFolder()

	return wireguard
}

func checkIfDeviceExists() bool {
	status, err := wireguard.Device.CheckExists()
	if err != nil {
		panic(err)
	}
	return status
}

func checkIfConnected(response Response) bool {
	log.Println("Send Ping to Server: ", response.Servers[0].InternalServerIP)
	pinger, err := ping.NewPinger(response.Servers[0].InternalServerIP)
	if err != nil {
		panic(err)
	}
	pinger.Count = 1
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()

	if stats.PacketsRecv > 1 {
		return true
	}
	return false
}

func connect(response Response) {
	wireguard.Device.Add()
	err := wireguard.Device.AddIPRange(response.Ticket.InternIpv4+"/32")
	//err := wireguard.Device.AddIPRange(`10.42.133.1/24`)
	if err != nil {
		log.Println("Cant add iprange, fatal, rollback")
		disconnect()
		panic(err)
	}

	err = wireguard.AddPrivateKey()
	if err != nil {
		panic(err)
	}
	err = wireguard.Device.Up()
	if err != nil {
		panic(err)
	}
	err = wireguard.ConnectToRemotePeer(response.Servers[0].PublicKey, response.Servers[0].AllowedIP, response.Servers[0].Hostname)
	if err != nil {
		panic(err)
	}
}

func disconnect() {
	wireguard.Device.Remove()
}


func main() {
	log.Println("Start")
	log.Println("Basic Wireguard Settings")
	configFile := flag.String("config", "/etc/go-wg/client.yaml", "Path to a config yaml file")

	flag.Parse()
	readConfig(*configFile)
	wireguard = NewWireguard()

	if s, _ := wireguard.PrivateKeyExists(); s == false {
		wireguard.GeneratePrivateKey()
	}
	if s, _ := wireguard.PublicKeyExists(); s == false {
		wireguard.GeneratePublicKey()
	}

	if checkTicketExists() == false {
		log.Println("Create new Ticket")
		createTicket()
	} else {
		log.Println("Ticket already exists")
	}

	response := getTicketInfo()
	log.Println("Ticket", response)

	disconnect()

	log.Println("Start Loop")

	if response.Ticket.Status == "approved" {
		if checkIfDeviceExists() == false {
			connect(response)
		}

		s := checkIfConnected(response)

		if s == false {
			log.Println("Not connected")
			//disconnect()
		}
	}
}
