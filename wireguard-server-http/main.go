package main

import (
	"database/sql"
	"github.com/go-chi/chi"
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

func setupRouter() chi.Router {
	router := chi.NewRouter()

	db := setubDB()

	ticketRepository := TicketRepository{}
	ticketRepository.DB = db

	a := api{}
	a.TicketRepository = ticketRepository
	//router.Route("/api/v1", a.createRoute)
	router.Route("/public/v1", a.createPublicRoute)

	return router
}

func main() {
	//

	r := setupRouter()
	http.ListenAndServe(":3333", r)
}