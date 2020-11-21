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
	// fmt.Println("copy called")
	// fmt.Println(viper.GetString("host"))
	// fmt.Println(viper.GetString("token"))
	// fmt.Println(viper.ConfigFileUsed())

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 { // we were piped into
		reader := bufio.NewReader(os.Stdin)
		var output []rune

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			output = append(output, input)
		}

		git, err := gitlab.NewClient(viper.GetString("token"), gitlab.WithBaseURL(viper.GetString("gitlab_url")))

		if err != nil {
			log.Fatalf("Failed connect to GitLab: %v", err)
		}

		snippet := &gitlab.CreateSnippetOptions{
			Title:    gitlab.String("pbsnip"),
			FileName: gitlab.String("pbsnip"),
			// Description: gitlab.String("Desc of snippet"),
			Content:    gitlab.String(string(output)),
			Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
		}

		_, _, err = git.Snippets.CreateSnippet(snippet)
		if err != nil {
			log.Fatal(err)
		}

		// for j := 0; j < len(output); j++ {
		// 	fmt.Printf("%c", output[j])
		// }

	} else {
		fmt.Println("ERROR: Please pipe something into STDIN")
		os.Exit(1)
	}
}
