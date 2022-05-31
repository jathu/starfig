package target

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jathu/starfig/internal/util"
	"github.com/sirupsen/logrus"
)

const StarfigFilename string = "STARFIG"

type BuildTarget struct {
	StarverseDir string
	Package      string
	TargetName   string
}

func (target BuildTarget) Target() string {
	return fmt.Sprintf("//%s:%s", target.Package, target.TargetName)
}

func (target BuildTarget) Path() string {
	return filepath.Join(target.StarverseDir, target.Package, StarfigFilename)
}

func ParseBuildTarget(
	starverseDir string, rawTargetInput string) ([]BuildTarget, error) {
	if !strings.HasPrefix(rawTargetInput, "//") {
		return []BuildTarget{}, fmt.Errorf(
			`"%s" is invalid because a target must start with //.`, rawTargetInput)
	}

	rawTargets := []string{}
	if strings.HasSuffix(rawTargetInput, "...") {
		searchDir := filepath.Join(starverseDir, rawTargetInput[:len(rawTargetInput)-3])
		packages, err := findStarfigPackages(starverseDir, searchDir)
		if err != nil {
			return []BuildTarget{}, err
		}
		rawTargets = append(rawTargets, packages...)
	} else {
		rawTargets = append(rawTargets, rawTargetInput)
	}

	targets := []BuildTarget{}
	for _, rawTarget := range rawTargets {
		rawTarget = rawTarget[2:]
		var colon_index, colon_count int
		for i, letter_ord := range rawTarget {
			if string(letter_ord) == ":" {
				colon_index = i
				colon_count += 1
			}
		}

		if colon_count != 1 {
			return []BuildTarget{}, fmt.Errorf("Invalid target %s.", rawTarget)
		}

		target := BuildTarget{starverseDir, "", ""}
		target.Package = rawTarget[:colon_index]
		target.TargetName = rawTarget[colon_index+1:]

		if util.PathExists(target.Path()) == false {
			return []BuildTarget{}, fmt.Errorf(
				"%s file for %s does not exist.", StarfigFilename, target.Target())
		}

		targets = append(targets, target)
	}

	return targets, nil
}

func findStarfigPackages(starverseDir string, searchRoot string) ([]string, error) {
	logrus.Debugf("searching for packages in %s", searchRoot)

	packages := []string{}
	err := filepath.WalkDir(searchRoot, func(path string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !entry.IsDir() && entry.Name() == StarfigFilename {
			rawTarget, err := filepath.Rel(starverseDir, filepath.Dir(path))
			if err != nil {
				return err
			}
			logrus.Debugf("found package: %s", rawTarget)
			packages = append(packages, fmt.Sprintf("//%s:...", rawTarget))
		}
		return nil
	})

	return packages, err
}
