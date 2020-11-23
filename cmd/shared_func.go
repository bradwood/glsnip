package cmd

import (
	"log"

	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

// BailOnError logs and exits if there is an error
func BailOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetGitlabClient connects to the Gitlab server or calls BailOnError()
func GetGitlabClient() gitlab.Client {
	git, err := gitlab.NewClient(viper.GetString("token"), gitlab.WithBaseURL(viper.GetString("gitlab_url")))
	BailOnError(err)
	return *git
}
