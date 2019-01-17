package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Version struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type Versions []Version

func (v Version) IsNotEmpty() bool {
	if v.Version == "" {
		return false
	}

	return true
}

func (list Versions) Get(repoName string) (*Version, error) {
	for _, version := range list {
		if strings.Contains(version.Name, repoName) {
			return &version, nil
		}
	}

	return &Version{}, fmt.Errorf("failed finding %s", repoName)
}

func (list *Versions) AddReplace(ver Version) {
	var idx = -1

	for i, version := range *list {
		if strings.Contains(version.Name, ver.Name) {
			idx = i
			break
		}
	}

	if idx == -1 {
		*list = append(*list, ver)
		return
	}

	(*list)[idx] = ver
}

func versionsFile() io.ReadWriteCloser {
	var (
		srcPath string
		err     error
	)

	if config.VersionsPath == "" {
		srcPath = filepath.Join(".", "data", "versions.json")
	} else {
		srcPath, err = filepath.Abs(config.VersionsPath)
	}

	file, err := os.OpenFile(srcPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func readVersions(list *Versions) error {
	file := versionsFile()

	defer file.Close()

	err := json.NewDecoder(file).Decode(&list)
	if err != nil {
		err = errors.Wrap(err, "failed reading in versions list")
	}

	return err
}

func saveVersions(list *Versions) error {
	file := versionsFile()

	err := json.NewEncoder(file).Encode(list)
	if err != nil {
		return errors.Wrap(err, "failed writing to versions file")
	}

	return nil
}
