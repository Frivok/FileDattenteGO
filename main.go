package main

import "fmt"

func main() {
	cfg := &Configuration{
		AdresseEcoute: ":4000",
		FonctionProducteurDeStockage: func() Stockeur {
			return NouveauMemoryStore()
		},
	}
	s, err := NouveauServeur(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Demarrer()
}
