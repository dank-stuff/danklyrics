package gonius

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"codeberg.org/dankstuff/danklyrics/internal/version"
	"codeberg.org/dankstuff/danklyrics/pkg/errors"
)

// Client is the Genius client that handles all the different API calls to api.genius.com
type Client struct {
	Account     *AccountService
	Annotations *AnnotationsService
	Artists     *ArtistsService
	Lyrics      *LyricsService
	Search      *SearchService
	Songs       *SongsService
}

func (c *Client) SetPageSize() {
	// per_page=20
}

// NewClient initializes the genius [Client] with the given access token to interact with different api.genius.com calls.
func NewClient(clientId, clientSecret string) *Client {
	baseGeniusUrl := "https://api.genius.com/"

	tokenFetcherInstance := &tokenFetcher{
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	lyricsService := &LyricsService{}

	c := &Client{}
	c.Account = &AccountService{}
	c.Annotations = &AnnotationsService{
		gClient: newApiClient(http.MethodGet, baseGeniusUrl+"annotations/", nil, tokenFetcherInstance),
	}
	c.Artists = &ArtistsService{
		gClient: newApiClient(http.MethodGet, baseGeniusUrl+"artists/", nil, tokenFetcherInstance),
	}
	c.Lyrics = lyricsService
	c.Search = &SearchService{
		gClient: newApiClient(http.MethodGet, baseGeniusUrl+"search/", nil, tokenFetcherInstance),
	}
	c.Songs = &SongsService{
		gClient: newApiClient(http.MethodGet, baseGeniusUrl+"songs/", nil, tokenFetcherInstance),
		lyrics:  lyricsService,
	}

	return c
}

// ApiResponse is the general response structure that is received from different api calls from api.genius.com
type ApiResponse struct {
	Meta *struct {
		Status int `json:"status,omitempty"`
	} `json:"meta,omitempty"`
	Response *struct {
		Annotation *Annotation `json:"annotation,omitempty"`
		Song       *Song       `json:"song,omitempty"`
		Songs      []Song      `json:"songs,omitempty"`
		Artist     *Artist     `json:"artist,omitempty"`
		Hits       []Hit       `json:"hits,omitempty"`
	} `json:"response,omitempty"`
}

type apiClient struct {
	client      *http.Client
	req         *http.Request
	initialPath string
	tokener     *tokenFetcher
}

func newApiClient(method, requestPath string, body io.Reader, tokener *tokenFetcher) *apiClient {
	req, err := http.NewRequest(method, requestPath, body)
	if err != nil {
		return nil
	}

	a := &apiClient{
		client:      &http.Client{},
		req:         req,
		initialPath: requestPath,
		tokener:     tokener,
	}
	a.setHeader("User-Agent", fmt.Sprintf("GONIUS %s (https://codeberg.org/dankstuff/danklyrics)", version.Version))
	a.reset()

	return a
}

func (a *apiClient) setPath(path string) error {
	if _, err := url.Parse(path); err != nil {
		return err
	}
	a.req.URL.Path = path

	return nil
}

func (a *apiClient) appendToPath(path string) error {
	a.req.URL.Path += path
	return nil
}

func (a *apiClient) setQueryParam(key, value string) error {
	q := a.req.URL.Query()
	q.Set(key, value)
	a.req.URL.RawQuery = q.Encode()

	return nil
}

func (a *apiClient) setHeader(key, value string) error {
	a.req.Header.Set(key, value)
	return nil
}

func (a *apiClient) callEndpoint() (ApiResponse, error) {
	// defer a.reset()

	var res ApiResponse
	resp, err := a.client.Do(a.req)
	if err != nil {
		return ApiResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ApiResponse{}, &errors.ErrInvalidToken{
			ProviderName: "Genius API",
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return ApiResponse{}, err
	}

	switch res.Meta.Status {
	case http.StatusOK:
		return res, nil
	case http.StatusNotFound:
		return ApiResponse{}, new(errors.ErrNotFound)
	default:
		return ApiResponse{}, new(errors.ErrApiError)
	}
}

func (a *apiClient) reset() {
	a.req.URL.Path = a.initialPath

	// setting response text_format to plain, so it's readable by the application,
	// rather than using dom or html which need furthor pasring.
	a.setQueryParam("text_format", "plain")

	token, _ := a.tokener.fetch()

	a.setHeader("Authorization", "Bearer "+token)
}

type tokenFetcher struct {
	clientId     string
	clientSecret string
	lastToken    string
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (t *tokenFetcher) fetch() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "manage_annotation")
	data.Set("client_id", t.clientId)
	data.Set("client_secret", t.clientSecret)

	req, err := http.NewRequest(http.MethodPost, "https://api.genius.com/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", new(errors.ErrInvalidToken)
	}

	defer resp.Body.Close()

	var respBody tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}

	t.lastToken = respBody.AccessToken

	return respBody.AccessToken, nil
}
