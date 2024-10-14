package config

import (
	"errors"
	"hscli/logging"
	"io/fs"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/urfave/cli/v2"
)

const EX_CONFIG = 78 // https://stackoverflow.com/questions/1101957/are-there-any-standard-exit-status-codes-in-linux

type Config struct {
	Root          string `yaml:"root"      env:"HS_ROOT" env-default:""`
	User          string `yaml:"user"      env:"HS_USER" env-default:""`
	Password      string `yaml:"password"  env:"HS_PASSWORD" env-default:""`
	CookieJarPath string `yaml:"cookiejar" env:"HS_COOKIEJAR" env-default:""`
}

// Attempts to load config from file and environment if any config parameter is not provided as a CLI argument
// Reads config from file (which in turn might be overwritten by environment in ReadConfig call)
// If the file can't be read only loads configuration from environment
func LoadConfig(cCtx *cli.Context, cfg *Config, cfgPath string) error {
	if cfg.Root != "" && cfg.User != "" && cfg.Password != "" && cfg.CookieJarPath != "" {
		return nil
	}

	if cfgPath != "" {
		logging.LogDebug("Attempting to load configuration from file %s", cfgPath)
		if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				logging.LogDebug("Failed reading configuration: %s", err)
				return cli.Exit("Failed loading configuration", EX_CONFIG)
			}
		}
	}

	logging.LogDebug("Attempting to load configuration from environment")
	if err := cleanenv.ReadEnv(cfg); err != nil {
		logging.LogDebug("Failed reading configuration from environment: %s", err)
		return cli.Exit("Failed loading configuration", EX_CONFIG)
	}

	// Check that all values were populated
	if cfg.Root == "" {
		return cli.Exit("Missing 'root' config parameter", EX_CONFIG)
	}
	if cfg.User == "" {
		return cli.Exit("Missing 'user' config parameter", EX_CONFIG)
	}
	if cfg.Password == "" {
		return cli.Exit("Missing 'password' config parameter", EX_CONFIG)
	}
	if cfg.CookieJarPath == "" {
		return cli.Exit("Missing 'cookiejar' config parameter", EX_CONFIG)
	}

	return nil
}
