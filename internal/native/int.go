package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - IntProvider

func IntProvider(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	provider := IntDescriptor{
		UUID:         uuid.New(),
		DefaultValue: starlark.MakeInt(0),
		Required:     false,
		Validations:  []starlark.Callable{},
	}

	if args.Len() > 0 {
		return starlark.None, fmt.Errorf("Invalid positional arguments %s in Int().", args)
	}

	for kwargName, kwargValue := range util.KwargsToMap(kwargs) {
		switch kwargName {
		case "default":
			defaultValue, ok := kwargValue.(starlark.Int)
			if !ok {
				return starlark.None, fmt.Errorf(
					"Expected default value to be int, but got %s.", kwargValue)
			}
			provider.DefaultValue = defaultValue
		case "required":
			requiredValue, ok := kwargValue.(starlark.Bool)
			if !ok {
				return starlark.None, fmt.Errorf(
					"Expected required value to be bool, but got %s.", kwargValue)
			}
			provider.Required = requiredValue
		case "validations":
			err := extractValidations(&provider.Validations, kwargValue)
			if err != nil {
				return starlark.None, err
			}
		default:
			return starlark.None, fmt.Errorf("Unknown keyword %s in Int().", kwargName)
		}
	}

	return provider, nil
}

// MARK: - IntDescriptor

type IntDescriptor struct {
	UUID         uuid.UUID
	DefaultValue starlark.Int
	Required     starlark.Bool
	Validations  []starlark.Callable
}

func (descriptor IntDescriptor) SKU() string {
	return fmt.Sprintf("starfig::descriptor:int:%s", descriptor.UUID)
}

func (descriptor IntDescriptor) Default() starlark.Value {
	return descriptor.DefaultValue
}

func (descriptor IntDescriptor) IsRequired() starlark.Bool {
	return descriptor.Required
}

func (descriptor IntDescriptor) Evaluate(
	thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	intValue, ok := value.(starlark.Int)
	if !ok {
		return starlark.None, fmt.Errorf("Expected int type but got %s.", value)
	}

	args := starlark.Tuple{intValue}
	kwargs := []starlark.Tuple{}
	err := runValidations(thread, args, kwargs, descriptor.Validations)

	return intValue, err
}

func (descriptor IntDescriptor) String() string {
	return jsonify(descriptor)
}

func (descriptor IntDescriptor) Type() string {
	return "IntDescriptor"
}

func (descriptor IntDescriptor) Freeze() {
	// no-op for now
}

func (descriptor IntDescriptor) Truth() starlark.Bool {
	return true
}

func (descriptor IntDescriptor) Hash() (uint32, error) {
	return hashify(descriptor)
}
