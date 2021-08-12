package gog

import (
	"context"
	"testing"
)

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	g, err := New(ctx)
	if err != nil {
		t.Fatalf("New: %+v", err)
	}
	err = g.Update(ctx)
	if err != nil {
		t.Fatalf("Update: %+v", err)
	}
}
