/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jaronnie/grum/pkg/githosting"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/cobra"
)

var (
	Type     string
	Insecure bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grum",
	Short: "git remote url modify",
	Long:  `git remote url modify.`,
	Run:   run,
}

func run(*cobra.Command, []string) {
	pwd, _ := os.Getwd()
	repo, err := git.PlainOpen(pwd)
	cobra.CheckErr(err)

	cfg, err := repo.Config()
	cobra.CheckErr(err)

	for i, v := range cfg.Remotes["origin"].URLs {
		cfg.Remotes["origin"].URLs[i], err = getNewRemoteUrl(v)
		cobra.CheckErr(err)

		user, err := User(v, Token(Type))
		cobra.CheckErr(err)
		cfg.User.Name = user.Username
		cfg.User.Email = user.Email
	}

	err = cfg.Validate()
	cobra.CheckErr(err)
	err = repo.Storer.SetConfig(cfg)
	cobra.CheckErr(err)
}

func getNewRemoteUrl(v string) (string, error) {
	token := Token(Type)
	if token == "" {
		return "", errors.New("empty token")
	}

	ep, err := Endpoint(v)
	cobra.CheckErr(err)

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
			return fmt.Sprintf("%s://oauth2:%s@%s:%d/%s", protocol, token, strings.TrimRight(ep.Host, "/"), ep.Port, strings.TrimLeft(ep.Path, "/")), nil
		} else {
			return fmt.Sprintf("%s://oauth2:%s@%s/%s", protocol, token, strings.TrimRight(ep.Host, "/"), strings.TrimLeft(ep.Path, "/")), nil
		}
	default:
		if ep.Port != 0 && ep.Port != 22 {
			return fmt.Sprintf("%s://%s@%s:%d/%s", protocol, token, strings.TrimRight(ep.Host, "/"), ep.Port, strings.TrimLeft(ep.Path, "/")), nil
		} else {
			return fmt.Sprintf("%s://%s@%s/%s", protocol, token, strings.TrimRight(ep.Host, "/"), strings.TrimLeft(ep.Path, "/")), nil
		}
	}
}

func Endpoint(url string) (*transport.Endpoint, error) {
	return transport.NewEndpoint(url)
}

func Token(t string) string {
	var token string
	switch t {
	case "github":
		token = os.Getenv("GITHUB_TOKEN")
	case "gitlab":
		token = os.Getenv("GITLAB_TOKEN")
	case "gitea":
		token = os.Getenv("GITEA_TOKEN")
	}
	return token
}

func User(url string, token string) (*githosting.UserInfo, error) {
	var apiUrl string
	ep, err := transport.NewEndpoint(url)
	cobra.CheckErr(err)

	if Type == githosting.GITHUB {
		apiUrl = "https://api.github.com"
	} else {
		apiUrl = fmt.Sprintf("%s://%s", ep.Protocol, ep.Host)
	}

	gh, err := githosting.New(githosting.Config{
		Type:  Type,
		Url:   apiUrl,
		Token: token,
	})
	cobra.CheckErr(err)

	return gh.GetUserInfo()
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
	rootCmd.PersistentFlags().StringVarP(&Type, "type", "", "github", "git remote type")
	rootCmd.PersistentFlags().BoolVarP(&Insecure, "insecure", "", false, "insecure")
}
