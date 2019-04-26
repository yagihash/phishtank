package phishtank

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.Equal(t, c.want, got)
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
				assert.Equal(t, c.want, got)
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
				assert.Equal(t, c.want, got)
			})
		}
	}

	testAPIKey(t)
	testAPIURL(t)
	testHttpClient(t)
}

func TestPost(t *testing.T) {
	t.Parallel()

	t.Run("CheckSuccessResponse", func(t *testing.T) {
		datatext := "data"

		ts := httptest.NewServer(http.HandlerFunc(createSuccessHandler(t, datatext)))
		defer ts.Close()

		client := New("apikey", OptionAPIURL(ts.URL))
		body, err := client.post()
		assert.NoError(t, err)

		got := &testResponse{}
		err = json.Unmarshal(body, &got)
		assert.NoError(t, err)

		want := createSuccessTestResponse(t, datatext)

		assert.Equal(t, want, got)
	})

	t.Run("CheckErrorResponse", func(t *testing.T) {
		errortext := "error"

		ts := httptest.NewServer(http.HandlerFunc(createErrorHandler(t, errortext)))
		defer ts.Close()

		client := New("apikey", OptionAPIURL(ts.URL))
		body, err := client.post()
		assert.NoError(t, err)

		got := &testResponse{}
		err = json.Unmarshal(body, &got)
		assert.NoError(t, err)

		want := createErrorTestResponse(t, errortext)

		assert.Equal(t, want, got)
	})

	t.Run("CheckXMLResponse", func(t *testing.T) {
		datatext := "data"

		ts := httptest.NewServer(http.HandlerFunc(createXMLHandler(t, datatext)))
		defer ts.Close()

		client := New("apikey", OptionAPIURL(ts.URL))
		_, err := client.post()
		assert.Error(t, err)
	})

	t.Run("CheckEchoResponse", func(t *testing.T) {
		datatext := "foobar"

		ts := httptest.NewServer(http.HandlerFunc(createEchoHandler(t)))
		defer ts.Close()

		client := New("apikey", OptionAPIURL(ts.URL))
		param := &Param{
			name:  "msg",
			value: datatext,
		}

		body, err := client.post(param)
		assert.NoError(t, err)

		got := &testResponse{}
		err = json.Unmarshal(body, &got)
		assert.NoError(t, err)

		want := createSuccessTestResponse(t, datatext)

		assert.Equal(t, want.Data.DataText, got.Data.DataText)
	})
}

func TestInvalidContentType(t *testing.T) {
	err := &InvalidContentType{}
	got := err.Error()
	want := "Response is not JSON"
	assert.Equal(t, want, got)
}

func createEchoHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		tr := createSuccessTestResponse(t, r.FormValue("msg"))
		body, err := json.Marshal(tr)
		if err != nil {
			t.Error(err)
		}
		w.Write([]byte(body))
	}
}

func createTestHandler(t *testing.T, r *testResponse, format string) func(http.ResponseWriter, *http.Request) {
	t.Helper()

	body, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", format)
		w.Write([]byte(body))
	}
}

func createSuccessHandler(t *testing.T, datatext string) func(http.ResponseWriter, *http.Request) {
	t.Helper()
	tr := createSuccessTestResponse(t, datatext)
	return createTestHandler(t, tr, "application/json")
}

func createErrorHandler(t *testing.T, errortext string) func(http.ResponseWriter, *http.Request) {
	t.Helper()
	tr := createErrorTestResponse(t, errortext)
	return createTestHandler(t, tr, "application/json")
}

func createXMLHandler(t *testing.T, datatext string) func(http.ResponseWriter, *http.Request) {
	t.Helper()
	tr := createSuccessTestResponse(t, datatext)
	return createTestHandler(t, tr, "application/xml")
}

func createSuccessTestResponse(t *testing.T, datatext string) *testResponse {
	t.Helper()
	return &testResponse{
		Metadata: createSuccessTestMetadata(t),
		Data:     createTestResponse(t, datatext),
	}
}

func createErrorTestResponse(t *testing.T, errortext string) *testResponse {
	t.Helper()
	return &testResponse{
		Metadata:  createErrorTestMetadata(t),
		ErrorText: errortext,
	}
}

func createSuccessTestMetadata(t *testing.T) ResponseMetadata {
	t.Helper()
	return ResponseMetadata{
		Timestamp: "2019-04-26T04:46:25+00:00",
		ServerID:  "deadbeef",
		Status:    "success",
		RequestID: "111.22.33.44.aaaaaaaaaaaaaa.11111111",
	}
}

func createErrorTestMetadata(t *testing.T) ResponseMetadata {
	t.Helper()
	return ResponseMetadata{
		Timestamp: "2019-04-26T04:46:25+00:00",
		ServerID:  "deadbeef",
		Status:    "error",
		RequestID: "111.22.33.44.aaaaaaaaaaaaaa.11111111",
	}
}

func createTestResponse(t *testing.T, datatext string) testResponseData {
	t.Helper()
	return testResponseData{
		DataText: datatext,
	}
}

type testResponse struct {
	Metadata  ResponseMetadata `json:"meta"`
	Data      testResponseData `json:"data"`
	ErrorText string           `json:"errortext"`
}

type testResponseData struct {
	DataText string `json:"datatext"`
}
