// Package firefly makes requests to the Firefly-III API
package firefly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
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

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type rawCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (f *Firefly) Categories() ([]Category, error) {
	const path = "/api/v1/autocomplete/categories?limit=1000"

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

type CategoryTotal struct {
	Category
	Spent  decimal.Decimal `json:"spent"`
	Earned decimal.Decimal `json:"earned"`
	Start  time.Time       `json:"start"`
	End    time.Time       `json:"end"`
}

type rawCategoryTotal struct {
	ID         string        `json:"id"`
	Attributes rawAttributes `json:"attributes"`
}

type rawAttributes struct {
	Name   string     `json:"name"`
	Spent  []rawTotal `json:"spent"`
	Earned []rawTotal `json:"earned"`
}

type rawTotal struct {
	Sum string `json:"sum"`
}

func (f *Firefly) ListCategoryTotals(start, end time.Time) ([]CategoryTotal, error) {
	const path = "/api/v1/categories/"
	params := fmt.Sprintf("?start=%s&end=%s", start.Format("2006-01-02"), end.Format("2006-01-02"))

	req, _ := http.NewRequest("GET", f.url+path+params, nil)
	req.Header.Add("Authorization", "Bearer "+f.token)
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Categories: %s", err)
	}
	defer resp.Body.Close()

	var rawResults struct {
		Data []rawCategoryTotal `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&rawResults)

	var results []CategoryTotal

	for _, r := range rawResults.Data {
		if len(r.Attributes.Spent) == 0 && len(r.Attributes.Earned) == 0 {
			continue
		}
		var c CategoryTotal
		var spent, earned decimal.Decimal
		id, err := strconv.Atoi(r.ID)
		if err != nil {
			return nil, fmt.Errorf("could not convert id to int: %s", err)
		}
		// TODO(davidschlachter): Maybe support multiple currencies one day
		if len(r.Attributes.Spent) == 1 {
			spent, err = decimal.NewFromString(r.Attributes.Spent[0].Sum)
			if err != nil {
				return nil, fmt.Errorf("could not convert spent sum to decimal: %s", err)
			}
		}
		if len(r.Attributes.Earned) == 1 {
			earned, err = decimal.NewFromString(r.Attributes.Earned[0].Sum)
			if err != nil {
				return nil, fmt.Errorf("could not convert earned sum to decimal: %s", err)
			}
		}

		c.ID = id
		c.Name = r.Attributes.Name
		c.Spent = spent
		c.Earned = earned
		c.Start = start
		c.End = end

		results = append(results, c)
	}

	return results, nil
}

func (f *Firefly) FetchCategoryTotal(catID int, start, end time.Time) ([]CategoryTotal, error) {
	const path = "/api/v1/categories/"
	params := fmt.Sprintf("?start=%s&end=%s", start.Format("2006-01-02"), end.Format("2006-01-02"))

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s%d%s", f.url, path, catID, params), nil)
	req.Header.Add("Authorization", "Bearer "+f.token)
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Category: %s", err)
	}
	defer resp.Body.Close()

	var rawResults struct {
		Data rawCategoryTotal `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&rawResults)

	var results []CategoryTotal

	r := &rawResults.Data
	var c CategoryTotal
	var spent, earned decimal.Decimal
	id, err := strconv.Atoi(r.ID)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to int: %s", err)
	}
	if len(r.Attributes.Spent) == 1 {
		spent, err = decimal.NewFromString(r.Attributes.Spent[0].Sum)
		if err != nil {
			return nil, fmt.Errorf("could not convert spent sum to decimal: %s", err)
		}
	}
	if len(r.Attributes.Earned) == 1 {
		earned, err = decimal.NewFromString(r.Attributes.Earned[0].Sum)
		if err != nil {
			return nil, fmt.Errorf("could not convert earned sum to decimal: %s", err)
		}
	}

	c.ID = id
	c.Name = r.Attributes.Name
	c.Spent = spent
	c.Earned = earned
	c.Start = start
	c.End = end

	results = append(results, c)

	return results, nil
}
