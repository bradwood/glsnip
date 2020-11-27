package cmd

import (
	"errors"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, profile string
var cfgFileFound bool = true

var rootCmd = &cobra.Command{
	Version: "0.2.0",
	Use:     "glsnip",
	Short:   "Copy and paste using GitLab Snippets",
	Long: `This app behaves like pbcopy(1) and pbpaste(1) on a Mac, or like xclip(1) on
Linux, except, instead of using a local clipboard, it uses GitLab Snippets.

Configuration:
  Create a YAML-formatted config file (default location $HOME/.glsnip). You must
  include at least a single server profile YAML block called 'default', like
  this:

    ---
    default:
      gitlab_url: https://url.of.gitlab.server/
      token: USERTOKEN
      clipboard_name: glsnip
    ...

  Multiple additional server profile blocks can be added using any block name,
  like this:

    ...
    work:
      gitlab_url: https://url.of.work.server/
      token: USERTOKENWORK
      clipboard_name: glsnip
    ...

Environment variables:
  Instead of using a configuration file, you may set environment variables by
  prefixing the key in a configuration file block with GLSNIP_ and then
  converting all alphabetic characters to UPPERCASE. Note that environment
  variables will override any configuration specified in the configuration file,
  regardless of the profile specified. You may specify a server profile by
  setting GLSNIP_PROFILE.`}

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
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "server profile")
}

func initConfig() {

	// determine config file path
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		BailOnError(err, "Could not find $HOME")

		viper.AddConfigPath(home)
		viper.SetConfigName(".glsnip")
		viper.SetConfigType("yaml")
	}

	// bind profile flag and env var to viper setting
	viper.SetEnvPrefix("glsnip")
	viper.BindEnv("profile")
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))

	// set defaults on top-level viper settings
	viper.SetDefault("clipboard_name", "glsnip")

	// try read in config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, ignoring as env vars may have been passed
			cfgFileFound = false
		} else {
			BailOnError(err, "Could not find $HOME")
			// Config file was found but another error was produced
		}
	}

	if cfgFileFound {
		if viper.GetString("profile") != "default" && !viper.IsSet(viper.GetString("profile")) {
			BailOnError(errors.New("Bad profile"), "Bad profile")
		}

		for _, key := range viper.AllKeys() {
			if keySplit := strings.Split(key, "."); keySplit[0] == viper.GetString("profile") {
				viper.Set(keySplit[1], viper.GetString(key))
			}
		}
	}

	viper.AutomaticEnv()

	if !viper.IsSet("gitlab_url") {
		BailOnError(errors.New("Bad or missing GitLab server URL"), "Bad or missing GitLab server URL")
	}
	if !viper.IsSet("token") {
		BailOnError(errors.New("Bad or missing GitLab server token"), "Bad or missing GitLab server token")
	}
}
