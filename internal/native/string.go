package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - StringProvider

func StringProvider(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	provider := StringDescriptor{
		UUID:         uuid.New(),
		DefaultValue: "",
		Required:     false,
		Validations:  []starlark.Callable{},
	}

	if args.Len() > 0 {
		return starlark.None, fmt.Errorf("Invalid positional arguments %s in String().", args)
	}

	for kwargName, kwargValue := range util.KwargsToMap(kwargs) {
		switch kwargName {
		case "default":
			defaultValue, ok := kwargValue.(starlark.String)
			if !ok {
				return starlark.None, fmt.Errorf(
					"Expected default value to be string, but got %s.", kwargValue)
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
			return starlark.None, fmt.Errorf("Unknown keyword %s in String().", kwargName)
		}
	}

	return provider, nil
}

// MARL: - StringDescriptor

type StringDescriptor struct {
	UUID         uuid.UUID
	DefaultValue starlark.String
	Required     starlark.Bool
	Validations  []starlark.Callable
}

func (descriptor StringDescriptor) SKU() string {
	return fmt.Sprintf("starfig::descriptor:string:%s", descriptor.UUID)
}

func (descriptor StringDescriptor) Default() starlark.Value {
	return descriptor.DefaultValue
}

func (descriptor StringDescriptor) IsRequired() starlark.Bool {
	return descriptor.Required
}

func (descriptor StringDescriptor) Evaluate(
	thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	stringValue, ok := value.(starlark.String)
	if !ok {
		return starlark.None, fmt.Errorf("Expected string type but got %s.", value)
	}

	args := starlark.Tuple{stringValue}
	kwargs := []starlark.Tuple{}
	err := runValidations(thread, args, kwargs, descriptor.Validations)

	return stringValue, err
}

func (descriptor StringDescriptor) String() string {
	return jsonify(descriptor)
}

func (descriptor StringDescriptor) Type() string {
	return "StringDescriptor"
}

func (descriptor StringDescriptor) Freeze() {
	// no-op for now
}

func (descriptor StringDescriptor) Truth() starlark.Bool {
	return true
}

func (descriptor StringDescriptor) Hash() (uint32, error) {
	return hashify(descriptor)
}
