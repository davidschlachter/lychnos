// Package firefly makes requests to the Firefly-III API
package firefly

import (
	"fmt"
	"net/http"
)

type Firefly struct {
	client     *http.Client
	token, url string
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
