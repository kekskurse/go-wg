package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"wg"
)

var wireguard wg.Wireguard

func setupWireguard() {
	shell := wg.WGShell{} //Just wrapper for exec.Command
	device := wg.Device{}
	device.Name = "wg-test"
	device.Shell = shell
	if l, _ := device.CheckExists(); l == false {
		err := device.Add()
		if err != nil {
			panic(err)
		}
	}

	if l, _ := device.IPRangeExists("10.42.133.0/24"); l == false {
		err := device.AddIPRange("10.42.133.0/24")
		if err != nil {
			panic(err)
		}
	}

	wireguard = wg.Wireguard{}
	wireguard.Device = device
	wireguard.Path = "/etc/go-wg"
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
		panic(err)
	}
	if port == 0 {
		wireguard.SetListenPort(51820)
	}
}

type dbClient struct {
	id int
	status string
	publicKey string
	internIpv4 string
}

//Get clients with status approved from the DB and add them to wireguard if they not exists
func addClientsFromDB()  {
	db, err := sql.Open("mysql", "root:example@tcp(127.0.0.1:4306)/wg")
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
	db, err := sql.Open("mysql", "root:example@tcp(127.0.0.1:4306)/wg")
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
	setupWireguard()
	log.Println("Get Clients from DB")
	addClientsFromDB()
	log.Println("Remove Clients not exists in DB")
	checkWireguardClientsAgainsDB()
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


	wireguard.SetListenPort(51820)
	wireguard.AddClient("+a5d1tw7YQC//FEmhAPsb1PKgWw18GIVFW2Bixm2nio=", "10.42.133.23/32")
}