package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-session/session"
	"html/template"
	"net/http"
)

type gui struct {
TicketRepository TicketRepository
}

func (g gui) createRoute(r chi.Router) {
	r.Get("/login", g.login)
	r.Post("/login", g.performLogin)
	r.Get("/ticket", g.tickets)
	r.Get("/ticket/approve", g.approveTicket)
}

func (g gui) approveTicket(w http.ResponseWriter, r *http.Request) {
	IDs, ok := r.URL.Query()["id"]
	if ok == false {
		panic("No ID given")
	}
	g.TicketRepository.ActivateTicket(IDs[0])
	http.Redirect(w, r, "/gui/ticket", 302)
}

func (g gui) login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		panic(err)
	}
}

func (g gui) tickets(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/base.html", "templates/tickets.html")
	if err != nil {
		panic(err)
	}

	type httpSturct struct {
		Tickets []Ticket
	}

	res := httpSturct{}

	res.Tickets = g.TicketRepository.List()


	err = tmpl.ExecuteTemplate(w, "base", res)
	if err != nil {
		panic(err)
	}
}

func (g gui) performLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}

	r.ParseForm()

	if len(r.Form["username"]) != 1 {
		err = tmpl.ExecuteTemplate(w, "base_fail", nil)
		if err != nil {
			panic(err)
		}
		return
	}
	if len(r.Form["password"]) != 1 {
		err = tmpl.ExecuteTemplate(w, "base_fail", nil)
		if err != nil {
			panic(err)
		}
		return
	}
	if r.Form["username"][0] != c.Username {
		err = tmpl.ExecuteTemplate(w, "base_fail", nil)
		if err != nil {
			panic(err)
		}
		return
	}

	if r.Form["password"][0] != c.Password {
		err = tmpl.ExecuteTemplate(w, "base_fail", nil)
		if err != nil {
			panic(err)
		}
		return
	}

	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		panic(err)
	}
	store.Set("login", 1)
	err = store.Save()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/gui/ticket", 302)


	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		panic(err)
	}


}