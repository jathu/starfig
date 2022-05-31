package tester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func GetTestStarverseDir(t *testing.T) string {
	nativeDir, err := os.Getwd()
	assert.Nil(t, err)
	internalDir := filepath.Dir(nativeDir)
	// Use the example in <starfig-dir>/internal/tester/data
	return filepath.Join(internalDir, "tester", "data")
}

func MockBuiltinWithName(name string) *starlark.Builtin {
	mock := func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return starlark.None, nil
	}
	return starlark.NewBuiltin(name, mock)
}

func MockBuiltinWithCallback(callback func(args starlark.Tuple, kwargs []starlark.Tuple)) *starlark.Builtin {
	return starlark.NewBuiltin(uuid.New().String(), func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		callback(args, kwargs)
		return starlark.None, nil
	})
}

func MockBuiltin() *starlark.Builtin {
	return MockBuiltinWithName(uuid.New().String())
}

func MockFailingBuiltin(errorMessage string) *starlark.Builtin {
	mock := func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return starlark.None, fmt.Errorf(errorMessage)
	}
	return starlark.NewBuiltin(uuid.New().String(), mock)
}

// Mock a failing validation function — which returns a string error message
func MockFailingFunction(errorMessage string) *starlark.Builtin {
	mock := func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return starlark.String(errorMessage), nil
	}
	return starlark.NewBuiltin(uuid.New().String(), mock)
}

func AssertSameValidations(t *testing.T, expected *starlark.List, actual []starlark.Callable) {
	assert.Equal(t, expected.Len(), len(actual))

	for i := 0; i < expected.Len(); i++ {
		// They should be pointing to the same func
		assert.Same(t, expected.Index(i), actual[i])
	}
}
