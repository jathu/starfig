package native

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jathu/starfig/internal/starverse"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestLoadProviderMissingTarget(t *testing.T) {
	thread := starlark.Thread{}
	thread.SetLocal(starverse.StarverseDirThreadKey, "/tmp")

	_, err := LoadProvider(&thread, "invalid_target")
	assert.ErrorContains(t, err,
		"Load source invalid_target is invalid because it must be absolute.")
}

func TestLoadProviderInvalidFile(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	thread := starlark.Thread{}
	thread.SetLocal(starverse.StarverseDirThreadKey, testStarverseDir)

	absolutePath := filepath.Join(testStarverseDir, "STARVERSE")
	_, err := LoadProvider(&thread, "//STARVERSE")

	expectedError := fmt.Sprintf("Only .star and STARFIG files can be loaded, %s is invalid.", absolutePath)
	assert.ErrorContains(t, err, expectedError)
}

func TestLoadProviderExecError(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	thread := starlark.Thread{}
	thread.SetLocal(starverse.StarverseDirThreadKey, testStarverseDir)

	absolutePath := filepath.Join(testStarverseDir, "invalid", "invalidSyntax.star")
	_, err := LoadProvider(&thread, "//invalid/invalidSyntax.star")

	expectedError := fmt.Sprintf("%s:1:1: undefined: this", absolutePath)
	assert.ErrorContains(t, err, expectedError)
}

func TestLoadProviderInvocationInStarFile(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	thread := starlark.Thread{Load: LoadProvider}
	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	thread.SetLocal(starverse.StarverseDirThreadKey, testStarverseDir)

	_, err := LoadProvider(&thread, "//invalid/invalidInvocation.star")

	expectedError := fmt.Sprintf("Schema types can only be instantiated in STARFIG files.")
	assert.ErrorContains(t, err, expectedError)
}

func TestLoadProviderNonResultInStarfigFile(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	thread := starlark.Thread{Load: LoadProvider}
	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	thread.SetLocal(starverse.StarverseDirThreadKey, testStarverseDir)

	_, err := LoadProvider(&thread, "//badfig/badfig.star")

	expectedError := fmt.Sprintf("STARFIG file can only contain schema instances.")
	assert.ErrorContains(t, err, expectedError)
}

func TestLoadProviderUpdatesContextManager(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	thread := starlark.Thread{Load: LoadProvider}
	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	thread.SetLocal(starverse.StarverseDirThreadKey, testStarverseDir)

	globals, err := LoadProvider(&thread, "//fruit/fruit.star")
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"Fruit"}, globals.Keys())

	contextSchemaNames := []string{}
	for _, schemaContextItem := range manager.builders {
		itemName, found := manager.GetSchemaName(schemaContextItem.SchemaDescriptor)
		assert.True(t, found)
		contextSchemaNames = append(contextSchemaNames, itemName)
	}
	assert.ElementsMatch(t, []string{"Fruit", "Color"}, contextSchemaNames)
}
