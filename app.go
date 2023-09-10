package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type RedirectionConfig struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func main() {

	// Read the port from an environment variable, or use a default value (e.g., 8080).
	port := os.Getenv("PAGABEIBIS_DOT_LINK_PORT")
	if port == "" {
		port = "8080"
	}

	// Path
	redirections_config := os.Getenv("PAGABEIBIS_DOT_LINK_REDIRECTION_CONFIG")
	if redirections_config == "" {
		redirections_config = "redirections.jsonl"
	}
	addr := ":" + port

	file, err := os.Open(redirections_config)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var entries []RedirectionConfig

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// Parse the JSON from each line and add it to the list
		var entry RedirectionConfig
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create an HTTP server that handles redirections.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Iterate over the redirection configurations.
		for _, entry := range entries {
			if r.URL.Path == entry.Path {
				fmt.Printf("path: %s, url: %s\n", entry.Path, entry.URL)
				http.Redirect(w, r, entry.URL, http.StatusSeeOther)
				return
			}
		}

		// If no matching path is found, handle as needed (e.g., return a 404 error).
		http.NotFound(w, r)
	})

	fmt.Printf("Server listening on %s...\n", addr)
	http.ListenAndServe(addr, nil)
}
