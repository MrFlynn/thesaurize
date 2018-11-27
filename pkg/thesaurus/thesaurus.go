package thesaurus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// BaseURL base url for BigHugeLabs Thesaurus API.
const BaseURL = "http://words.bighugelabs.com/api"

type (
	// WordClass stores list of lexigraphically classified responses to words.
	WordClass struct {
		Antonym []string `json:"ant"`
		Related []string `json:"rel"`
		Similar []string `json:"sim"`
		Synonym []string `json:"syn"`
		User    []string `json:"usr"`
	}

	// Response stores classified list of lexigraphically similar words classified by overal type.
	Response struct {
		Adjective WordClass `json:"adjective"`
		Adverb    WordClass `json:"adverb"`
		Noun      WordClass `json:"noun"`
		Verb      WordClass `json:"verb"`
	}

	// API contains API key and boolean to shutdown
	API struct {
		Key           string
		usageExceeded bool
		lastCalled    time.Time
	}
)

// Call method requests word synonyms from thesaurus API.
func (a *API) Call(word string) (Response, error) {
	// Construct url.
	url := fmt.Sprintf("%s/2/%s/%s/json", BaseURL, a.Key, word)

	// Call function to handle.
	body, err := doRequest(url)
	if err != nil {
		return Response{}, err
	}

	// Transform JSON into struct for access.
	resp := Response{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response{}, err
	}

	return resp, err
}

func doRequest(url string) ([]byte, error) {
	// net.http Client.
	client := http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	switch c := resp.StatusCode; c {
	case 303:
		// This occurs when a word similar to the one provided is found.
		url, _ := resp.Location()

		newResp, err := client.Get(url.String())
		if err != nil {
			return []byte{}, err
		}

		resp = newResp
	case 404:
		// Return empty content if no word could be found.
		return []byte{}, nil
	case 500:
		// Special cases. Usually this gets triggered if you exceed API limits or get IP banned.
		errorString := fmt.Sprintf("API returned: %s", resp.Status)
		return []byte{}, errors.New(errorString)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
