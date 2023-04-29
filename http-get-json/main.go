package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w *Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o *Occurrence) GetResponse() string {
	out := []string{}
	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occurrence))
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

type Response interface {
	GetResponse() string
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}
	resp, err := doRequest(args[1])
	if err != nil {
		log.Fatal(err)
	}

	if resp == nil {
		fmt.Printf("No Response\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", resp.GetResponse())
}

func doRequest(requestURL string) (Response, error) {
	if _, err := url.ParseRequestURI(requestURL); err != nil {
		return nil, fmt.Errorf("Usage: ./http-get <url>\n\nURL is not valid URL: %s\n", requestURL)
	}

	response, err := http.Get(requestURL)

	if err != nil {
		return nil, fmt.Errorf("Resposne Error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("Read all error: %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output (HTTP Code %d): %s\n", response.StatusCode, string(body))
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %s", err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal error: %s", err)
		}
		return &words, nil

	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal error: %s", err)
		}
		return &occurrence, nil
	}

	return nil, nil
}
