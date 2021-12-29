// Package firefly makes requests to the Firefly-III API
package firefly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type rawCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (f *Firefly) Categories() ([]Category, error) {
	const path = "/api/v1/autocomplete/categories"

	req, _ := http.NewRequest("GET", f.url+path, nil)
	req.Header.Add("Authorization", "Bearer "+f.token)
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Categories: %s", err)
	}
	defer resp.Body.Close()

	var rawResults []rawCategory
	json.NewDecoder(resp.Body).Decode(&rawResults)

	var results []Category

	for _, r := range rawResults {
		id, err := strconv.Atoi(r.ID)
		if err != nil {
			return nil, fmt.Errorf("could not convert id to int: %s", err)
		}
		c := Category{ID: id, Name: r.Name}
		results = append(results, c)
	}

	return results, nil
}
