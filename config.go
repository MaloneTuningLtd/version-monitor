package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	DisableScheduler bool   `env:"DISABLE_SCHEDULER"`
	VersionsPath     string `env:"VERSIONS_PATH"`
	SourcesPath      string `env:"SOURCES_PATH"`
	SlackHook        string `env:"SLACK_HOOK"`
	GithubUser       string `env:"GITHUB_USER"`
	GithubToken      string `env:"GITHUB_TOKEN"`
}

func findSecrets() (values map[string]string) {
	values = make(map[string]string)

	var secrets []string
	const secretsPath = "/run/secrets"

	filepath.Walk(secretsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("SECRETS: accessing %s failed with: %s\n", path, err.Error())
			return nil
		}

		if !info.IsDir() {
			secrets = append(secrets, path)
			log.Printf("SECRETS: found %s secret\n", filepath.Base(path))
		}

		return nil
	})

	for _, secretPath := range secrets {
		f, err := os.Open(secretPath)
		if err != nil {
			log.Printf("SECRETS: failed to read %s\n", secretPath)
			continue
		}

		defer f.Close()

		if b, err := ioutil.ReadAll(f); err == nil {
			b = bytes.Replace(b, []byte("\n"), []byte(""), -1)
			values[filepath.Base(secretPath)] = string(b)
		}
	}

	return
}

func (c *Config) inferENV() (envMap map[string]string) {
	envMap = make(map[string]string)

	cType := reflect.TypeOf(c)
	cStruct := cType.Elem()

	for i := 0; i < cStruct.NumField(); i++ {
		field := cStruct.Field(i)

		if env, ok := field.Tag.Lookup("env"); ok && env != "" {
			envMap[field.Name] = env
		}
	}

	return
}

func loadConfig(c *Config) {
	keyMap := c.inferENV()
	secrets := findSecrets()

	cs := reflect.ValueOf(c).Elem()

	if cs.Kind() != reflect.Struct {
		return
	}

	for k, v := range keyMap {
		field := cs.FieldByName(k)

		secret, ok := secrets[v]
		if !ok && secret != "" {
			secret = os.Getenv(v)
		}

		if secret != "" && field.IsValid() && field.CanSet() {
			log.Printf("CONFIG: using %s\n", v)

			if field.Kind() == reflect.String {
				field.SetString(secret)
			}

			if field.Kind() == reflect.Bool && strings.Contains(strings.ToLower(secret), "yes") {
				field.SetBool(true)
			}
		}
	}
}
