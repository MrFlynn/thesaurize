package loader

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

type word struct {
	ID       string         `json:"id"`
	Match    *regexp.Regexp `json:"match"`
	Tags     []string       `json:"tags"`
	Severity int            `json:"severity"`
}

func (w *word) UnmarshalJSON(data []byte) error {
	type alias word
	aux := &struct {
		Match string `json:"match"`
		*alias
	}{
		alias: (*alias)(w),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	w.Match, err = regexp.Compile(aux.Match)
	return err
}

type profanityFilter struct {
	index      []word
	categories map[string]struct{}
}

func (f *profanityFilter) init(url string) error {
	response, err := http.Get(url)
	if err != nil || response.StatusCode == http.StatusNotFound {
		return errors.New("unable to get filter index")
	}

	return json.NewDecoder(response.Body).Decode(&f.index)
}

func (f *profanityFilter) match(text string) bool {
	for _, word := range f.index {
		for _, tag := range word.Tags {
			if _, ok := f.categories[tag]; ok && word.Match.MatchString(text) {
				return true
			}
		}
	}

	return false
}
