package starverse

import (
	"os"
	"path/filepath"

	"github.com/jathu/starfig/internal/util"
)

const StarverseFilename string = "STARVERSE"
const StarverseDirThreadKey string = "starverseDir"

func FindStarverseDirectory() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	starversePath, err := util.WalkUpFind(StarverseFilename, currentPath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(starversePath), nil
}
