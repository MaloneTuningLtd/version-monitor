package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type GithubProvider struct {
	SourceProvider
	User  string
	Token string
}

func NewGithubProvider() *GithubProvider {
	return &GithubProvider{
		SourceProvider: SourceProvider{
			Name: "Github",
			Client: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
	}
}

func (g *GithubProvider) WithCredentials(user, pass string) {
	g.User = user
	g.Token = pass
}

func (g *GithubProvider) Fetch(repoName string) (repo Repository) {
	const githubURL = "https://api.github.com/"

	var (
		reply struct {
			Name        string `json:"full_name"`
			Description string `json:"description"`
			HTMLURL     string `json:"html_url"`
			TagURL      string `json:"tags_url"`
		}

		c = make(chan string)
	)

	client := g.SourceProvider.Client

	go func() {
		var reply []TagReply

		tagsURL := <-c

		req, err := http.NewRequest("GET", tagsURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		if g.User != "" && g.Token != "" {
			req.SetBasicAuth(g.User, g.Token)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&reply)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to parse tags"))
		}

		if len(reply) >= 1 {
			tags := sortMakeVersionTags(reply)
			c <- tags[0].Original()
		} else {
			c <- "None"
		}
	}()

	req, err := http.NewRequest("GET", githubURL+"repos/"+repoName, nil)
	if err != nil {
		log.Fatal(req)
	}

	if g.User != "" && g.Token != "" {
		req.SetBasicAuth(g.User, g.Token)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to parse response"))
	}

	log.Printf("LOG: checking %s \n", reply.Name)
	c <- reply.TagURL

	// set repo details
	repo.Description = reply.Description
	repo.URL = reply.HTMLURL
	repo.TagURL = reply.TagURL

	// set repo version
	if tag := <-c; tag != "None" {
		repo.Version.Name = reply.Name
		repo.Version.Version = tag
	}

	return
}
