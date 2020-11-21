package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.0.1",
	Use:     "glsnip",
	Short:   "Copy and Paste using GitLab Snippets",
	Long: `This app behaves like pbcopy(1) and pbpaste(1) on a Mac, or like xclip(1) on
Linux, except, instead of using a local clipboard, it uses GitLab Snippets.

Configuration:
  Create a YAML-formatted config file at $HOME/.glsnip like this:

    gitlab_url: https://url.of.gitlab.server/
    token: USERTOKEN
    clipboard_name: glsnip

Environment variables:
  Instead of using a configuration file, you may set environment variables by
  prefixing the key in the configuration file with GLSNIP_ and then converting
  all alphabetic characters to UPPERCASE.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default $HOME/.glsnip)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".glsnip" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".glsnip")
		viper.SetConfigType("yaml")

		viper.SetDefault("clipboard_name", "glsnip")
	}

	viper.SetEnvPrefix("glsnip")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()

	// if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}
