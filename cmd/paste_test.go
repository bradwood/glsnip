package cmd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaste(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v4/snippets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		mustWriteHTTPResponse(t, w, "../testdata/list_snippets_paste.json")
	})

	mux.HandleFunc("/api/v4/snippets/41/raw", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "PASTED_DATA")
	})

	var emptyStringArray []string

	output := paste(emptyStringArray, *client, "glsnip")

	assert.Equal(t, "PASTED_DATA", output)
}
