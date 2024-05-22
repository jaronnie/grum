/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

func clone(cmd *cobra.Command, args []string) error {
	url := args[0]

	s, err := getNewRemoteUrl(url)
	if err != nil {
		return err
	}
	command := fmt.Sprintf("git clone %s", s)
	err = RunCommand(command)
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

func RunCommand(arg string) error {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case "darwin", "linux":
		cmd = exec.Command("sh", "-c", arg)
	case "windows":
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return fmt.Errorf("unexpected os: %v", goos)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Dir = pwd

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		_, _ = io.Copy(os.Stdout, stdout)
	}()

	go func() {
		_, _ = io.Copy(os.Stderr, stderr)
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
