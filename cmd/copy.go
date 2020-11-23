package cmd

import (
	"bufio"
	"fmt"
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

func init() {
	rootCmd.AddCommand(copyCmd)
}

// Copy implements the copy command
func Copy(cmd *cobra.Command, args []string) {
	git := GetGitlabClient()

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 { // we were piped into
		reader := bufio.NewReader(os.Stdin)
		copy(args, git, viper.GetString("clipboard_name"), reader)
	} else { // invoked without a pipe or redirect
		fmt.Println("ERROR: Please pipe something into STDIN")
		os.Exit(1)
	}
}

func copy(args []string, git gitlab.Client, clipboardName string, reader RuneReader) {

	// read stdin
	// reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	// search snippets for a clipboard with the correct name to update
	snippets, _, err := git.Snippets.ListSnippets(&gitlab.ListSnippetsOptions{})

	BailOnError(err)

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
			Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
		}

		_, _, err = git.Snippets.CreateSnippet(snippetoptions)

		BailOnError(err)

	} else { // update existing snippet
		snippetoptions := &gitlab.UpdateSnippetOptions{
			Title:      gitlab.String(clipboardName),
			FileName:   gitlab.String(clipboardName),
			Content:    gitlab.String(string(output)),
			Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
		}

		_, _, err = git.Snippets.UpdateSnippet(clipboardID, snippetoptions)

		BailOnError(err)
	}

}
