package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - ObjectProvider

func ObjectProvider(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	provider := ObjectDescriptor{
		UUID:        uuid.New(),
		Required:    false,
		Validations: []starlark.Callable{},
	}

	if args.Len() == 0 {
		return starlark.None, fmt.Errorf("Object requires a schema type. i.e. Object(Foo).")
	} else if args.Len() > 1 {
		return starlark.None, fmt.Errorf("Object can only have one type. i.e. Object(Foo).")
	}

	schemaBuilderFunction, ok := args[0].(*starlark.Builtin)
	if !ok {
		return starlark.None, fmt.Errorf("Object can only be another schema, not %s.", args[0])
	} else {
		contextManager := thread.Local(SchemaContextManagerThreadKey).(SchemaContextManager)
		descriptor, found := contextManager.GetDescriptor(schemaBuilderFunction.Name())
		if !found {
			return provider, fmt.Errorf("Unable to find schema %s.", schemaBuilderFunction.Name())
		}
		provider.WrappedDescriptor = descriptor
	}

	for kwargName, kwargValue := range util.KwargsToMap(kwargs) {
		switch kwargName {
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
			return starlark.None, fmt.Errorf("Unknown keyword %s in Object().", kwargName)
		}
	}

	return provider, nil
}

// MARK: - ObjectDescriptor

type ObjectDescriptor struct {
	UUID              uuid.UUID
	WrappedDescriptor Descriptor
	Required          starlark.Bool
	Validations       []starlark.Callable
}

func (descriptor ObjectDescriptor) SKU() string {
	return fmt.Sprintf("starfig::descriptor:object:%s", descriptor.UUID)
}

func (descriptor ObjectDescriptor) Default() starlark.Value {
	return descriptor.WrappedDescriptor.Default()
}

func (descriptor ObjectDescriptor) IsRequired() starlark.Bool {
	return descriptor.Required
}

func (descriptor ObjectDescriptor) Evaluate(
	thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	evaluatedValue, err := descriptor.WrappedDescriptor.Evaluate(thread, value)
	if err != nil {
		return starlark.None, err
	}

	args := starlark.Tuple{evaluatedValue}
	kwargs := []starlark.Tuple{}
	err = runValidations(thread, args, kwargs, descriptor.Validations)

	return evaluatedValue, err
}

func (descriptor ObjectDescriptor) String() string {
	return jsonify(descriptor)
}

func (descriptor ObjectDescriptor) Type() string {
	return "ObjectDescriptor"
}

func (descriptor ObjectDescriptor) Freeze() {
	// no-op for now
}

func (descriptor ObjectDescriptor) Truth() starlark.Bool {
	return true
}

func (descriptor ObjectDescriptor) Hash() (uint32, error) {
	return hashify(descriptor)
}
