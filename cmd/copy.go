package cmd

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

// RuneReader interface to allow STDIN dependency injection
type RuneReader interface {
	ReadRune() (rune, int, error)
}

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy GitLab Snippet from STDIN",
	Args:  cobra.NoArgs,
	Run:   Copy,
}

var visibility string

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVarP(&visibility, "visibility", "v", "private", "visibility level")
}

// Copy implements the copy command
func Copy(cmd *cobra.Command, args []string) {
	git := GetGitlabClient()

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 { // we were piped into
		reader := bufio.NewReader(os.Stdin)
		copy(args, git, viper.GetString("clipboard_name"), viper.GetInt("project_id"), visibility, reader)
	} else { // invoked without a pipe or redirect
		println("ERROR: Please pipe something into STDIN")
		os.Exit(1)
	}
}

func copy(args []string, git gitlab.Client, clipboardName string, projectID int, visibility string, reader RuneReader) {

	// read stdin
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	// check visibility
	if (visibility != "private") && (visibility != "internal") && (visibility != "public") {
		BailOnError(errors.New("Bad visibility"), "Bad visibility. Must be `private`, `internal or `public`")
	}

	if projectID > 0 { // Project Snippet
		// search snippets for a clipboard with the correct name to update
		snippets, _, err := git.ProjectSnippets.ListSnippets(projectID, &gitlab.ListProjectSnippetsOptions{})
		BailOnError(err, "Could not read Project Snippets")
		var clipboardFound bool = false
		var clipboardID int

		for _, item := range snippets {

			if item.Title == clipboardName {
				clipboardFound = true
				clipboardID = item.ID
				break
			}
		}

		// create a new snippet
		if !clipboardFound {
			snippetoptions := &gitlab.CreateProjectSnippetOptions{
				Title:      gitlab.String(clipboardName),
				FileName:   gitlab.String(clipboardName),
				Content:    gitlab.String(string(output)),
				Visibility: gitlab.Visibility(gitlab.VisibilityValue(visibility)),
			}

			_, _, err = git.ProjectSnippets.CreateSnippet(projectID, snippetoptions)

			BailOnError(err, "Could not create Snippet")

		} else { // update existing snippet
			snippetoptions := &gitlab.UpdateProjectSnippetOptions{
				Title:      gitlab.String(clipboardName),
				FileName:   gitlab.String(clipboardName),
				Content:    gitlab.String(string(output)),
				Visibility: gitlab.Visibility(gitlab.VisibilityValue(visibility)),
			}

			_, _, err = git.ProjectSnippets.UpdateSnippet(projectID, clipboardID, snippetoptions)

			BailOnError(err, "Could not update snippet")
		}

	} else { // Personal Snippet
		// search snippets for a clipboard with the correct name to update
		snippets, _, err := git.Snippets.ListSnippets(&gitlab.ListSnippetsOptions{})
		BailOnError(err, "Could not read Personal Snippets")
		var clipboardFound bool = false
		var clipboardID int

		for _, item := range snippets {

			if item.Title == clipboardName {
				clipboardFound = true
				clipboardID = item.ID
				break
			}
		}

		// create a new snippet
		if !clipboardFound {
			snippetoptions := &gitlab.CreateSnippetOptions{
				Title:      gitlab.String(clipboardName),
				FileName:   gitlab.String(clipboardName),
				Content:    gitlab.String(string(output)),
				Visibility: gitlab.Visibility(gitlab.VisibilityValue(visibility)),
			}

			_, _, err = git.Snippets.CreateSnippet(snippetoptions)

			BailOnError(err, "Could not create Snippet")

		} else { // update existing snippet
			snippetoptions := &gitlab.UpdateSnippetOptions{
				Title:      gitlab.String(clipboardName),
				FileName:   gitlab.String(clipboardName),
				Content:    gitlab.String(string(output)),
				Visibility: gitlab.Visibility(gitlab.VisibilityValue(visibility)),
			}

			_, _, err = git.Snippets.UpdateSnippet(clipboardID, snippetoptions)

			BailOnError(err, "Could not update snippet")
		}

	}

}
