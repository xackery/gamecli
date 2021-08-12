package gog

import (
	"context"
)

// Gog represents a gog instance
type Gog struct {
}

// New creates a new gog instance
func New(ctx context.Context) (*Gog, error) {
	s := &Gog{}
	return s, nil
}
