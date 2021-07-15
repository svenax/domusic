package cmd

import (
	"fmt"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "domusic",
	Short: "Handles the music library at svenax.net",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	errExit(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	filename := ".domusic.yaml"
	if runtime.GOOS == "windows" {
		filename = "domusic.ini"
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", filename))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		errExit(err)

		viper.AddConfigPath(home)

		viper.SetConfigName(".domusic")

		if runtime.GOOS == "windows" {
			viper.SetConfigName("domusic")
		}
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	errExit(err)
}
