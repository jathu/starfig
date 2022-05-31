package native

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
)

var Predeclared = starlark.StringDict{
	"Bool":   starlark.NewBuiltin("Bool", BoolProvider),
	"Int":    starlark.NewBuiltin("Int", IntProvider),
	"Float":  starlark.NewBuiltin("Float", FloatProvider),
	"String": starlark.NewBuiltin("String", StringProvider),
	"Object": starlark.NewBuiltin("Object", ObjectProvider),
	"List":   starlark.NewBuiltin("List", ListProvider),
	"Schema": starlark.NewBuiltin("Schema", SchemaProvider),
}

type Descriptor interface {
	SKU() string
	Default() starlark.Value
	IsRequired() starlark.Bool
	Evaluate(thread *starlark.Thread, value starlark.Value) (starlark.Value, error)
	// Conform to starlark.Value
	String() string
	Type() string
	Freeze()
	Truth() starlark.Bool
	Hash() (uint32, error)
}

func jsonify(value starlark.Value) string {
	str, err := json.Marshal(&struct {
		Type       string
		Descriptor starlark.Value
	}{
		Type:       value.Type(),
		Descriptor: value,
	})
	if err != nil {
		logrus.Panic(err)
	}
	return string(str)
}

func hashify(value starlark.Value) (uint32, error) {
	hash := fnv.New32a()
	hash.Write([]byte(value.String()))
	return hash.Sum32(), nil
}

func runValidations(
	thread *starlark.Thread,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
	validations []starlark.Callable) error {

	for _, validation := range validations {
		result, err := validation.CallInternal(thread, args, kwargs)
		if err != nil {
			return err
		} else if result.Type() != starlark.None.Type() {
			return errors.New(result.String())
		}
	}

	return nil
}

func extractValidations(validations *[]starlark.Callable, rawInputValue starlark.Value) error {
	validationsValue, ok := rawInputValue.(*starlark.List)
	if !ok {
		return fmt.Errorf(
			"Expected validations value to be a list of functions, but got %s.", rawInputValue)
	}

	for i := 0; i < validationsValue.Len(); i++ {
		item := validationsValue.Index(i)
		itemFunc, ok := item.(starlark.Callable)
		if !ok {
			return fmt.Errorf("Expected validation to be a functions, but got %s.", item)
		}
		*validations = append(*validations, itemFunc)
	}

	return nil
}
