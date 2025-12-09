package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/urfave/cli/v3"
)

const (
	version           = "2.2.0"
	lowestLilyVersion = "2.24.0"
)

var (
	buildTime = "unknown"
	gitSha1   = "dev"
)

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Show version of this program and of Lilypond",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Printf("domusic v%s (%s) built %s %s/%s\n", version, gitSha1, buildTime, runtime.GOOS, runtime.GOARCH)
		fmt.Println(lilyVersion())
		fmt.Println("Version cmd:", lowestLilyVersion)
		fmt.Println("Config path:", configPath)
		return nil
	},
}

func lilyVersion() string {
	c := exec.Command("lilypond", "--version")
	out, err := c.Output()
	if err != nil {
		return "(lilypond not found or failed to run)"
	}
	if len(out) == 0 {
		return "(no version output)"
	}
	return string(bytes.Split(out, []byte("\n"))[0])
}
