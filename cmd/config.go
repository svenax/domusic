package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration values for domusic
type Config struct {
	Root        string   `yaml:"root" env:"DOMUSIC_ROOT"`
	LyEditor    string   `yaml:"ly-editor" env:"DOMUSIC_LY_EDITOR"`
	LyViewer    string   `yaml:"ly-viewer" env:"DOMUSIC_LY_VIEWER"`
	EnNotebook  string   `yaml:"en-notebook" env:"DOMUSIC_EN_NOTEBOOK"`
	FontInclude string   `yaml:"font-include" env:"DOMUSIC_FONT_INCLUDE"`
	SyncServer  string   `yaml:"sync-server" env:"DOMUSIC_SYNC_SERVER"`
	SyncUser    string   `yaml:"sync-user" env:"DOMUSIC_SYNC_USER"`
	SyncPath    string   `yaml:"sync-path" env:"DOMUSIC_SYNC_PATH"`
	SyncSSHKey  string   `yaml:"sync-ssh-key" env:"DOMUSIC_SYNC_SSH_KEY"`
	SyncInclude []string `yaml:"sync-include" env:"DOMUSIC_SYNC_INCLUDE"`
	SyncExclude []string `yaml:"sync-exclude" env:"DOMUSIC_SYNC_EXCLUDE"`
}

var config *Config
var testMode bool // For testing - prevents config file loading

// loadFromEnv uses reflection to populate struct fields from environment variables
// based on the 'env' struct tags
func loadFromEnv(cfg *Config) {
	v := reflect.ValueOf(cfg).Elem()
	t := reflect.TypeOf(cfg).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the env tag
		envVar := fieldType.Tag.Get("env")
		if envVar == "" {
			continue
		}

		// Get environment variable value
		envValue := os.Getenv(envVar)
		if envValue == "" {
			continue
		}

		// Handle string fields
		if field.Kind() == reflect.String {
			field.SetString(envValue)
		} else if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.String {
			// Handle string slices - split by comma
			parts := strings.Split(envValue, ",")
			// Trim whitespace from each part
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			field.Set(reflect.ValueOf(parts))
		}
	}

	// Special case: use EDITOR as fallback for LyEditor if LyEditor is empty
	if cfg.LyEditor == "" {
		if val := os.Getenv("EDITOR"); val != "" {
			cfg.LyEditor = val
		}
	}
}

// loadConfig loads configuration from file and environment variables
func loadConfig(configFile string) (*Config, error) {
	cfg := &Config{}

	// Skip config file loading in test mode, but still load environment variables
	if testMode {
		loadFromEnv(cfg)
		return cfg, nil
	}

	// Try to load from config file
	if configFile == "" {
		// Determine default config file
		home, err := os.UserHomeDir()
		if err == nil {
			if runtime.GOOS == "windows" {
				configFile = filepath.Join(home, "domusic.yaml")
				if _, err := os.Stat(configFile); os.IsNotExist(err) {
					configFile = filepath.Join(home, "domusic.ini")
				}
			} else {
				configFile = filepath.Join(home, ".domusic.yaml")
				if _, err := os.Stat(configFile); os.IsNotExist(err) {
					configFile = filepath.Join(home, ".domusic")
				}
			}
		}
	}

	// Load from file if it exists
	if configFile != "" {
		if data, err := os.ReadFile(configFile); err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
		}
		// Ignore file not found errors - we'll use environment variables and defaults
	}

	// Override with environment variables using reflection
	loadFromEnv(cfg)

	return cfg, nil
}

// initConfig initializes the global config
func initConfig() {
	var err error
	config, err = loadConfig(cfgFile)
	if err != nil {
		// Don't exit on config errors, just use defaults
		config = &Config{}
	}
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	if config == nil {
		initConfig()
	}
	return config
}

// Configuration getter methods for backward compatibility
func getString(key string) string {
	cfg := GetConfig()
	switch key {
	case "root":
		return cfg.Root
	case "ly-editor":
		return cfg.LyEditor
	case "ly-viewer":
		return cfg.LyViewer
	case "en-notebook":
		return cfg.EnNotebook
	case "font-include":
		return cfg.FontInclude
	case "sync-server":
		return cfg.SyncServer
	case "sync-user":
		return cfg.SyncUser
	case "sync-path":
		return cfg.SyncPath
	case "sync-ssh-key":
		return cfg.SyncSSHKey
	default:
		return ""
	}
}

func isSet(key string) bool {
	return getString(key) != ""
}

func setString(key, value string) {
	cfg := GetConfig()
	switch key {
	case "root":
		cfg.Root = value
	case "ly-editor":
		cfg.LyEditor = value
	case "ly-viewer":
		cfg.LyViewer = value
	case "en-notebook":
		cfg.EnNotebook = value
	case "font-include":
		cfg.FontInclude = value
	case "sync-server":
		cfg.SyncServer = value
	case "sync-user":
		cfg.SyncUser = value
	case "sync-path":
		cfg.SyncPath = value
	case "sync-ssh-key":
		cfg.SyncSSHKey = value
	}
}

// getStringSlice returns a string slice configuration value
func getStringSlice(key string) []string {
	cfg := GetConfig()
	switch key {
	case "sync-include":
		return cfg.SyncInclude
	case "sync-exclude":
		return cfg.SyncExclude
	default:
		return []string{}
	}
}

// setStringSlice sets a string slice configuration value
func setStringSlice(key string, value []string) {
	cfg := GetConfig()
	switch key {
	case "sync-include":
		cfg.SyncInclude = value
	case "sync-exclude":
		cfg.SyncExclude = value
	}
}
