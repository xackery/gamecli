package steam

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// Run starts an appID
func (s *Steam) Run(ctx context.Context, appID string, name string, isDirect bool) error {

	acfs, err := s.parse(ctx)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	acf := find(appID, name, acfs)
	if acf == nil {
		return fmt.Errorf("no games found matching criteria (appid: %s, name: %s)", appID, name)
	}

	c, err := s.prepareCommand(ctx, acf, isDirect)
	if err != nil {
		return fmt.Errorf("prepareCommand: %w", err)
	}

	log.Info().Msgf("executing %s with args %s", c.Path, strings.Join(c.Args, " "))
	err = c.Start()
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}

	err = c.Wait()
	if err != nil && err.Error() != "exit status 42" {
		return fmt.Errorf("wait: %w", err)
	}

	log.Debug().Msgf("waiting for process to start")
	return nil
}
