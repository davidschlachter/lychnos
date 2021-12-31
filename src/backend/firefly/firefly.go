// Package firefly makes requests to the Firefly-III API
package firefly

import (
	"fmt"
	"net/http"
)

type Firefly struct {
	client     *http.Client
	token, url string
	cache      Cache
}

func New(client *http.Client, token, url string) (*Firefly, error) {
	if len(token) == 0 || len(url) == 0 || client == nil {
		return nil, fmt.Errorf("must provide valid client, token and url")
	}
	return &Firefly{
		client: client,
		token:  token,
		url:    url,
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
