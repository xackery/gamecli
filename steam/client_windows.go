package steam

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/xackery/gamecli/config"
	"golang.org/x/sys/windows/registry"
)

func (s *Steam) findBinary(ctx context.Context) error {
	path, err := findSteamPathInRegistry()
	if err != nil {
		return fmt.Errorf("findSteamPathInRegistry")
	}

	paths := []string{
		fmt.Sprintf("%s/%s", path, "steam.exe"),
		`C:/Program Files/Steam/steam.exe`,
	}
	var lastError error
	for _, path := range paths {
		log.Debug().Msgf("os.Stat %s", path)
		_, err := os.Stat(path)
		if err == nil {
			log.Debug().Msgf("found steam binary at %s", path)

			cfg, err := config.Get(ctx)
			if err != nil {
				return fmt.Errorf("config.Get: %w", err)
			}
			cfg.SteamBinaryPath = path
			err = config.Set(ctx, cfg)
			if err != nil {
				return fmt.Errorf("config.Set: %w", err)
			}
			return nil
		}
		if err != nil && !os.IsNotExist(err) {
			log.Debug().Err(err).Msgf("failed stat")
			lastError = err
			continue
		}
	}
	if lastError != nil {
		return fmt.Errorf("fallback: %w", lastError)
	}
	return fmt.Errorf("no paths valid")
}

func (s *Steam) findAppPath(ctx context.Context) error {
	path, err := findSteamPathInRegistry()
	if err != nil {
		return fmt.Errorf("findSteamPathInRegistry")
	}

	paths := []string{
		fmt.Sprintf("%s", path),
		`C:/Program Files/Steam`,
	}
	var lastError error
	for _, path := range paths {
		log.Debug().Msgf("os.Stat %s", path)
		_, err := os.Stat(path)
		if err == nil {
			log.Debug().Msgf("found steam app path at %s", path)

			cfg, err := config.Get(ctx)
			if err != nil {
				return fmt.Errorf("config.Get: %w", err)
			}
			cfg.SteamAppPath = path
			err = config.Set(ctx, cfg)
			if err != nil {
				return fmt.Errorf("config.Set: %w", err)
			}
			return nil
		}
		if err != nil && !os.IsNotExist(err) {
			log.Debug().Err(err).Msgf("failed stat")
			lastError = err
			continue
		}
	}
	if lastError != nil {
		return fmt.Errorf("fallback: %w", lastError)
	}
	return fmt.Errorf("no paths valid")
}

func (s *Steam) steamAppsDir() (string, error) {
	mu.Lock()
	defer mu.Unlock()
	if s.steamAppsDirCache != "" {
		return s.steamAppsDirCache, nil
	}
	cfg, err := config.Get(context.Background())
	if err != nil {
		return "", fmt.Errorf("config.Get: %w", err)
	}

	s.steamAppsDirCache = fmt.Sprintf("%s/SteamApps", cfg.SteamAppPath)
	fi, err := os.Stat(s.steamAppsDirCache)
	if err != nil {
		return "", fmt.Errorf("stat steamAppsDir: %w", err)
	}
	if !fi.IsDir() {
		return "", fmt.Errorf("steamAppsDir is not a directory")
	}
	return s.steamAppsDirCache, nil
}

func (s *Steam) prepareCommand(ctx context.Context, acf *Acf, isDirect bool) (exec.Cmd, error) {
	buf := &bytes.Buffer{}
	mw := io.MultiWriter(os.Stdout, buf)
	c := exec.Cmd{}

	cfg, err := config.Get(ctx)
	if err != nil {
		return c, fmt.Errorf("config.Get: %w", err)
	}

	binaryName := filepath.Base(cfg.SteamBinaryPath)
	binaryDir := strings.TrimSuffix(cfg.SteamBinaryPath, binaryName)

	c.Path = binaryName
	c.Dir = binaryDir
	c.Args = []string{binaryName, fmt.Sprintf("steam://rungameid/%s", acf.AppID)}
	c.Stdout = mw
	c.Stderr = mw

	if isDirect {
		c.Path = "/usr/bin/open"
		c.Args = []string{"open", fmt.Sprintf("%s/%s.app", acf.GamePath, acf.Installdir)}
	}
	return c, nil
}

func findSteamPathInRegistry() (string, error) {
	path, err := findSteam(false)
	if err != nil {
		path, err = findSteam(true)
		if err != nil {
			return "", fmt.Errorf("findSteam: %w", err)
		}
		return path, nil
	}
	return path, nil
}

func findSteam(is32Bit bool) (string, error) {
	path := `SOFTWARE\Wow6432Node\Valve\Steam`
	if is32Bit {
		path = `SOFTWARE\Valve\Steam`
	}
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("OpenKey: %w", err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("InstallPath")
	if err != nil {
		return "", fmt.Errorf("GetStringValue: %w", err)
	}
	return s, nil
}
