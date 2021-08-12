package steam

import (
	"context"
	"fmt"
)

// Update will request from steam an update of local games supported
func (s *Steam) Update(ctx context.Context) error {
	acfs, err := s.parse(ctx)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	for _, a := range acfs {
		fmt.Println(a.AppID, a.Name, a.StateName)
	}
	return nil
}
