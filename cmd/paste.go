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
	output := paste(args, git, viper.GetString("clipboard_name"))
	fmt.Print(output)
}

func paste(args []string, git gitlab.Client, clipboardName string) string {

	var output string

	snippets, _, err := git.Snippets.ListSnippets(&gitlab.ListSnippetsOptions{})

	BailOnError(err)

	for _, item := range snippets {

		if item.Title == clipboardName {
			snip, _, err := git.Snippets.SnippetContent(item.ID)
			BailOnError(err)
			output = string(snip)
			break
		}
	}

	return output
}
