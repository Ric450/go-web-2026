package main

import (
    "fmt"
    "net/http"
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
        id := r.PathValue("id")
        fmt.Fprintln(w, "Jeu numero", id)
    })

    http.ListenAndServe(":8080", mux)
}