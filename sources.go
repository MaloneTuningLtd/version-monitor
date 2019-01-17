package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Sources map[string][]string

func sourcesFile() io.ReadWriteCloser {
	var (
		srcPath string
		err     error
	)

	if config.SourcesPath == "" {
		srcPath = filepath.Join(".", "data", "sources.json")
	} else {
		srcPath, err = filepath.Abs(config.SourcesPath)
	}

	file, err := os.OpenFile(srcPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func readSources(src *Sources) error {
	file := sourcesFile()
	defer file.Close()

	err := json.NewDecoder(file).Decode(src)
	if err != nil {
		err = errors.Wrap(err, "failed reading in sources")
		return err
	}

	return nil
}
