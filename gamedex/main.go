package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	port := flag.String("port", "8080", "Port du serveur")
	flag.Parse()

	templates, err := template.ParseGlob("templates/*.html")

	if err != nil {
		fmt.Println("Erreur templates :", err)
		return
	}

	render := func(w http.ResponseWriter, name string, data any, status int) {
		var page bytes.Buffer

		err := templates.ExecuteTemplate(&page, name, data)
		if err != nil {
			fmt.Println("Erreur template :", err)
			http.Error(w, "Erreur interne", 500)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		w.Write(page.Bytes())
	}

	afficherErreur := func(w http.ResponseWriter, code int, message string) {
		data := map[string]any{
			"Code":    code,
			"Message": message,
		}

		render(w, "erreur", data, code)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.HandleFunc("GET /jeux", func(w http.ResponseWriter, r *http.Request) {
		games, err := FetchGames()

		if err != nil {
			afficherErreur(w, 503, "API indisponible")
			return
		}

		q := r.URL.Query().Get("q")

		resultats := games[:0]

		for _, game := range games {
			if strings.Contains(
				strings.ToLower(game.Name),
				strings.ToLower(q),
			) {
				resultats = append(resultats, game)
			}
		}

		data := map[string]any{
			"Games": resultats,
			"Query": q,
		}

		render(w, "liste", data, 200)
	})

	mux.HandleFunc("GET /jeux/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil || id <= 0 {
			afficherErreur(w, 400, "Identifiant invalide")
			return
		}

		games, err := FetchGames()

		if err != nil {
			afficherErreur(w, 503, "API indisponible")
			return
		}

		for _, game := range games {
			if game.ID == id {

				tournaments, err := FetchTournaments()

				if err != nil {
					afficherErreur(w, 503, "API indisponible")
					return
				}

				resultats := tournaments[:0]

				for _, tournament := range tournaments {
					if tournament.GameID == game.ID {
						resultats = append(resultats, tournament)
					}
				}

				data := map[string]any{
					"Game":        game,
					"Tournaments": resultats,
				}

				render(w, "jeu", data, 200)
				return
			}
		}

		afficherErreur(w, 404, "Jeu introuvable")
	})

	mux.Handle("GET /static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	fmt.Println("GameDex : http://localhost:" + *port)

	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		fmt.Println("Erreur serveur :", err)
	}
}