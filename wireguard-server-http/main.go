package main

import (
	"database/sql"
	"flag"
	"github.com/go-chi/chi"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type httpInterface interface {
	createRoute(r chi.Router)
}

func setubDB() *sql.DB {
	db, err := sql.Open("mysql", "root:example@tcp(127.0.0.1:4306)/wg")

	if err != nil {
		panic(err)
	}
	return db
}


type config struct {
	DBConnectionString string `yaml:"DBConnectionString"`
	ListenPort int `yaml:"listenPort"`
	IPRange string `yaml:"ipRange"`
	ServerCertificatePath string `yaml:"serverCertificatePath"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
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

func setupRouter() chi.Router {
	router := chi.NewRouter()

	db := setubDB()

	ticketRepository := TicketRepository{}
	ticketRepository.DB = db

	a := api{}
	a.TicketRepository = ticketRepository

	g := gui{}
	g.TicketRepository = ticketRepository
	//router.Route("/api/v1", a.createRoute) //todo add api stuff later
	router.Route("/gui", g.createRoute)

	router.Route("/public/v1", a.createPublicRoute)

	router.Get("/", func (w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/gui/login", 301)
	})

	return router
}

func main() {
	//
	configFile := flag.String("config", "/etc/go-wg/server.yaml", "Path to a config yaml file")

	flag.Parse()
	readConfig(*configFile)
	r := setupRouter()
	log.Println("Run at http://localhost:3333")
	http.ListenAndServe(":3333", r)
}