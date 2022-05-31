package evaluator

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jathu/starfig/internal/target"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func TestEvaluateBuildTarget(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	buildTarget := target.BuildTarget{
		StarverseDir: testStarverseDir,
		Package:      "fruit",
		TargetName:   "apple",
	}
	evaluateResults, err := EvaluateBuildTarget(testStarverseDir, buildTarget)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(evaluateResults))

	expectedApple := new(starlark.Dict)
	expectedApple.SetKey(starlark.String("name"), starlark.String("Apple"))
	expectedApple.SetKey(starlark.String("colors"), starlark.NewList([]starlark.Value{
		makeColor(255, 0, 0),
		makeColor(0, 255, 0),
		makeColor(255, 255, 0),
	}))

	assert.Equal(t, buildTarget, evaluateResults[0].Target)
	same, err := expectedApple.CompareSameType(syntax.EQL, evaluateResults[0].Result.Evaluated, 10)
	assert.Nil(t, err)
	assert.True(t, same)
}

func TestEvaluateBuildTargetExecError(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	buildTarget := target.BuildTarget{
		StarverseDir: testStarverseDir,
		Package:      "evalerror",
		TargetName:   "apple",
	}
	_, err := EvaluateBuildTarget(testStarverseDir, buildTarget)
	absolutePath := filepath.Join(testStarverseDir, "evalerror", "STARFIG")
	expected := fmt.Sprintf("%s:4: Invalid field colors in Fruit: Expected list type but got 416.", absolutePath)
	assert.ErrorContains(t, err, expected)
}

func TestEvaluateBuildTargetExecErrorNonEval(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	buildTarget := target.BuildTarget{
		StarverseDir: testStarverseDir,
		Package:      "nonexistentPackage",
		TargetName:   "",
	}
	_, err := EvaluateBuildTarget(testStarverseDir, buildTarget)
	absolutePath := filepath.Join(testStarverseDir, "nonexistentPackage", "STARFIG")
	expected := fmt.Sprintf(`open %s: no such file or directory`, absolutePath)
	assert.ErrorContains(t, err, expected)
}

func TestEvaluateBuildTargetSpread(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	buildTarget := target.BuildTarget{
		StarverseDir: testStarverseDir,
		Package:      "trait",
		TargetName:   "...",
	}
	evaluateResults, err := EvaluateBuildTarget(testStarverseDir, buildTarget)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(evaluateResults))

	expectedGreen := makeColor(0, 255, 0)
	sameGreen, err := expectedGreen.CompareSameType(syntax.EQL, evaluateResults[0].Result.Evaluated, 2)
	assert.Nil(t, err)
	assert.True(t, sameGreen)
	assert.Equal(t, "green", evaluateResults[0].Target.TargetName)

	expectedRed := makeColor(255, 0, 0)
	sameRed, err := expectedRed.CompareSameType(syntax.EQL, evaluateResults[1].Result.Evaluated, 2)
	assert.Nil(t, err)
	assert.True(t, sameRed)
	assert.Equal(t, "red", evaluateResults[1].Target.TargetName)

	expectedYellow := makeColor(255, 255, 0)
	sameYellow, err := expectedYellow.CompareSameType(syntax.EQL, evaluateResults[2].Result.Evaluated, 2)
	assert.Nil(t, err)
	assert.True(t, sameYellow)
	assert.Equal(t, "yellow", evaluateResults[2].Target.TargetName)
}

func TestEvaluateBuildTargetNotFoundTarget(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	buildTarget := target.BuildTarget{
		StarverseDir: testStarverseDir,
		Package:      "fruit",
		TargetName:   "tangerine",
	}
	_, err := EvaluateBuildTarget(testStarverseDir, buildTarget)
	assert.ErrorContains(t, err, "//fruit:tangerine not found.")
}

func TestEvaluateBuildTargetNotSchemaResult(t *testing.T) {
	testStarverseDir := tester.GetTestStarverseDir(t)
	buildTarget := target.BuildTarget{
		StarverseDir: testStarverseDir,
		Package:      "badfig",
		TargetName:   "BadFig",
	}
	_, err := EvaluateBuildTarget(testStarverseDir, buildTarget)
	assert.ErrorContains(t, err, "//badfig:BadFig is not a schema result.")
}

// MARK: - Helpers

func makeColor(red int, green int, blue int) *starlark.Dict {
	result := new(starlark.Dict)
	result.SetKey(starlark.String("red"), starlark.MakeInt(red))
	result.SetKey(starlark.String("green"), starlark.MakeInt(green))
	result.SetKey(starlark.String("blue"), starlark.MakeInt(blue))
	return result
}
