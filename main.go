package main

import (
	"fmt"
	"log"
)

var (
	config Config
	srcs   Sources
	vers   Versions
)

func changeNotify(updated []*Version) {
	for _, version := range updated {
		vers.AddReplace(*version)
		notify(version)
	}
}

func check() (updated []*Version) {
	github := NewGithubProvider()
	github.WithCredentials(config.GithubUser, config.GithubToken)

	process := func(name string) {
		repo := github.Fetch(name)
		currentVersion, err := vers.Get(name)

		if err != nil {
			log.Println(err.Error())
			updated = append(updated, &repo.Version)

			return
		}

		if IsNewer(repo.Version, *currentVersion) {
			updated = append(updated, &repo.Version)
		}
	}

	for requireProvider, repos := range srcs {
		if requireProvider == "github" {
			for _, repoName := range repos {
				process(repoName)
			}
		}
	}

	return
}

func main() {
	fmt.Println("Github Version Monitor")

	// load in state
	loadConfig(&config)
	srcs = make(Sources)

	srcErr := readSources(&srcs)
	versErr := readVersions(&vers)

	if srcErr != nil {
		log.Fatal(srcErr)
	}

	if versErr != nil {
		log.Println(versErr)
	}

	// check for updated repositories
	updated := check()
	changeNotify(updated)

	err := saveVersions(&vers)
	if err != nil {
		log.Fatal(err)
	}
}
