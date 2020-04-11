package main

import (
	"database/sql"
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"wg"
)

var wireguard wg.Wireguard

func setupWireguard(delete bool) {
	shell := wg.WGShell{} //Just wrapper for exec.Command
	device := wg.Device{}
	device.Name = "wg-test"
	device.Shell = shell

	if delete {
		if l, _ := device.CheckExists(); l == true {
			device.Remove()
		}
	}

	if l, _ := device.CheckExists(); l == false {
		err := device.Add()
		if err != nil {
			panic(err)
		}
	}

	if l, _ := device.IPRangeExists(c.IPRange); l == false {
		err := device.AddIPRange(c.IPRange)
		if err != nil {
			panic(err)
		}
	}

	wireguard = wg.Wireguard{}
	wireguard.Device = device
	wireguard.Path = c.ServerCertificatePath
	wireguard.Shell = shell
	wireguard.SetupFolder()

	if l, _ := wireguard.PrivateKeyExists(); l == false {
		wireguard.GeneratePrivateKey()
	}

	if l, _ := wireguard.PublicKeyExists(); l == false {
		wireguard.GeneratePublicKey()
	}

	port, err := wireguard.GetListenPort()
	if err != nil {
		if err.Error() != "No Listen Port found" {
			panic(err)
		}
	}
	if port == 0 {
		wireguard.SetListenPort(c.ListenPort)
	}
	wireguard.AddPrivateKey()
	wireguard.Device.Up()
}

type dbClient struct {
	id int
	status string
	publicKey string
	internIpv4 string
}

type config struct {
	DBConnectionString string `yaml:"DBConnectionString"`
	ListenPort int `yaml:"listenPort"`
	IPRange string `yaml:"ipRange"`
	ServerCertificatePath string `yaml:"serverCertificatePath"`
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

//Get clients with status approved from the DB and add them to wireguard if they not exists
func addClientsFromDB()  {
	db, err := sql.Open("mysql", c.DBConnectionString)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	results, err := db.Query("SELECT id, status, publicKey, internIpv4  FROM tickets WHERE `status` = \"approved\"")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var clients []dbClient

	for results.Next() {
		var client dbClient
		// for each row, scan the result into our tag composite object
		err = results.Scan(&client.id, &client.status, &client.publicKey, &client.internIpv4)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(client.publicKey)
		clients = append(clients, client)
	}

	for _, client := range clients {
		if l, _ := wireguard.ClientExists(client.publicKey); l == false {
			log.Println("Add new Client: ", client.publicKey)
			wireguard.AddClient(client.publicKey, client.internIpv4)
		}
	}
}

// Check if clients are connected to the wireguard which are not with status approved in DB
func checkWireguardClientsAgainsDB() {
	db, err := sql.Open("mysql", c.DBConnectionString)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	results, err := db.Query("SELECT id, status, publicKey, internIpv4  FROM tickets WHERE `status` = \"approved\"")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var clients []dbClient

	for results.Next() {
		var client dbClient
		// for each row, scan the result into our tag composite object
		err = results.Scan(&client.id, &client.status, &client.publicKey, &client.internIpv4)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(client.publicKey)
		clients = append(clients, client)
	}

	wgClients, err := wireguard.ListClients()

	if err != nil {
		panic(err)
	}

	for _, wgClient := range wgClients {
		exists := false
		for _, dbClient := range clients {
			if wgClient == dbClient.publicKey {
				exists = true
			}
		}

		if exists == false {
			log.Println("Need to remove client ", wgClient)
			wireguard.RemoveClient(wgClient)
		}
	}
}

func main() {
	log.Println("Start")
	log.Println("Basic Wireguard Settings")
	configFile := flag.String("config", "/etc/go-wg/server.yaml", "Path to a config yaml file")
	delete := flag.String("delete", "no", "if yes it will reset all")
	flag.Parse()

	del := false
	if *delete == "yes" {
		del = true
	}

	log.Println(*delete)
	log.Println("Delete Device before create", del)

	readConfig(*configFile)
	setupWireguard(del)
	log.Println("Get Clients from DB")
	addClientsFromDB()
	log.Println("Remove Clients not exists in DB")
	checkWireguardClientsAgainsDB()
	log.Println("Done")

}