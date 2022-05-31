package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
)

func JoinErrors(errs []error, prefix string, delimiter string) string {
	messages := []string{}
	for _, err := range errs {
		messages = append(messages, fmt.Sprintf("%s%s", prefix, err))
	}
	return strings.Join(messages, delimiter)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func WalkUpFind(targetFilename string, currentPath string) (string, error) {
	candidate := filepath.Join(currentPath, targetFilename)

	if PathExists(candidate) {
		return candidate, nil
	}

	parentPath := filepath.Dir(currentPath)
	if parentPath == currentPath {
		return "", fmt.Errorf("Unable to find a %s file in the working path.", targetFilename)
	}

	return WalkUpFind(targetFilename, parentPath)
}

func KwargsToMap(kwargs []starlark.Tuple) map[string]starlark.Value {
	result := map[string]starlark.Value{}

	for _, kwarg := range kwargs {
		name := kwarg.Index(0).(starlark.String).GoString()
		result[name] = kwarg.Index(1)
	}

	return result
}
