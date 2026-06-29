package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/internal/config"
)

func makeApiPostRequest[T any](path, token string, body T) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, config.Env().ApiAddress+path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("api responded with a non 200 status")
	}

	return nil
}

func makeApiGetRequest[T any](path, token string) (T, error) {
	req, err := http.NewRequest(http.MethodGet, config.Env().ApiAddress+path, http.NoBody)
	if err != nil {
		return *new(T), err
	}

	req.Header.Set("Authorization", token)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return *new(T), err
	}

	if resp.StatusCode != http.StatusOK {
		return *new(T), errors.New("api responded with a non 200 status")
	}

	var respBody T
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return *new(T), err
	}

	return respBody, nil
}
