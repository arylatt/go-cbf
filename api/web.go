package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	CambridgeBeerFestivalHost = "https://data.cambridgebeerfestival.com"
	DefaultUserAgent          = "go-cbf (https://github.com/arylatt/go-cbf)"

	CategoryAppleJuice        = "apple-juice"
	CategoryBeer              = "beer"
	CategoryCider             = "cider"
	CategoryInternationalBeer = "international-beer"
	CategoryMead              = "mead"
	CategoryPerry             = "perry"
	CategoryWine              = "wine"
)

var (
	ErrEventOrCategoryNotFound = errors.New("the event or category does not appear to exist")

	KnownCategories = []string{
		CategoryAppleJuice,
		CategoryBeer,
		CategoryCider,
		CategoryInternationalBeer,
		CategoryMead,
		CategoryPerry,
		CategoryWine,
	}
)

type WebClient struct {
	client    *http.Client
	host      *url.URL
	userAgent string
}

func NewWebClient(opts ...WebClientOption) (wc *WebClient, err error) {
	hostURL, _ := url.Parse(CambridgeBeerFestivalHost)

	wc = &WebClient{
		client:    &http.Client{},
		host:      hostURL,
		userAgent: DefaultUserAgent,
	}

	for _, opt := range opts {
		if err = opt(wc); err != nil {
			return
		}
	}

	return
}

func (wc *WebClient) Get(event, category string) (r *Response, err error) {
	reqURL, err := wc.host.Parse(fmt.Sprintf("/%s/%s.json", event, category))
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return
	}

	req.Header.Set("user-agent", wc.userAgent)

	resp, err := wc.client.Do(req)
	if err != nil {
		return r, &WebClientHTTPError{
			Err:      err,
			Response: resp,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return r, &WebClientHTTPError{
			Err:      ErrEventOrCategoryNotFound,
			Response: resp,
		}
	}

	defer resp.Body.Close()

	r = &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)

	return
}

type WebClientHTTPError struct {
	Err      error
	Response *http.Response
}

func (err *WebClientHTTPError) Error() string {
	return err.Err.Error()
}

type WebClientOption func(wc *WebClient) error

func WithHost(host string) WebClientOption {
	return func(wc *WebClient) (err error) {
		wc.host, err = url.Parse(host)
		return
	}
}

func WithClient(client *http.Client) WebClientOption {
	return func(wc *WebClient) (err error) {
		wc.client = client
		return
	}
}

func WithUserAgent(userAgent string) WebClientOption {
	return func(wc *WebClient) (err error) {
		if userAgent != "" {
			wc.userAgent = userAgent
		}
		return
	}
}
