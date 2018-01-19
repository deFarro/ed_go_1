package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Type for config
type devConfig struct {
	Endpoint     string
	QueryMapping map[string]string
	APIKey       int
}

// Method to map query keys
func (config devConfig) mapQueryKeys(src url.Values) url.Values {
	queries := url.Values{}
	for new, old := range config.QueryMapping {
		queries.Add(new, src.Get(old))
	}
	return queries
}

// Method to handle routes
func (config devConfig) handle(w http.ResponseWriter, r *http.Request) {
	// Ignore favicon requests
	if strings.Contains(r.URL.Path, "favicon") {
		return
	}
	r.ParseForm()
	fmt.Printf("Got a request with this query: %s\n", r.Form.Encode())

	newQuery := config.mapQueryKeys(r.Form)
	redirectURL := config.Endpoint + "?" + newQuery.Encode()

	response, err := http.Get(redirectURL)
	fmt.Printf("GET request is send with URL: %s\n", redirectURL)

	if err != nil {
		fmt.Fprintf(w, "Error occured while requesting data: %s", err)
		return
	}
	fmt.Printf("Succesfully got a response with status code: %s\n", response.Status)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Fprintf(w, "Error occured while parsing response: %s", err)
		return
	}
	fmt.Println("Response was succesfully parsed")

	fmt.Fprintf(w, string(body))
}

func main() {
	// Importing config from a file
	source, err := ioutil.ReadFile("devConfig.json")
	if err != nil {
		fmt.Printf("File reading error: %s\n", err)
		os.Exit(1)
	}

	// Implementing config
	config := devConfig{}
	json.Unmarshal(source, &config)
	router := devConfig{Endpoint: config.Endpoint, QueryMapping: config.QueryMapping, APIKey: config.APIKey}

	// Start server
	http.HandleFunc("/", router.handle)
	http.ListenAndServe(":3000", nil)
}
