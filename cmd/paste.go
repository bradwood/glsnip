package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Paste GitLab Snippet to STDOUT",
	Args:  cobra.NoArgs,
	Run:   Paste,
}

func init() {
	rootCmd.AddCommand(pasteCmd)
}

// Paste implements the paste command
func Paste(cmd *cobra.Command, args []string) {
	git := GetGitlabClient()
	paste(args, git)
}

// TODO: write test for this
func paste(args []string, git gitlab.Client) {

	snippets, _, err := git.Snippets.ListSnippets(&gitlab.ListSnippetsOptions{})

	BailOnError(err)

	for _, item := range snippets {

		if item.Title == viper.GetString("clipboard_name") {
			snip, _, err := git.Snippets.SnippetContent(item.ID)
			BailOnError(err)
			fmt.Print(string(snip))
			break
		}
	}

}
