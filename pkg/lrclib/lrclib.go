package lrclib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"codeberg.org/dankstuff/danklyrics/internal/version"
	"codeberg.org/dankstuff/danklyrics/pkg/errors"
)

// Client is the lrclibgo client that handles all the different API calls to lrclib.net/api
type Client struct {
	Search *SearchService
	Songs  *SongsService
}

// NewClient initializes the genius [Client] to interact with different lrclib.net/api calls.
func NewClient() *Client {
	baseLrcLibUrl := "https://lrclib.net/api/"

	c := &Client{}
	c.Search = &SearchService{
		gClient: newApiClient[[]Song](http.MethodGet, baseLrcLibUrl+"search"),
	}
	c.Songs = &SongsService{
		gClient: newApiClient[Song](http.MethodGet, baseLrcLibUrl+""),
	}

	return c
}

type errorResponse struct {
	Message    string `json:"message"`
	Name       string `json:"name"`
	StatusCode int    `json:"statusCode"`
	Status     int    `json:"status"`
}

type apiClient[R Song | []Song] struct {
	client      *http.Client
	req         *http.Request
	initialPath string
}

func newApiClient[RespType Song | []Song](method, requestPath string) *apiClient[RespType] {
	req, err := http.NewRequest(method, requestPath, http.NoBody)
	if err != nil {
		return nil
	}

	a := &apiClient[RespType]{
		client:      &http.Client{},
		req:         req,
		initialPath: requestPath,
	}
	a.setHeader("User-Agent", fmt.Sprintf("LRCLIB-GO %s (https://codeberg.org/dankstuff/danklyrics/pkg/lrclib)", version.Version))

	return a
}

func (a *apiClient[_]) appendToPath(path string) error {
	a.req.URL.Path += path
	return nil
}

func (a *apiClient[_]) setQueryParam(key, value string) error {
	q := a.req.URL.Query()
	q.Set(key, value)
	a.req.URL.RawQuery = q.Encode()

	return nil
}

func (a *apiClient[_]) setHeader(key, value string) error {
	a.req.Header.Set(key, value)
	return nil
}

func (a *apiClient[R]) callEndpoint() (R, error) {
	defer a.reset()

	var res R
	resp, err := a.client.Do(a.req)
	if err != nil {
		return res, err
	}

	if resp.StatusCode != http.StatusOK {
		return res, &errors.ErrApiError{
			StatusCode: resp.StatusCode,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		var errRes errorResponse
		_ = json.Unmarshal(body, &errRes)
		errStatusCode := errRes.StatusCode
		if errStatusCode == 0 {
			errStatusCode = errRes.Status
		}
		if errStatusCode != http.StatusOK && errStatusCode != 0 {
			return res, &errors.ErrApiError{
				StatusCode: errStatusCode,
				Message:    errRes.Message,
			}
		}

		return res, err
	}

	return res, nil
}

func (a *apiClient[_]) reset() {
	a.req.URL.Path = a.initialPath
}
