package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

var syncCmd = &cli.Command{
	Name:  "sync",
	Usage: "Sync files from _output directory to external web server",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"n"},
			Usage:   "perform a trial run with no changes made",
		},
		&cli.BoolFlag{
			Name:    "progress",
			Aliases: []string{"p"},
			Usage:   "show progress during transfer",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "increase verbosity",
		},
		&cli.BoolFlag{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete extraneous files from destination dirs",
		},
		&cli.StringFlag{
			Name:    "exclude",
			Aliases: []string{"e"},
			Usage:   "exclude files matching pattern",
		},
		&cli.StringFlag{
			Name:    "include",
			Aliases: []string{"i"},
			Usage:   "include files matching pattern",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		syncer := &syncer{cmd}
		return syncer.run()
	},
}

type syncer struct {
	cmd *cli.Command
}

func (s *syncer) run() error {
	// Validate configuration
	if err := s.validateConfig(); err != nil {
		return err
	}

	// Build source path
	sourcePath := pathFromRoot(outputDir) + "/"

	// Check if source directory exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source directory %s does not exist", sourcePath)
	}

	// Build destination path
	destPath := s.buildDestinationPath()

	// Build rsync command
	args := s.buildRsyncArgs(sourcePath, destPath)

	fmt.Printf("Syncing %s to %s\n", sourcePath, destPath)
	if s.cmd.Bool("dry-run") {
		fmt.Println("Dry run mode - no changes will be made")
	}

	// Execute rsync
	rsyncCmd := exec.Command("rsync", args...)
	rsyncCmd.Stdout = os.Stdout
	rsyncCmd.Stderr = os.Stderr

	if s.cmd.Bool("verbose") {
		fmt.Printf("Executing: rsync %s\n", strings.Join(args, " "))
	}

	return rsyncCmd.Run()
}

func (s *syncer) validateConfig() error {
	config := GetConfig()
	server := config.Sync.Server
	if server == "" {
		return fmt.Errorf("sync-server not configured - please set it in your config file or DOMUSIC_SYNC_SERVER environment variable")
	}

	user := config.Sync.User
	if user == "" {
		return fmt.Errorf("sync-user not configured - please set it in your config file or DOMUSIC_SYNC_USER environment variable")
	}

	path := config.Sync.Path
	if path == "" {
		return fmt.Errorf("sync-path not configured - please set it in your config file or DOMUSIC_SYNC_PATH environment variable")
	}

	return nil
}

func (s *syncer) buildDestinationPath() string {
	config := GetConfig()
	user := config.Sync.User
	server := config.Sync.Server
	path := config.Sync.Path

	return fmt.Sprintf("%s@%s:%s", user, server, path)
}

func (s *syncer) buildRsyncArgs(source, dest string) []string {
	args := []string{
		"-az", // archive mode, compress
	}
	config := GetConfig()

	// Add SSH key if configured
	if sshKey := config.Sync.SshKey; sshKey != "" {
		// Expand tilde to home directory if needed
		if strings.HasPrefix(sshKey, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				sshKey = filepath.Join(home, sshKey[2:])
			}
		}
		args = append(args, "-e", fmt.Sprintf("ssh -i %s", sshKey))
	}

	if s.cmd.Bool("dry-run") {
		args = append(args, "--dry-run")
	}
	if s.cmd.Bool("progress") || s.cmd.Bool("verbose") {
		args = append(args, "--progress")
	}
	if s.cmd.Bool("delete") {
		args = append(args, "--delete")
	}

	if exclude := s.cmd.String("exclude"); exclude != "" {
		args = append(args, "--exclude", exclude)
	}
	for _, exclude := range config.Sync.Exclude {
		if exclude != "" {
			args = append(args, "--exclude", exclude)
		}
	}
	if include := s.cmd.String("include"); include != "" {
		args = append(args, "--include", include)
	}
	for _, include := range config.Sync.Include {
		if include != "" {
			args = append(args, "--include", include)
		}
	}

	// Add verbose flag (rsync has different levels)
	if s.cmd.Bool("verbose") {
		args = append(args, "-v")
	}

	// Add source and destination
	args = append(args, source, dest)

	return args
}
