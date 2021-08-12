package steam

import (
	"context"
	"testing"
)

func TestDiscoverSteamBinary(t *testing.T) {
	_, err := New(context.Background())
	if err != nil {
		t.Fatalf("new: %v", err)
	}

}
