package main

import (
	"log"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

func sortMakeVersionTags(tags []TagReply) []*semver.Version {
	vs := make([]*semver.Version, len(tags))

	for i, r := range tags {
		v := strings.Replace(r.Name, "v", "", 1)

		ver, err := semver.NewVersion(v)
		if err != nil {
			errors.Wrapf(err, "error parsing version: %s", v)
		}

		vs[i] = ver
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	return vs
}

func IsNewer(recent, current Version) bool {
	recentVer, err := semver.NewVersion(recent.Version)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse recent version: %s", recent.Version)
		log.Println(err)

		return false
	}

	currentVer, err := semver.NewVersion(current.Version)
	if err != nil {
		if current.Version != "" {
			err = errors.Wrapf(err, "failed to parse current version: %s", current.Version)
			log.Println(err)
		}

		return true
	}

	return recentVer.GreaterThan(currentVer)
}