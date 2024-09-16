package speller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const serviceURL = "http://speller.yandex.net/services/spellservice.json/checkText"

type Speller interface {
	Check(ctx context.Context, text string) error
}

type YandexSpeller struct{}

func NewYandexSpeller() *YandexSpeller {
	return &YandexSpeller{}
}

type Misspell struct {
	Pos  int    `json:"pos"`
	Word string `json:"word"`
	// Row  int    `json:"row"`
	// Col  int    `json:"col"`
	// Suggestions []string `json:"s"`
}

func (y *YandexSpeller) Check(ctx context.Context, text string) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceURL, nil)
	if err != nil {
		return err
	}
	q := r.URL.Query()
	q.Add("text", text)
	r.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	var misspells []Misspell
	if err = json.Unmarshal(body, &misspells); err != nil {
		return err
	}

	if len(misspells) > 0 {
		return ErrorSpell{misspells}
	}

	return nil
}
