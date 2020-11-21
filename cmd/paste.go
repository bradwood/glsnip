package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	fmt.Println("paste called")
}
