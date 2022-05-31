package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - BoolProvider

func BoolProvider(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	provider := BoolDescriptor{
		UUID:         uuid.New(),
		DefaultValue: false,
		Required:     false,
		Validations:  []starlark.Callable{},
	}

	if args.Len() > 0 {
		return starlark.None, fmt.Errorf("Invalid positional arguments %s in Bool().", args)
	}

	for kwargName, kwargValue := range util.KwargsToMap(kwargs) {
		switch kwargName {
		case "default":
			defaultValue, ok := kwargValue.(starlark.Bool)
			if !ok {
				return starlark.None, fmt.Errorf(
					"Expected default value to be bool, but got %s.", kwargValue)
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
			return starlark.None, fmt.Errorf("Unknown keyword %s in Bool().", kwargName)
		}
	}

	return provider, nil
}

// MARK: - BoolDescriptor

type BoolDescriptor struct {
	UUID         uuid.UUID
	DefaultValue starlark.Bool
	Required     starlark.Bool
	Validations  []starlark.Callable
}

func (descriptor BoolDescriptor) SKU() string {
	return fmt.Sprintf("starfig::descriptor:bool:%s", descriptor.UUID)
}

func (descriptor BoolDescriptor) Default() starlark.Value {
	return descriptor.DefaultValue
}

func (descriptor BoolDescriptor) IsRequired() starlark.Bool {
	return descriptor.Required
}

func (descriptor BoolDescriptor) Evaluate(
	thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	boolValue, ok := value.(starlark.Bool)
	if !ok {
		return starlark.None, fmt.Errorf("Expected bool type but got %s.", value)
	}

	args := starlark.Tuple{boolValue}
	kwargs := []starlark.Tuple{}
	err := runValidations(thread, args, kwargs, descriptor.Validations)

	return boolValue, err
}

func (descriptor BoolDescriptor) String() string {
	return jsonify(descriptor)
}

func (descriptor BoolDescriptor) Type() string {
	return "BoolDescriptor"
}

func (descriptor BoolDescriptor) Freeze() {}

func (descriptor BoolDescriptor) Truth() starlark.Bool {
	return true
}

func (descriptor BoolDescriptor) Hash() (uint32, error) {
	return hashify(descriptor)
}
