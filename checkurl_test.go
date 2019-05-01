package phishtank

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		input  string
		indb   bool
		format string
		want   bool
	}{
		{name: "InDatabase", input: "https://example.com", indb: true, format: "application/json", want: true},
		{name: "NotInDatabase", input: "https://example.com", indb: false, format: "application/json", want: false},
		{name: "InvalidContentType", input: "https://example.com", indb: false, format: "application/xml", want: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(createCheckURLHandler(t, c.indb, c.format)))
			defer ts.Close()

			client := New("apikey", OptionAPIURL(ts.URL))
			body, err := client.CheckURL("https://hoge.com")

			if c.format == "application/json" {
				assert.NoError(t, err)

				got := body.Results.InDatabase
				assert.Equal(t, c.want, got)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func createCheckURLHandler(t *testing.T, indb bool, format string) func(http.ResponseWriter, *http.Request) {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", format)
		tr := createCheckURLResponse(t, r.FormValue("url"), indb)
		body, err := json.Marshal(tr)
		if err != nil {
			t.Error(err)
		}
		w.Header().Set(HEADER_REQCOUNTINTERVAL, "300 seconds")
		w.Header().Set(HEADER_REQLIMIT, "10")
		w.Header().Set(HEADER_REQCOUNT, "1")

		if _, err := w.Write([]byte(body)); err != nil {
			t.Error(err)
		}
	}
}

func createCheckURLResponse(t *testing.T, url string, indb bool) CheckURLResponse {
	return CheckURLResponse{
		Meta: createSuccessTestMetadata(t),
		Results: CheckURLResults{
			URL:        url,
			InDatabase: indb,
		},
	}
}
