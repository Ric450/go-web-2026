package main

type Game struct {
	ID     int      `json:"id"`
    Name   string   `json:"name"`
    Year   int      `json:"year"`
    Studio string   `json:"studio"`
    Genres []string `json:"genres"`
}

type Tournament struct {
	ID     int     `json:"id"`
    GameID int     `json:"gameId"`
    Name   string  `json:"name"`
    City   string  `json:"city"`
    Date   string  `json:"date"`
    Lat    float64 `json:"lat"`
    Lng    float64 `json:"lng"`
}