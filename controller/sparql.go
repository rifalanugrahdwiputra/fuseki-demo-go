package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const FusekiEndpoint = "http://localhost:3030/univ/sparql"

// SparqlResult untuk parse JSON dari Fuseki
type SparqlResult struct {
	Head struct {
		Vars []string `json:"vars"`
	} `json:"head"`
	Results struct {
		Bindings []map[string]map[string]string `json:"bindings"`
	} `json:"results"`
}

// QuerySPARQL menjalankan query ke Fuseki
func QuerySPARQL(query string) ([]string, [][]string, error) {
	data := url.Values{}
	data.Set("query", query)

	req, err := http.NewRequest("POST", FusekiEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/sparql-results+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var result SparqlResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, fmt.Errorf("Failed to parse JSON: %v\nRaw response: %s", err, string(body))
	}

	// Ambil data tabel
	var rows [][]string
	for _, binding := range result.Results.Bindings {
		row := make([]string, len(result.Head.Vars))
		for i, v := range result.Head.Vars {
			if val, ok := binding[v]; ok {
				row[i] = val["value"]
			} else {
				row[i] = ""
			}
		}
		rows = append(rows, row)
	}

	return result.Head.Vars, rows, nil
}
