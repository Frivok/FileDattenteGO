package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type Producteur interface {
	Demarrer() error
}

type ProducteurHTTP struct {
	adresseEcoute    string
	serveur          *http.Server
	cancalProduction chan Message
}

func (p *ProducteurHTTP) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	var (
		chemin  = strings.TrimPrefix(r.URL.Path, "/")
		parties = strings.Split(chemin, "/")
	)

	if r.Method == "GET" {
		// Traitement GET
	}

	if r.Method == "POST" {
		if len(parties) != 2 {
			fmt.Println("action invalide")
			return
		}
		p.cancalProduction <- Message{
			Data:  []byte("test"),
			Topic: parties[1],
		}
	}

	fmt.Println(parties)
}

func NouveauProducteurHTTP(adresseEcoute string, canalProduction chan Message) *ProducteurHTTP {
	return &ProducteurHTTP{
		adresseEcoute:    adresseEcoute,
		cancalProduction: canalProduction,
	}
}

func (p *ProducteurHTTP) Demarrer() error {
	slog.Info("Transport HTTP démarré", "port", p.adresseEcoute)
	return http.ListenAndServe(p.adresseEcoute, p)
}
