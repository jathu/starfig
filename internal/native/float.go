package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - FloatProvider

func FloatProvider(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	provider := FloatDescriptor{
		UUID:         uuid.New(),
		DefaultValue: 0.0,
		Required:     false,
		Validations:  []starlark.Callable{},
	}

	if args.Len() > 0 {
		return starlark.None, fmt.Errorf("Invalid positional arguments %s in Float().", args)
	}

	for kwargName, kwargValue := range util.KwargsToMap(kwargs) {
		switch kwargName {
		case "default":
			defaultValue, ok := kwargValue.(starlark.Float)
			if !ok {
				return starlark.None, fmt.Errorf(
					"Expected default value to be float, but got %s.", kwargValue)
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
			return starlark.None, fmt.Errorf("Unknown keyword %s in Float().", kwargName)
		}
	}

	return provider, nil
}

// MARK: - FloatDescriptor

type FloatDescriptor struct {
	UUID         uuid.UUID
	DefaultValue starlark.Float
	Required     starlark.Bool
	Validations  []starlark.Callable
}

func (descriptor FloatDescriptor) SKU() string {
	return fmt.Sprintf("starfig::descriptor:float:%s", descriptor.UUID)
}

func (descriptor FloatDescriptor) Default() starlark.Value {
	return descriptor.DefaultValue
}

func (descriptor FloatDescriptor) IsRequired() starlark.Bool {
	return descriptor.Required
}

func (descriptor FloatDescriptor) Evaluate(
	thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	floatValue, ok := value.(starlark.Float)
	if !ok {
		return starlark.None, fmt.Errorf("Expected float type but got %s.", value)
	}

	args := starlark.Tuple{floatValue}
	kwargs := []starlark.Tuple{}
	err := runValidations(thread, args, kwargs, descriptor.Validations)

	return floatValue, err
}

func (descriptor FloatDescriptor) String() string {
	return jsonify(descriptor)
}

func (descriptor FloatDescriptor) Type() string {
	return "FloatDescriptor"
}

func (descriptor FloatDescriptor) Freeze() {
	// no-op for now
}

func (descriptor FloatDescriptor) Truth() starlark.Bool {
	return true
}

func (descriptor FloatDescriptor) Hash() (uint32, error) {
	return hashify(descriptor)
}
