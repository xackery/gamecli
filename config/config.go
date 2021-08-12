package config

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/jbsmith7741/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	mu  sync.RWMutex
	cfg *Config
)

const (
	configFileName = "gamecli.conf"
)

// Config represents a configuration parse
type Config struct {
	Debug           bool   `toml:"debug" desc:"TalkEQ Configuration\n\n# Debug messages are displayed. This will cause console to be more verbose, but also more informative"`
	SteamBinaryPath string `toml:"steam_binary_path" desc:"Path where steam is installed"`
	SteamAppPath    string `toml:"steam_app_path" desc:"Path where steam stores cache data"`
}

func Get(ctx context.Context) (Config, error) {
	mu.Lock()
	defer mu.Unlock()
	var err error
	var newConfig Config
	if cfg != nil {
		return *cfg, nil
	}

	cfg, err = new(ctx)
	if err != nil {
		return newConfig, fmt.Errorf("new: %w", err)
	}
	return *cfg, nil
}

func Set(ctx context.Context, newConfig Config) error {
	mu.Lock()
	defer mu.Unlock()
	cfg = &newConfig
	return newConfig.save()
}

// new creates a new configuration
func new(ctx context.Context) (*Config, error) {
	var f *os.File
	cfg := Config{}

	isNewConfig := false
	fi, err := os.Stat(configFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("config info: %w", err)
		}
		f, err = os.Create(configFileName)
		if err != nil {
			return nil, fmt.Errorf("create talkeq.conf: %w", err)
		}
		fi, err = os.Stat(configFileName)
		if err != nil {
			return nil, fmt.Errorf("new config info: %w", err)
		}
		isNewConfig = true
	}
	if !isNewConfig {
		f, err = os.Open(configFileName)
		if err != nil {
			return nil, fmt.Errorf("open config: %w", err)
		}
	}

	defer f.Close()
	if fi.IsDir() {
		return nil, fmt.Errorf("%s is a directory, should be a file", configFileName)
	}

	if isNewConfig {
		enc := toml.NewEncoder(f)
		enc.Encode(getDefaultConfig())

		log.Debug().Msgf("a new %s was created", configFileName)
		if runtime.GOOS == "windows" {
			option := ""
			fmt.Println("press a key then enter to exit.")
			fmt.Scan(&option)
		}
		os.Exit(0)
	}

	_, err = toml.DecodeReader(f, &cfg)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", configFileName, err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	err = cfg.Verify()
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}

	return &cfg, nil
}

func (c *Config) save() error {
	fw, err := os.Create(configFileName)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer fw.Close()

	enc := toml.NewEncoder(fw)
	err = enc.Encode(c)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

// Verify returns an error if configuration appears off
func (c *Config) Verify() error {

	return nil
}

func getDefaultConfig() Config {
	cfg := Config{
		Debug: true,
	}
	return cfg
}
