package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

type api struct {
 TicketRepository TicketRepository
}

func (a api) createRoute(r chi.Router) {
	r.Get("/ticket", a.listTickets)
}

func (a api) createPublicRoute(r chi.Router) {
	r.Post("/ticket", a.createTicket)
	r.Get("/ticket/{ticketID}", a.getPublicTicket)
}

func (a api) listTickets(w http.ResponseWriter, r *http.Request) {
	tickets := a.TicketRepository.List()

	ren := render.New()
	ren.JSON(w, http.StatusOK, tickets)
}

type apiResponse struct {
	Ticket PublicTicket `json:"ticket"`
	Server []Server `json:"servers"`
}

func (a api) getPublicTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketID")
	ticket := a.TicketRepository.GetByTicketID(ticketID)
	keys, _ := r.URL.Query()["publicKey"]

	if len(keys) != 1 {
		http.Error(w, http.StatusText(400)+": Public Key is missing", 400)
		return
	}

	if ticket.PublicKey != keys[0] {
		http.Error(w, http.StatusText(400)+": Public Key dont match", 400)
		return
	}

	publicTicket := ticket.PublicTicket()

	res := apiResponse{}
	res.Ticket = publicTicket
	res.Server = GetServer()

	ren := render.New()
	ren.JSON(w, http.StatusOK, res)
}

func (a api) createTicket(w http.ResponseWriter, r *http.Request) {
	ticketRequest := TicketRequest{}
	err := json.NewDecoder(r.Body).Decode(&ticketRequest)
	if err != nil {
		panic(err)
	}

	if ticketRequest.PublicKey == "" {
		http.Error(w, http.StatusText(400) + ": Empty Public Key", 400)
		return
	}

	if ticketRequest.Hostname == "" {
		http.Error(w, http.StatusText(400) + ": Empty Hostname", 400)
		return
	}

	ticket := Ticket{}
	ticket.Status = "new"
	ticket.Hostname = ticketRequest.Hostname
	ticket.PublicKey = ticketRequest.PublicKey
	ticket.TicketID = uuid.Must(uuid.NewV4()).String()
	ticket.PublicIP = r.RemoteAddr

	ticket, err = a.TicketRepository.CreateTicket(ticket)

	publicTicket := ticket.PublicTicket()

	if err != nil {
		panic(err)
	}

	ren := render.New()
	ren.JSON(w, http.StatusOK, publicTicket)
}
