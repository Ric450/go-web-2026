package main

import (
    "fmt"
    "net/http"
    "strconv"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Bienvenue sur GameDex")
    })

    mux.HandleFunc("GET /jeux", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Liste des jeux video")
    })

    mux.HandleFunc("GET /jeux/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))

        if err != nil {
            http.Error(w, "Identifiant invalide", 400)
            return
        }

        if id == 999 {
            http.Error(w, "Jeu introuvable", 404)
            return
        }

        fmt.Fprintln(w, "Jeu numero", id)
    })

    http.ListenAndServe(":8080", mux)
}