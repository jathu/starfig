package evaluator

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/native"
	"github.com/jathu/starfig/internal/starverse"
	"github.com/jathu/starfig/internal/target"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
)

var emptySrc interface{}

type EvaluateResult struct {
	Target target.BuildTarget
	Result native.SchemaResult
}

func EvaluateBuildTarget(starverseDir string, buildTarget target.BuildTarget) ([]EvaluateResult, error) {
	thread := &starlark.Thread{
		Name: fmt.Sprintf("EvaluateBuildTarget:%s", uuid.New()),
		Load: native.LoadProvider,
		Print: func(thread *starlark.Thread, msg string) {
			logrus.Info(msg)
		},
	}
	thread.SetLocal(starverse.StarverseDirThreadKey, starverseDir)
	thread.SetLocal(native.SchemaContextManagerThreadKey, native.NewSchemaContextManager())

	globals, err := starlark.ExecFile(thread, buildTarget.Path(), emptySrc, native.Predeclared)
	if err != nil {
		evalErr, ok := err.(*starlark.EvalError)
		if ok {
			var lastFrame starlark.CallFrame
			for _, callstack := range evalErr.CallStack {
				// https://github.com/google/starlark-go/blob/d1966c6b9fcd6631f48f5155f47afcd7adcc78c2/starlark/eval.go#L197
				if callstack.Pos.Filename() != "<builtin>" {
					lastFrame = callstack
					break
				}
			}
			pos := lastFrame.Pos
			return []EvaluateResult{}, fmt.Errorf(
				"%s:%d: %s", pos.Filename(), pos.Line, evalErr.Msg)
		} else {
			return []EvaluateResult{}, err
		}
	}

	results := []EvaluateResult{}

	if buildTarget.TargetName == "..." {
		for _, targetName := range globals.Keys() {
			thisBuildTarget := target.BuildTarget{
				StarverseDir: buildTarget.StarverseDir,
				Package:      buildTarget.Package,
				TargetName:   targetName,
			}
			value := globals[targetName]
			result, ok := value.(native.SchemaResult)
			if ok {
				results = append(results, EvaluateResult{
					Target: thisBuildTarget,
					Result: result,
				})
			}
		}
	} else {
		value, ok := globals[buildTarget.TargetName]
		if !ok {
			return []EvaluateResult{}, fmt.Errorf("%s not found.", buildTarget.Target())
		}
		result, ok := value.(native.SchemaResult)
		if !ok {
			return []EvaluateResult{}, fmt.Errorf("%s is not a schema result.", buildTarget.Target())
		}
		results = append(results, EvaluateResult{Target: buildTarget, Result: result})
	}

	return results, nil
}
