package version

import "testing"

func TestContentHashIsStableForEquivalentMaps(t *testing.T) {
	left, err := ContentHash(map[string]any{"network": "n", "species": []string{"H2", "O2"}})
	if err != nil {
		t.Fatal(err)
	}
	right, err := ContentHash(map[string]any{"species": []string{"H2", "O2"}, "network": "n"})
	if err != nil {
		t.Fatal(err)
	}
	if left != right {
		t.Fatalf("equivalent payloads produced different hashes: %s != %s", left, right)
	}
}
