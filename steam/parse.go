package steam

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (s *Steam) parse(ctx context.Context) ([]*Acf, error) {
	steamAppsDir, err := s.steamAppsDir()
	if err != nil {
		return nil, fmt.Errorf("steamAppsDir: %w", err)
	}

	libs, err := parseVdf(ctx, fmt.Sprintf("%s/libraryfolders.vdf", steamAppsDir))
	if err != nil {
		return nil, fmt.Errorf("parseVdf: %w", err)
	}
	libs = append(libs, steamAppsDir)
	acfs := []*Acf{}
	for _, path := range libs {
		as, err := parseAcfDir(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("parseAcfDir: %w", err)
		}
		acfs = append(acfs, as...)
	}
	log.Debug().Msgf("found %d total games", len(acfs))
	return acfs, nil
}
