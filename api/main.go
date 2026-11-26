
package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	addr := flag.String("addr", ":8081", "adresse d'écoute")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/{$}", index)
	mux.HandleFunc("GET /api/docs", docs)
	mux.HandleFunc("GET /api/openapi.yaml", spec)
	mux.HandleFunc("GET /api/games", listGames)
	mux.HandleFunc("POST /api/recherche", recherche)
	mux.HandleFunc("GET /api/games/{id}", getGame)
	mux.HandleFunc("GET /api/games/{id}/tournaments", gameTournaments)
	mux.HandleFunc("GET /api/tournaments", listTournaments)
	mux.HandleFunc("GET /api/stats", stats)
	mux.HandleFunc("GET /api/broken", broken)
	mux.HandleFunc("GET /api/slow", slow)

	log.Printf("Mini API sur http://localhost%s/api/", *addr)
	log.Printf("Documentation sur http://localhost%s/api/docs", *addr)
	log.Fatal(http.ListenAndServe(*addr, logRequests(mux)))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func index(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"games":           "http://localhost:8081/api/games",
		"game":            "http://localhost:8081/api/games/{id}",
		"tournaments":     "http://localhost:8081/api/tournaments",
		"gameTournaments": "http://localhost:8081/api/games/{id}/tournaments",
		"broken":          "http://localhost:8081/api/broken (répond toujours 500)",
		"slow":            "http://localhost:8081/api/slow (répond après 5 secondes)",
	})
}

func listGames(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(r.URL.Query().Get("q"))
	genre := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year")

	result := []Game{}
	for _, g := range games {
		if q != "" && !strings.Contains(strings.ToLower(g.Name), q) {
			continue
		}
		if genre != "" && !slices.Contains(g.Genres, genre) {
			continue
		}
		if year != "" && strconv.Itoa(g.Year) != year {
			continue
		}
		result = append(result, g)
	}
	writeJSON(w, http.StatusOK, result)
}

func findGame(w http.ResponseWriter, r *http.Request) (Game, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "l'identifiant doit être un nombre")
		return Game{}, false
	}
	for _, g := range games {
		if g.ID == id {
			return g, true
		}
	}
	writeError(w, http.StatusNotFound, "jeu introuvable")
	return Game{}, false
}

func getGame(w http.ResponseWriter, r *http.Request) {
	if g, ok := findGame(w, r); ok {
		writeJSON(w, http.StatusOK, g)
	}
}

func gameTournaments(w http.ResponseWriter, r *http.Request) {
	g, ok := findGame(w, r)
	if !ok {
		return
	}
	out := []Tournament{}
	for _, t := range tournaments {
		if t.GameID == g.ID {
			out = append(out, t)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func listTournaments(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, tournaments)
}

func stats(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("X-API-Key") {
	case "prof":
		writeJSON(w, http.StatusOK, map[string]int{
			"games":       len(games),
			"tournaments": len(tournaments),
		})
	case "eleve":
		writeError(w, http.StatusForbidden, "clé reconnue, mais pas le droit de voir les statistiques")
	default:
		writeError(w, http.StatusUnauthorized, "clé absente ou inconnue : ajoutez l'en-tête X-API-Key")
	}
}

func broken(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusInternalServerError, "panne simulée")
}

func slow(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	writeJSON(w, http.StatusOK, games)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

//go:embed openapi.yaml
var openapiSpec []byte

func spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.Write(openapiSpec)
}

const docsPage = `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <title>Mini API GameDex · documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({ url: "/api/openapi.yaml", dom_id: "#swagger-ui" });
  </script>
</body>
</html>`

func docs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(docsPage))
}

// Nouvelle fonction POST
func recherche(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Query string `json:"query"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		writeError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	resultats := []Game{}

	for _, game := range games {
		if strings.Contains(
			strings.ToLower(game.Name),
			strings.ToLower(data.Query),
		) {
			resultats = append(resultats, game)
		}
	}

	writeJSON(w, http.StatusOK, resultats)
}
