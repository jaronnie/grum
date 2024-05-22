/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaronnie/grum/pkg/execx"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "grum clone means git clone but do something to help you",
	Long:  `grum clone means git clone but do something to help you.`,
	Args:  cobra.ExactArgs(1),
	RunE:  clone,
}

func clone(_ *cobra.Command, args []string) error {
	url := args[0]

	s, err := getNewRemoteUrl(url)
	if err != nil {
		return err
	}
	command := fmt.Sprintf("git clone %s", s)
	err = execx.Run(command)
	cobra.CheckErr(err)

	ep, err := Endpoint(url)
	cobra.CheckErr(err)

	pwd, _ := os.Getwd()
	repo, err := git.PlainOpen(filepath.Join(pwd, strings.TrimRight(filepath.Base(ep.Path), ".git")))
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
	return nil
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
