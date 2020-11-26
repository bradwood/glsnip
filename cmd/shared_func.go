package cmd

import (
	"os"

	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

// BailOnError exits if there is an error
func BailOnError(err error, msg string) {
	if err != nil {
		// println(msg, err)
		println("Error:", msg)
		os.Exit(1)
	}
}

// GetGitlabClient connects to the Gitlab server or calls BailOnError()
func GetGitlabClient() gitlab.Client {
	git, err := gitlab.NewClient(viper.GetString("token"), gitlab.WithBaseURL(viper.GetString("gitlab_url")))
	BailOnError(err, "Could not connect to server")
	return *git
}
