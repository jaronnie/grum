/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Type string
var Insecure bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grum",
	Short: "git reomote url modify",
	Long:  `git reomote url modify.`,
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	pwd, _ := os.Getwd()
	repo, err := git.PlainOpen(pwd)
	if err != nil {
		return err
	}

	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	for i, v := range cfg.Remotes["origin"].URLs {
		ep, err := transport.NewEndpoint(v)
		if err != nil {
			return err
		}

		var token string
		switch Type {
		case "github":
			token = os.Getenv("GITHUB_TOKEN")
		case "gitlab":
			token = os.Getenv("GITLAB_TOKEN")
		}
		if token == "" {
			return errors.New("empty token")
		}
		protocol := "https"
		if Insecure {
			protocol = "http"
		}
		if ep.Protocol == "http" {
			protocol = "http"
		}

		switch Type {
		case "gitlab":
			if ep.Port != 0 && ep.Port != 22 {
				cfg.Remotes["origin"].URLs[i] = fmt.Sprintf("%s://oauth2:%s@%s:%d/%s", protocol, token, strings.TrimRight(ep.Host, "/"), ep.Port, strings.TrimLeft(ep.Path, "/"))
			} else {
				cfg.Remotes["origin"].URLs[i] = fmt.Sprintf("%s://oauth2:%s@%s/%s", protocol, token, strings.TrimRight(ep.Host, "/"), strings.TrimLeft(ep.Path, "/"))
			}
		default:
			if ep.Port != 0 && ep.Port != 22 {
				cfg.Remotes["origin"].URLs[i] = fmt.Sprintf("%s://%s@%s:%d/%s", protocol, token, strings.TrimRight(ep.Host, "/"), ep.Port, strings.TrimLeft(ep.Path, "/"))
			} else {
				cfg.Remotes["origin"].URLs[i] = fmt.Sprintf("%s://%s@%s/%s", protocol, token, strings.TrimRight(ep.Host, "/"), strings.TrimLeft(ep.Path, "/"))
			}
		}
	}

	err = cfg.Validate()
	if err != nil {
		return err
	}
	err = repo.Storer.SetConfig(cfg)
	if err != nil {
		return err
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.grum.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVarP(&Type, "type", "", "github", "git remote type")
	rootCmd.Flags().BoolVarP(&Insecure, "insecure", "", false, "insecure")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".grum" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".grum")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
