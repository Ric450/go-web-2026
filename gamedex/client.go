package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

var client = &http.Client{
    Timeout: 3 * time.Second,
}
 
func FetchGames() ([]Game, error) {
    resp, err := client.Get("http://localhost:8081/api/games")

    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("erreur HTTP : %d", resp.StatusCode)
    }

    var games []Game

    err = json.NewDecoder(resp.Body).Decode(&games)

    if err != nil {
        return nil, err
    }

    return games, nil
}

func FetchTournaments() ([]Tournament, error) {
    resp, err := client.Get("http://localhost:8081/api/tournaments")

    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("erreur HTTP : %d", resp.StatusCode)
    }

    var tournaments []Tournament

    err = json.NewDecoder(resp.Body).Decode(&tournaments)

    if err != nil {
        return nil, err
    }

    return tournaments, nil
}