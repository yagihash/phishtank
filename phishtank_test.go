package phishtank

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	testAPIKey := func(t *testing.T) {
		t.Helper()
		cases := []struct {
			name  string
			input string
			want  string
		}{
			{name: "CheckAPIKey", input: strings.Repeat("a", 64), want: strings.Repeat("a", 64)},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				client := New(c.input)
				got := client.apikey
				if diff := cmp.Diff(got, c.want); diff != "" {
					t.Errorf("(-got +want)%s", diff)
				}
			})
		}
	}

	testAPIURL := func(t *testing.T) {
		t.Helper()
		cases := []struct {
			name  string
			input Option
			want  string
		}{
			{name: "CheckDefaultAPIURL", input: nil, want: APIURL},
			{name: "CheckOptionAPIURL", input: OptionAPIURL("https://example.com/foo"), want: "https://example.com/foo"},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				var client *Client
				if c.input != nil {
					client = New("apikey", c.input)
				} else {
					client = New("apikey")
				}
				got := client.endpoint
				if diff := cmp.Diff(got, c.want); diff != "" {
					t.Errorf("(-got +want)%s", diff)
				}
			})
		}
	}

	testHttpClient := func(t *testing.T) {
		t.Helper()
		cases := []struct {
			name  string
			input Option
			want  *http.Client
		}{
			{name: "CheckDefaultClient", input: nil, want: &http.Client{}},
			{name: "CheckOptionClient", input: OptionHttpClient(&http.Client{Timeout: 100}), want: &http.Client{Timeout: 100}},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				var client *Client
				if c.input != nil {
					client = New("apikey", c.input)
				} else {
					client = New("apikey")
				}
				got := client.httpclient
				if diff := cmp.Diff(got, c.want); diff != "" {
					t.Errorf("(-got +want)%s", diff)
				}
			})
		}
	}

	testAPIKey(t)
	testAPIURL(t)
	testHttpClient(t)
}
