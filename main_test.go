package main

import (
	"testing"
)

func TestStashAtIndex(t *testing.T) {
	stash := stashAtIndex(1)
	expected := "stash@{1}"

	if stash != expected {
		t.Fatalf("Stash was incorrect, got: %s, want: %s", stash, expected)
	}
}
