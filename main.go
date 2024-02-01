package main

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Définition de la structure Sheep

type Sheep struct {
	ID     int
	Name   string
	Age    int
	Weight float64
}

// Fonction pour créer une nouvelle instance de Sheep
func NewSheep(id int, name string, age int, weight float64) *Sheep {
	return &Sheep{
		ID:     id,
		Name:   name,
		Age:    age,
		Weight: weight,
	}
}

// Fonction pour afficher les informations d'un Sheep
func PrintSheep(sheep *Sheep) {
	fmt.Printf("Nom : %s\n", sheep.Name)
	fmt.Printf("Âge : %d\n", sheep.Age)
	fmt.Printf("Poids : %.2f\n", sheep.Weight)
}

// Fonction pour compter le nombre de Sheep dans une ferme
func CountSheeps(farm []*Sheep) int {
	return len(farm)
}

// Fonction pour ajouter un Sheep à la ferme
func AddSheep(farm map[string]*Sheep, sheep *Sheep) {
	farm[sheep.Name] = sheep
}

// Fonction pour supprimer un Sheep de la ferme
func RemoveSheep(farm map[string]*Sheep, name string) {
	delete(farm, name)
}

// Fonction pour obtenir un Sheep de la ferme par son nom
func GetSheep(farm map[string]*Sheep, name string) *Sheep {
	return farm[name]
}

// Création de la ferme de Sheep
var farm = make(map[string]*Sheep)

func CreateSheepHandler(w http.ResponseWriter, r *http.Request) {
	var sheep Sheep
	err := json.NewDecoder(r.Body).Decode(&sheep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	AddSheep(farm, &sheep)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Mouton créé avec succès")
}

// Handler pour mettre à jour un Sheep
func UpdateSheepHandler(w http.ResponseWriter, r *http.Request) {
	var sheep Sheep
	err := json.NewDecoder(r.Body).Decode(&sheep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	RemoveSheep(farm, sheep.Name)
	AddSheep(farm, &sheep)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mouton mis à jour avec succès")
}

// Handler pour supprimer un Sheep
func DeleteSheepHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	RemoveSheep(farm, name)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mouton supprimé avec succès")
}

// Handler pour obtenir un Sheep
func GetSheepHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("Nom")

	sheep := GetSheep(farm, name)
	if sheep == nil {
		http.Error(w, "Mouton introuvable", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(sheep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/sheep", CreateSheepHandler)
	r.Get("/sheep", UpdateSheepHandler)
	r.Get("/sheep", DeleteSheepHandler)
	r.Get("/sheep", GetSheepHandler)

	http.ListenAndServe(":45002", r)
	
}
