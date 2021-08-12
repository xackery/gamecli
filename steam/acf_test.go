package steam

import (
	"context"
	"testing"
)

func TestParseAcf(t *testing.T) {
	ctx := context.Background()
	a, err := parseAcf(ctx, "acf_test.txt")
	if err != nil {
		t.Fatalf("parseAcf: %v", err)
	}
	if a == nil {
		t.Fatalf("expected acf return, but it was empty")
	}
	//fmt.Printf("%+v\n", a)
}
