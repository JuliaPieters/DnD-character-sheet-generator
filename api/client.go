package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ---------- Structs ----------

type APIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type APIListResponse struct {
	Results []APIResource `json:"results"`
}

// ---------- Helpers ----------

func getJSON(url string, target interface{}) error {
	url = strings.ToLower(strings.ReplaceAll(url, " ", "-"))

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("fout %d bij ophalen van %s: %s", resp.StatusCode, url, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !json.Valid(body) {
		return fmt.Errorf("ongeldige JSON van %s: %s", url, string(body[:min(100, len(body))]))
	}

	return json.Unmarshal(body, target)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
