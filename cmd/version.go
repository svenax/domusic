package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

const version = "1.7.1"

var (
	BuildTime = "unknown"
	GitSha1   = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of this program and of Lilypond",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s-%s built %s %s/%s\n", cmd.Root().Name(), version, GitSha1, BuildTime, runtime.GOOS, runtime.GOARCH)
		fmt.Println(lilyVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func lilyVersion() string {
	c := exec.Command("lilypond", "--version")
	out, _ := c.Output()

	return string(bytes.Split(out, []byte("\n"))[0])
}
