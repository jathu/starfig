package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - ListProvider

func ListProvider(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	provider := ListDescriptor{
		UUID:        uuid.New(),
		Validations: []starlark.Callable{},
	}

	if args.Len() == 0 {
		return starlark.None, fmt.Errorf(
			"List requires a type. i.e. List(String), List(Foo).")
	} else if args.Len() > 1 {
		return starlark.None, fmt.Errorf(
			"List can only have one type. i.e. List(String), List(Foo).")
	}

	schemaBuilderFunction, ok := args[0].(*starlark.Builtin)
	if ok {
		switch schemaBuilderFunction.Name() {
		case "Bool":
			provider.WrappedDescriptor = BoolDescriptor{}
		case "Float":
			provider.WrappedDescriptor = FloatDescriptor{}
		case "Int":
			provider.WrappedDescriptor = IntDescriptor{}
		case "String":
			provider.WrappedDescriptor = StringDescriptor{}
		default:
			contextManager := thread.Local(SchemaContextManagerThreadKey).(SchemaContextManager)
			descriptor, found := contextManager.GetDescriptor(schemaBuilderFunction.Name())
			if found {
				provider.WrappedDescriptor = descriptor
			} else {
				return provider, fmt.Errorf("Unable to find %s.", schemaBuilderFunction.Name())
			}
		}
	} else {
		return starlark.None, fmt.Errorf("Invalid list object %s.", args[0])
	}

	for kwargName, kwargValue := range util.KwargsToMap(kwargs) {
		switch kwargName {
		case "validations":
			err := extractValidations(&provider.Validations, kwargValue)
			if err != nil {
				return starlark.None, err
			}
		default:
			return starlark.None, fmt.Errorf("Unknown keyword %s in List().", kwargName)
		}
	}

	return provider, nil
}

// MARK: - ListDescriptor

type ListDescriptor struct {
	UUID              uuid.UUID
	WrappedDescriptor Descriptor
	Validations       []starlark.Callable
}

func (descriptor ListDescriptor) SKU() string {
	return fmt.Sprintf("starfig::descriptor:list:%s", descriptor.UUID)
}

func (descriptor ListDescriptor) Default() starlark.Value {
	return starlark.NewList([]starlark.Value{})
}

func (descriptor ListDescriptor) IsRequired() starlark.Bool {
	return false
}

func (descriptor ListDescriptor) Evaluate(
	thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	listValue, ok := value.(*starlark.List)
	if !ok {
		return starlark.None, fmt.Errorf("Expected list type but got %s.", value)
	}

	evaluatedValues := []starlark.Value{}
	for i := 0; i < listValue.Len(); i++ {
		evaluatedValue, err := descriptor.WrappedDescriptor.Evaluate(thread, listValue.Index(i))
		if err != nil {
			return starlark.None, err
		}
		evaluatedValues = append(evaluatedValues, evaluatedValue)
	}
	evaluatedValuesList := starlark.NewList(evaluatedValues)
	args := starlark.Tuple{evaluatedValuesList}
	kwargs := []starlark.Tuple{}
	err := runValidations(thread, args, kwargs, descriptor.Validations)

	return evaluatedValuesList, err
}

func (descriptor ListDescriptor) String() string {
	return jsonify(descriptor)
}

func (descriptor ListDescriptor) Type() string {
	return "ListDescriptor"
}

func (descriptor ListDescriptor) Freeze() {
	// no-op for now
}

func (descriptor ListDescriptor) Truth() starlark.Bool {
	return true
}

func (descriptor ListDescriptor) Hash() (uint32, error) {
	return hashify(descriptor)
}
