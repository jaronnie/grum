/*
Copyright © 2023 jaronnie <jaron@jaronnie.com>

*/

package main

import "github.com/jaronnie/grum/cmd"

// ldflags
var (
	version = "1.3.0"
	commit  string
	date    string
)

func main() {
	cmd.Version = version
	cmd.Date = date
	cmd.Commit = commit

	cmd.Execute()
}
