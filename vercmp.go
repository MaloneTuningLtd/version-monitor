package main

import (
	"log"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

func sortMakeVersionTags(tags []TagReply) []*semver.Version {
	var vs []*semver.Version

	for _, r := range tags {
		v := strings.Replace(r.Name, "v", "", 1)

		ver, err := semver.NewVersion(v)
		if err != nil {
			err = errors.Wrapf(err, "WARNING: error parsing version: %s", v)
			log.Println(err)
			continue
		}

		// vs[i] = ver
		vs = append(vs, ver)
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	return vs
}

func IsNewer(recent, current Version) bool {
	if recent.Version == "" || current.Version == "" {
		return false
	}

	recentVer, err := semver.NewVersion(recent.Version)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse (%s) recent version: %s", recent.Name, recent.Version)
		log.Println(err)

		return false
	}

	currentVer, err := semver.NewVersion(current.Version)
	if err != nil {
		if current.Version != "" {
			err = errors.Wrapf(err, "failed to parse (%s) current version: %s", current.Name, current.Version)
			log.Println(err)
		}

		return true
	}

	return recentVer.GreaterThan(currentVer)
}
