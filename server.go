package main

import (
	"fmt"
	"log/slog"
)

type Message struct {
	Topic string
	Data  []byte
}

type Configuration struct {
	AdresseEcoute                string
	FonctionProducteurDeStockage FonctionProducteurDeStockage
}

type Serveur struct {
	*Configuration

	// topics est une carte (map) qui associe des sujets (topics) à des stockeurs (storers).
	topics map[string]Stockeur

	consumers       []Consommateur
	producteurs     []Producteur
	canalProduction chan Message
	canalQuitter    chan struct{}
}

func NouveauServeur(cfg *Configuration) (*Serveur, error) {
	canalProduction := make(chan Message)
	return &Serveur{
		Configuration:   cfg,
		topics:          make(map[string]Stockeur),
		canalQuitter:    make(chan struct{}),
		canalProduction: canalProduction,
		producteurs: []Producteur{
			NouveauProducteurHTTP(cfg.AdresseEcoute, canalProduction),
		},
	}, nil
}

func (s *Serveur) Demarrer() {
	for _, producteur := range s.producteurs {
		go func(p Producteur) {
			if err := p.Demarrer(); err != nil {
				fmt.Println(err)
			}
		}(producteur)
	}
	s.boucle()
}

func (s *Serveur) boucle() {
	for {
		select {
		case <-s.canalQuitter:
			return
		case msg := <-s.canalProduction:
			offset, err := s.publier(msg)
			if err != nil {
				slog.Error("échec de la publication", err)
			} else {
				slog.Info("message publié", "offset", offset)
			}
		}
	}
}

func (s *Serveur) publier(msg Message) (int, error) {
	stockeur := s.obtenirStockeurPourSujet(msg.Topic)
	return stockeur.Pousser(msg.Data)
}

func (s *Serveur) obtenirStockeurPourSujet(sujet string) Stockeur {
	if _, ok := s.topics[sujet]; !ok {
		s.topics[sujet] = s.FonctionProducteurDeStockage()
		slog.Info("nouveau sujet créé", "sujet", sujet)
	}
	return s.topics[sujet]
}
