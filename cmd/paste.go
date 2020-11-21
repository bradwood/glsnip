package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Paste GitLab Snippet to STDOUT",
	Run:   paste,
}

func init() {
	rootCmd.AddCommand(pasteCmd)
}

func paste(cmd *cobra.Command, args []string) {

	git, err := gitlab.NewClient(viper.GetString("token"), gitlab.WithBaseURL(viper.GetString("gitlab_url")))

	if err != nil {
		log.Fatalf("Failed connect to GitLab: %v", err)
	}

	snippets, _, err := git.Snippets.ListSnippets(&gitlab.ListSnippetsOptions{})

	if err != nil {
		log.Fatal(err)
	}

	for _, item := range snippets {

		if item.Title == viper.GetString("clipboard_name") {
			snip, _, err := git.Snippets.SnippetContent(item.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(string(snip))
			break
		}
	}

}
