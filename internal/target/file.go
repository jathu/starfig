package target

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/jathu/starfig/internal/util"
)

type FileTarget struct {
	StarverseDir string
	Package      string
	Filename     string
}

func (target FileTarget) String() string {
	return target.Path()
}

func (target FileTarget) Path() string {
	return filepath.Join(target.StarverseDir, target.Package, target.Filename)
}

func (target FileTarget) Exists() bool {
	return util.PathExists(target.Path())
}

func (target FileTarget) Target() string {
	return fmt.Sprintf("//%s", filepath.Join(target.Package, target.Filename))
}

func (target FileTarget) IsStarFile() bool {
	return filepath.Ext(target.Filename) == ".star"
}

func (target FileTarget) IsStarFigFile() bool {
	return target.Filename == StarfigFilename
}

func ParseFileTarget(starverseDir string, rawLabel string) (FileTarget, error) {
	target := FileTarget{starverseDir, "", ""}

	if strings.HasPrefix(rawLabel, "//") {
		rawLabel = rawLabel[2:]
		target.Package = path.Dir(rawLabel)
		target.Filename = filepath.Base(rawLabel)
	} else {
		return target, fmt.Errorf("Load source %s is invalid because it must be absolute.", rawLabel)
	}

	return target, nil
}
