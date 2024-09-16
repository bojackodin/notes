package speller

import (
	"context"
	"encoding/json"
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
}

func (y *YandexSpeller) Check(ctx context.Context, text string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("text", text)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var misspells []Misspell
	if err = json.NewDecoder(resp.Body).Decode(&misspells); err != nil {
		return err
	}

	if len(misspells) > 0 {
		return SpellError{misspells}
	}

	return nil
}
