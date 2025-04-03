// Package firefly makes requests to the Firefly-III API
package firefly

import (
	"fmt"
	"net/http"
)

type Config struct {
	Token, URL string
}

type Firefly struct {
	client *http.Client
	config Config
	cache  Cache
}

func New(client *http.Client, c Config) (*Firefly, error) {
	if len(c.Token) == 0 || len(c.URL) == 0 || client == nil {
		return nil, fmt.Errorf("must provide valid client, token and url")
	}
	return &Firefly{
		client: client,
		config: c,
	}, nil
}

type meta struct {
	Pagination pagination `json:"pagination"`
}

type pagination struct {
	Total       int `json:"total"`
	Count       int `json:"count"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

type links struct {
	Self  string `json:"self"`
	First string `json:"first"`
	Next  string `json:"next"`
	Last  string `json:"last"`
}
