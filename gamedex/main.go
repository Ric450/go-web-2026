package main

import (
    "flag"
    "fmt"
    "net/http"
    "strconv"
)

func main() {
    port := flag.String("port", "8080", "Port du serveur")
    flag.Parse()

    games, err := FetchGames()

    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }

    tournaments, err := FetchTournaments()

    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }

    for _, game := range games {
        count := 0

        for _, tournament := range tournaments {
            if tournament.GameID == game.ID {
                count++
            }
        }

        fmt.Println(game.Name, game.Year, "-", count, "tournois")
    }

    mux := http.NewServeMux()

    mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
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

    mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.ListenAndServe(":"+*port, mux)
}