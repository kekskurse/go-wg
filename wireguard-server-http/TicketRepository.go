package main

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)


type TicketRepositoryInterface interface {
	List() (tickets []Ticket)
	GetByTicketID(ticketID string) (ticket Ticket)
	CreateTicket (ticket Ticket) (retTicket Ticket, err error)
}

type Ticket struct {
	ID int64 `json:"id"`
	TicketID string `json:"ticketID"`
	Status string `json:"status"`
	PublicKey string `json:"publicKey"`
	PublicIP string `json:"publicIP"`
	Hostname string `json:"hostname"`
	InternIpv4 string `json:"internIpv4"`
}

type PublicTicket struct {
	TicketID string `json:"ticketID"`
	Status string `json:"status"`
	InternIpv4 string `json:"internIpv4"`
}

type TicketRequest struct {
	PublicKey string `json:"publicKey"`
	Hostname string `json:"hostname"`
}

type Server struct {
	Hostname string `json:"hostname"`
	AllowdIP string `json:"allowdIP"`
	PublicKey string `json:"publicKey"`
}

func GetServer () (s []Server) {
	server := Server{}
	server.Hostname = "vpn.n6e.de"
	server.AllowdIP = "10.42.0.0/16"
	server.PublicKey = "abc"
	s = append(s, server)
	return
}

type TicketRepository struct{
	DB *sql.DB
	IPHandler *IPHandler
}

func (t TicketRepository) Setup() {
	t.IPHandler.Setup(t)
}

func (t TicketRepository) List () (tickets []Ticket) {
	log.Println("List Tickets")
	results, err := t.DB.Query("SELECT id, ticketID, status, publicKey, publicIP, hostname, internIpv4  FROM tickets")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var ticket Ticket
		// for each row, scan the result into our tag composite object
		err = results.Scan(&ticket.ID, &ticket.TicketID, &ticket.Status, &ticket.PublicKey, &ticket.PublicIP, &ticket.Hostname, &ticket.InternIpv4)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		tickets = append(tickets, ticket)
	}

	return
}

func (t TicketRepository) CreateTicket (ticket Ticket) (retTicket Ticket, err error) {
	retTicket = ticket
	res, err := t.DB.Exec(`INSERT INTO tickets (ticketID, status, publicKey, publicIP, hostname, internIpv4) VALUES(?, ?, ?, ?, ?, "")`, ticket.TicketID, ticket.Status, ticket.PublicKey, ticket.PublicIP, ticket.Hostname)
	if err != nil {
		return
	} else {
		var id int64
		id, err = res.LastInsertId()
		if err != nil {
			return
		}

		retTicket.ID = id
	}

	return
}

func (t TicketRepository) GetByTicketID(ticketID string) (ticket Ticket) {
	results, err := t.DB.Query("SELECT id, ticketID, status, publicKey, publicIP, hostname, internIpv4  FROM tickets WHERE ticketID = ?", ticketID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&ticket.ID, &ticket.TicketID, &ticket.Status, &ticket.PublicKey, &ticket.PublicIP, &ticket.Hostname, &ticket.InternIpv4)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		return
	}

	return
}
func (t TicketRepository) SaveTicket(ticket Ticket) {
	log.Println("Save Ticket", ticket.ID)
	_, err := t.DB.Exec(`UPDATE tickets SET status = ?, internIpv4 = ? WHERE id = ?;`, ticket.Status, ticket.InternIpv4, ticket.ID)
	if err != nil {
		panic(err)
	}
}

func (t TicketRepository) ChangeStatus(ticketID string, status string) () {
	ticket := t.GetByTicketID(ticketID)
	ticket.Status = status
	t.SaveTicket(ticket)
}
func (t TicketRepository) ActivateTicket(ticketID string) () {
	ticket := t.GetByTicketID(ticketID)
	if ticket.ID == 0 {
		panic("Ticket not found")
	}
	if ticket.Status == "approved" {
		log.Println("Nothing to do")
		return
	}
	ticket.InternIpv4 = t.IPHandler.GetRandomFreeIP()
	ticket.Status = "approved"
	t.SaveTicket(ticket)
}

func (t Ticket) PublicTicket () (ticket PublicTicket) {
	ticket.Status = t.Status
	ticket.TicketID = t.TicketID
	ticket.InternIpv4 = t.InternIpv4
	return
}
