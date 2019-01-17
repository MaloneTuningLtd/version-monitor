package main

import "net/http"

type Provider interface {
	Fetch(repoName string) Repository
}

type SourceProvider struct {
	Name   string
	Client *http.Client
}

type TagReply struct {
	Name string `json:"name"`
}

type Repository struct {
	Version

	Description string
	URL         string
	TagURL      string
}
