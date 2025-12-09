package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration values for domusic
type Config struct {
	Root        string         `yaml:"root" env:"DOMUSIC_ROOT"`
	LyEditor    string         `yaml:"ly-editor" env:"DOMUSIC_LY_EDITOR"`
	LyViewer    string         `yaml:"ly-viewer" env:"DOMUSIC_LY_VIEWER"`
	FontInclude string         `yaml:"font-include" env:"DOMUSIC_FONT_INCLUDE"`
	Sync        SyncConfig     `yaml:"sync"`
	Template    TemplateConfig `yaml:"template"`
}

// SyncConfig holds configuration for the sync command
type SyncConfig struct {
	Server  string   `yaml:"server" env:"DOMUSIC_SYNC_SERVER"`
	User    string   `yaml:"user" env:"DOMUSIC_SYNC_USER"`
	Path    string   `yaml:"path" env:"DOMUSIC_SYNC_PATH"`
	SshKey  string   `yaml:"ssh-key" env:"DOMUSIC_SYNC_SSH_KEY"`
	Include []string `yaml:"include" env:"DOMUSIC_SYNC_INCLUDE"`
	Exclude []string `yaml:"exclude" env:"DOMUSIC_SYNC_EXCLUDE"`
}

// TemplateConfig holds configuration for common file templates used with Lilypond
type TemplateConfig struct {
	Common     string `yaml:"common"`
	Collection string `yaml:"collection"`
	Make       string `yaml:"make"`
}

var config *Config
var testMode bool // For testing - prevents config file loading

// loadFromEnv uses reflection to populate struct fields from environment variables
// based on the 'env' struct tags
func loadFromEnv(cfg *Config) {
	loadEnvIntoStruct(reflect.ValueOf(cfg).Elem(), reflect.TypeOf(cfg).Elem())

	// Special case: use EDITOR as fallback for LyEditor if LyEditor is empty
	if cfg.LyEditor == "" {
		if val := os.Getenv("EDITOR"); val != "" {
			cfg.LyEditor = val
		}
	}
}

// loadEnvIntoStruct recursively loads environment variables into a struct
func loadEnvIntoStruct(v reflect.Value, t reflect.Type) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Handle nested structs recursively
		if field.Kind() == reflect.Struct {
			loadEnvIntoStruct(field, fieldType.Type)
			continue
		}

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
}

// findLocalConfig searches for a config file in the current directory and parent directories
// up to the filesystem root. Returns the first config file found, or empty string if none found.
func findLocalConfig() string {
	configNames := []string{".domusic.yaml", ".domusic"}

	// Start from current working directory
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Search up the directory tree
	for {
		for _, name := range configNames {
			configPath := filepath.Join(dir, name)
			if _, err := os.Stat(configPath); err == nil {
				return configPath
			}
		}

		parent := filepath.Dir(dir)

		// Stop if we've reached the root
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// loadConfig loads configuration from file and environment variables
func loadConfig() (*Config, error) {
	cfg := &Config{}

	// Skip config file loading in test mode, but still load environment variables
	if testMode {
		loadFromEnv(cfg)
		return cfg, nil
	}

	// Try to load from config file
	if configPath == "" {
		// First try local config search
		configPath = findLocalConfig()

		// If no local config found, search in standard locations using XDG Base Directory Specification
		if configPath == "" {
			configPath, _ = xdg.SearchConfigFile("domusic/config.yaml")
		}

		// Fallback to legacy home directory locations if XDG search fails
		if configPath == "" {
			if home, err := os.UserHomeDir(); err == nil {
				legacyPaths := []string{
					filepath.Join(home, ".domusic.yaml"),
					filepath.Join(home, ".domusic"),
				}
				for _, path := range legacyPaths {
					if _, err := os.Stat(path); err == nil {
						configPath = path
						break
					}
				}
			}
		}
	}

	// Load from file if it exists
	if configPath != "" {
		if data, err := os.ReadFile(configPath); err == nil {
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
	config, err = loadConfig()
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
