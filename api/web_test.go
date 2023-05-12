package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const nonExistantEvent = "this-event-does-not-exist.json"

var hostURL *url.URL

type expectedEvent struct {
	event string
}

func testServer(t *testing.T, expectedEvent *expectedEvent) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
		event, category := path[0], path[1]

		assert.Equal(t, expectedEvent.event, event)

		fileName := fmt.Sprintf("testdata/%s", category)

		_, statErr := os.Stat(fileName)
		if category == nonExistantEvent && statErr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Test Error - File should not exist!"))
			return
		} else if statErr != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Event/Category does not exist"))
			return
		}

		file, err := os.ReadFile(fileName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Test Error - %s", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(file)
	}))
}

func TestConstCambridgeBeerFestivalHost(t *testing.T) {
	var err error

	hostURL, err = url.Parse(CambridgeBeerFestivalHost)

	assert.NoError(t, err)
}

func TestNewWebClient(t *testing.T) {
	TestConstCambridgeBeerFestivalHost(t)

	localhostStr := "http://localhost"
	localhost, _ := url.Parse(localhostStr)
	client := &http.Client{}

	tests := []struct {
		opts   []WebClientOption
		client *WebClient
	}{
		{
			client: &WebClient{client: &http.Client{}, host: hostURL},
		},
		{
			opts:   []WebClientOption{WithClient(client)},
			client: &WebClient{client: client, host: hostURL},
		},
		{
			opts:   []WebClientOption{WithHost(localhostStr)},
			client: &WebClient{client: &http.Client{}, host: localhost},
		},
	}

	for _, test := range tests {
		wc, err := NewWebClient(test.opts...)

		assert.NoError(t, err)
		assert.Equal(t, test.client, wc)
	}
}

func TestWebClientGet(t *testing.T) {
	expectedEvent := &expectedEvent{}
	httpErr := &WebClientHTTPError{}

	srv := testServer(t, expectedEvent)

	defer srv.Close()

	tests := []struct {
		event    string
		category string
		assert   func(*Response, error)
	}{
		{
			event:    "cbf2023",
			category: "beer",
			assert: func(r *Response, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
			},
		},
		{
			event:    "cbf2023",
			category: nonExistantEvent,
			assert: func(r *Response, err error) {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &httpErr)
				assert.Nil(t, r)
			},
		},
	}

	srvURL, _ := url.Parse(srv.URL)

	wc := WebClient{
		client: srv.Client(),
		host:   srvURL,
	}

	for _, test := range tests {
		expectedEvent.event = test.event
		test.assert(wc.Get(test.event, test.category))
	}
}
