package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestCopyUpdate(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v4/snippets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		mustWriteHTTPResponse(t, w, "../testdata/list_snippets_paste.json")
	})

	mux.HandleFunc("/api/v4/snippets/41", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testURL(t, r, "/api/v4/snippets/41") // this line verifies the update PUT occurs
		fmt.Fprint(w,
			`{
			"id": 1,
			"title": "test",
			"description": "description of snippet",
			"visibility": "internal",
			"author": {
				"id": 1,
				"username": "john_smith",
				"email": "john@example.com",
				"name": "John Smith",
				"state": "active",
				"created_at": "2012-05-23T08:00:58Z"
			},
			"expires_at": null,
			"updated_at": "2012-06-28T10:52:04Z",
			"created_at": "2012-06-28T10:52:04Z",
			"project_id": null,
			"web_url": "http://example.com/snippets/1",
			"raw_url": "http://example.com/snippets/1/raw",
			"ssh_url_to_repo": "ssh://git@gitlab.example.com:snippets/1.git",
			"http_url_to_repo": "https://gitlab.example.com/snippets/1.git",
			"file_name": "renamed.md",
			"files": [
				{
				"path": "renamed.md",
				"raw_url": "https://gitlab.example.com/-/snippets/1/raw/master/renamed.md"
				}
			]
			}`)
	})

	var emptyStringArray []string

	copy(emptyStringArray, *client, "glsnip", -1, "internal", bytes.NewBufferString("std in"))

}

func TestCopyCreate(t *testing.T) {
	mux, server, client := setup(t)
	defer teardown(server)

	mux.HandleFunc("/api/v4/snippets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			testMethod(t, r, "GET")
			mustWriteHTTPResponse(t, w, "../testdata/list_snippets_paste_no_glsnip.json")
		} else if r.Method == "POST" {
			fmt.Fprint(w,
				`{
				"id": 1,
				"title": "test",
				"description": "description of snippet",
				"visibility": "internal",
				"author": {
					"id": 1,
					"username": "john_smith",
					"email": "john@example.com",
					"name": "John Smith",
					"state": "active",
					"created_at": "2012-05-23T08:00:58Z"
				},
				"expires_at": null,
				"updated_at": "2012-06-28T10:52:04Z",
				"created_at": "2012-06-28T10:52:04Z",
				"project_id": null,
				"web_url": "http://example.com/snippets/1",
				"raw_url": "http://example.com/snippets/1/raw",
				"ssh_url_to_repo": "ssh://git@gitlab.example.com:snippets/1.git",
				"http_url_to_repo": "https://gitlab.example.com/snippets/1.git",
				"file_name": "renamed.md",
				"files": [
					{
					"path": "renamed.md",
					"raw_url": "https://gitlab.example.com/-/snippets/1/raw/master/renamed.md"
					}
				]
				}`)
		}
	})

	var emptyStringArray []string

	copy(emptyStringArray, *client, "glsnip", -1, "private", bytes.NewBufferString("std in"))

}
