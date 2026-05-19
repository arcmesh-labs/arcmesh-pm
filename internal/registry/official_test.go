package registry

import (
	"testing"
)

func TestSearchOfficial(t *testing.T) {
	results, err := SearchOfficial("github")
	if err != nil {
		t.Fatalf("SearchOfficial returned error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result for query 'github', got 0")
	}
	t.Logf("found %d results", len(results))
	for _, r := range results {
		if r.Source != "official" {
			t.Errorf("entry %q: Source = %q, want 'official'", r.Name, r.Source)
		}
		if r.FullName == "" {
			t.Errorf("entry %q: FullName is empty", r.Name)
		}
	}
}
