package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy GitLab Snippet from STDIN",
	Run:   copy,
}

func init() {
	rootCmd.AddCommand(copyCmd)
}

// TODO: do we even need args here?
func copy(cmd *cobra.Command, args []string) {

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 { // we were piped into

		// read stdin
		reader := bufio.NewReader(os.Stdin)
		var output []rune

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			output = append(output, input)
		}

		// connect to gitlab
		git, err := gitlab.NewClient(viper.GetString("token"), gitlab.WithBaseURL(viper.GetString("gitlab_url")))

		if err != nil {
			log.Fatalf("Failed connect to GitLab: %v", err)
		}

		// search snippets for a clipboard with the correct name to update
		snippets, _, err := git.Snippets.ListSnippets(&gitlab.ListSnippetsOptions{})

		if err != nil {
			log.Fatal(err)
		}

		var clipboard_found bool = false
		var clipboard_id int

		for _, item := range snippets {

			if item.Title == viper.GetString("clipboard_name") {
				clipboard_found = true
				clipboard_id = item.ID
				break
			}
		}

		// create a new snippet
		if !clipboard_found {
			snippetoptions := &gitlab.CreateSnippetOptions{
				Title:      gitlab.String(viper.GetString("clipboard_name")),
				FileName:   gitlab.String(viper.GetString("clipboard_name")),
				Content:    gitlab.String(string(output)),
				Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
			}
			_, _, err = git.Snippets.CreateSnippet(snippetoptions)
			if err != nil {
				log.Fatal(err)
			}
		} else { // update existing snippet
			snippetoptions := &gitlab.UpdateSnippetOptions{
				Title:      gitlab.String(viper.GetString("clipboard_name")),
				FileName:   gitlab.String(viper.GetString("clipboard_name")),
				Content:    gitlab.String(string(output)),
				Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
			}
			_, _, err = git.Snippets.UpdateSnippet(clipboard_id, snippetoptions)
			if err != nil {
				log.Fatal(err)
			}
		}

	} else {
		fmt.Println("ERROR: Please pipe something into STDIN")
		os.Exit(1)
	}
}
