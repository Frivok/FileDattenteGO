package main

import (
	"fmt"
	"sync"
)

type FonctionProducteurDeStockage func() Stockeur

type Stockeur interface {
	Pousser([]byte) (int, error)
	Obtenir(int) ([]byte, error)
}

// MemoryStore est une implémentation de l'interface Stockeur qui stocke les données en mémoire.
type MemoryStore struct {
	mu   sync.RWMutex
	data [][]byte
}

// NouveauMemoryStore crée une nouvelle instance de MemoryStore.
func NouveauMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make([][]byte, 0),
	}
}

// Pousser ajoute un tableau de bytes à la mémoire et renvoie son offset.
func (s *MemoryStore) Pousser(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, b)
	return len(s.data) - 1, nil
}

// Obtenir récupère le tableau de bytes à l'offset spécifié.
func (s *MemoryStore) Obtenir(offset int) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset < 0 {
		return nil, fmt.Errorf("l'offset ne peut pas être inférieur à 0")
	}
	if len(s.data)-1 < offset {
		return nil, fmt.Errorf("offset (%d) trop élevé", offset)
	}
	return s.data[offset], nil
}
