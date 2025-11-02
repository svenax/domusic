package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/urfave/cli/v3"
)

const version = "2.1.2"

var (
	BuildTime = "unknown"
	GitSha1   = "dev"
)

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Show version of this program and of Lilypond",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		// Version command doesn't need config
		fmt.Printf("domusic v%s (%s) built %s %s/%s\n", version, GitSha1, BuildTime, runtime.GOOS, runtime.GOARCH)
		fmt.Println(lilyVersion())
		return nil
	},
}

func lilyVersion() string {
	c := exec.Command("lilypond", "--version")
	out, _ := c.Output()

	return string(bytes.Split(out, []byte("\n"))[0])
}
