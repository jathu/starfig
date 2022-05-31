package native

import (
	"fmt"
	"strings"

	"github.com/jathu/starfig/internal/starverse"
	"github.com/jathu/starfig/internal/target"
	"go.starlark.net/starlark"
)

var emptySrc interface{}

func LoadProvider(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	results := starlark.StringDict{}
	starverseDir := thread.Local(starverse.StarverseDirThreadKey).(string)

	fileTarget, err := target.ParseFileTarget(starverseDir, module)
	if err != nil {
		return results, err
	}

	if !fileTarget.IsStarFile() && !fileTarget.IsStarFigFile() {
		return results, fmt.Errorf(
			"Only .star and STARFIG files can be loaded, %s is invalid.", fileTarget.String())
	}

	globals, err := starlark.ExecFile(thread, fileTarget.Path(), emptySrc, Predeclared)
	if err != nil {
		return results, err
	}

	contextManager := thread.Local(SchemaContextManagerThreadKey).(SchemaContextManager)
	for name, value := range globals {
		schemaBuilder, ok := value.(*starlark.Builtin)
		if ok {
			contextManager.UpdateRecognizedSchema(schemaBuilder, name, fileTarget)
		}

		if fileTarget.IsStarFile() {
			_, ok := value.(SchemaResult)
			if ok {
				return results, fmt.Errorf("Schema types can only be instantiated in STARFIG files.")
			}
		} else if fileTarget.IsStarFigFile() {
			_, ok := value.(SchemaResult)
			if !ok {
				return results, fmt.Errorf("STARFIG file can only contain schema instances.")
			}
		}

		// Don't export globals that start with an underscore. Those are private.
		if !strings.HasPrefix(name, "_") {
			results[name] = value
		}
	}

	return results, nil
}
