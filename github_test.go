package main

import "testing"

func TestFetch(t *testing.T) {
	g := NewGithubProvider()
	repo := g.Fetch("nlopes/slack")

	if repo.Name == "" {
		t.Error("expected Name to not be empty")
	}

	if repo.URL == "" {
		t.Error("expected URL to not be empty")
	}

	if repo.TagURL == "" {
		t.Error("expected TagURL to not be empty")
	}
}
