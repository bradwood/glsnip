package cmd

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/xanzy/go-gitlab"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// setup sets up a test HTTP server along with a gitlab.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(t *testing.T) (*http.ServeMux, *httptest.Server, *gitlab.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the Gitlab client being tested.
	client, err := gitlab.NewClient("", gitlab.WithBaseURL(server.URL))
	if err != nil {
		server.Close()
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, server, client
}

// teardown closes the test HTTP server.
func teardown(server *httptest.Server) {
	server.Close()
}

func testURL(t *testing.T, r *http.Request, want string) {
	if got := r.RequestURI; got != want {
		t.Errorf("Request url: %+v, want %s", got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)

	if err != nil {
		t.Fatalf("Failed to Read Body: %v", err)
	}

	if got := buffer.String(); got != want {
		t.Errorf("Request body: %s, want %s", got, want)
	}
}

func mustWriteHTTPResponse(t *testing.T, w io.Writer, fixturePath string) {
	f, err := os.Open(fixturePath)
	if err != nil {
		t.Fatalf("error opening fixture file: %v", err)
	}

	if _, err = io.Copy(w, f); err != nil {
		t.Fatalf("error writing response: %v", err)
	}
}

//TODO: is this needed
func errorOption(*retryablehttp.Request) error {
	return errors.New("RequestOptionFunc returns an error")
}

//TODO: is this needed
func loadFixture(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return content
}

func patchEnv(key, value string) func() {
	bck := os.Getenv(key)
	deferFunc := func() {
		os.Setenv(key, bck)
	}

	if value != "" {
		os.Setenv(key, value)
	} else {
		os.Unsetenv(key)
	}

	return deferFunc
}

func TestInitConfig(t *testing.T) {

	// create HOME dir and associated env var
	homeDir, _ := ioutil.TempDir("", "glsnip*")
	// set up cleanup function for HOME dir
	defer func() { os.RemoveAll(homeDir) }()

	// set up homedir
	defer patchEnv("HOME", homeDir)()

	// don't cache the homedir for the duration of this test
	homedir.DisableCache = true
	defer func() { homedir.DisableCache = false }()

	// inject  config file
	confFile, _ := os.Create(homeDir + "/.glsnip")
	confFile.WriteString(`
default:
  gitlab_url: https://some.url/
  token: USERTOKEN
`)
	confFile.Close()
	// inject flags
	// inject env vars

	println(os.Getenv("HOME"))
	initConfig()

	assert.Equal(t, viper.GetString("gitlab_url"), "https://some.url/")
	// test values for viper

}
