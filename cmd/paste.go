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
	output := paste(args, git, viper.GetString("clipboard_name"), viper.GetInt("project_id"))
	fmt.Print(output)
}

func paste(args []string, git gitlab.Client, clipboardName string, projectID int) string {

	var output string

	if projectID > 0 { // Project Snippet
		snippets, _, err := git.ProjectSnippets.ListSnippets(projectID, &gitlab.ListProjectSnippetsOptions{})
		BailOnError(err, "Could not read Project Snippets")
		for _, item := range snippets {

			if item.Title == clipboardName {
				snip, _, err := git.Snippets.SnippetContent(item.ID)
				BailOnError(err, "Could not read Snippet contents")
				output = string(snip)
				break
			}
		}
		return output

	} else { // Personal Snippet
		snippets, _, err := git.Snippets.ListSnippets(&gitlab.ListSnippetsOptions{})
		BailOnError(err, "Could not read Personal Snippets")
		for _, item := range snippets {

			if item.Title == clipboardName {
				snip, _, err := git.Snippets.SnippetContent(item.ID)
				BailOnError(err, "Could not read Snippet contents")
				output = string(snip)
				break
			}
		}
		return output

	}
}
