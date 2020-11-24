package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Version: "0.0.3",
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

// Execute runs the cli's main root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default $HOME/.glsnip)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".glsnip")
		viper.SetConfigType("yaml")

		viper.SetDefault("clipboard_name", "glsnip")
	}

	viper.SetEnvPrefix("glsnip")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
