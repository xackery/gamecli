package steam

import (
	"context"
	"fmt"
)

// Steam represents a steam instance
type Steam struct {
	// use s.steamAppsDir() to access
	steamAppsDirCache string
}

// New creates a new steam instance
func New(ctx context.Context) (*Steam, error) {
	s := &Steam{}
	err := s.findBinary(ctx)
	if err != nil {
		return nil, fmt.Errorf("findSteamBinary: %w", err)
	}
	err = s.findAppPath(ctx)
	if err != nil {
		return nil, fmt.Errorf("findAppPath: %w", err)
	}
	return s, nil
}
