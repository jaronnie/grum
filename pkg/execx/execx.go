package execx

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func Run(arg string) error {
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
