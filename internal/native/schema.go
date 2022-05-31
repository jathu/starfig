package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - SchemaProvider

func SchemaProvider(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	provider := SchemaDescriptor{
		UUID:        uuid.New(),
		Fields:      new(starlark.Dict),
		Validations: []starlark.Callable{},
	}

	if args.Len() > 0 {
		return starlark.None, fmt.Errorf("Invalid positional arguments %s in Schema().", args)
	}

	for kwargName, kwargValue := range util.KwargsToMap(kwargs) {
		switch kwargName {
		case "fields":
			fields, ok := kwargValue.(*starlark.Dict)
			if !ok {
				return starlark.None, fmt.Errorf(
					"Expected fields to be a dict, but got %s.", kwargValue)
			}
			for _, tuple := range fields.Items() {
				key := tuple.Index(0)
				_, ok := tuple.Index(1).(Descriptor)
				if !ok {
					return starlark.None, fmt.Errorf(
						"Expected a descriptor for %s, but got %s.", key, tuple.Index(1))
				}
			}

			provider.Fields = fields
		case "validations":
			err := extractValidations(&provider.Validations, kwargValue)
			if err != nil {
				return starlark.None, err
			}
		default:
			return starlark.None, fmt.Errorf("Unknown keyword %s in Schema().", kwargName)
		}
	}

	contextManager := thread.Local(SchemaContextManagerThreadKey).(SchemaContextManager)
	contextManager.QueueSeenDescriptor(provider)

	return createSchemaBuilder(provider)
}

func createSchemaBuilder(descriptor SchemaDescriptor) (*starlark.Builtin, error) {
	builder := func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		result := SchemaResult{
			UUID:             uuid.New(),
			SchemaDescriptor: descriptor,
			Evaluated:        descriptor.Default().(*starlark.Dict),
		}
		err := result.Evaluate(thread, args, kwargs)
		return result, err
	}

	return starlark.NewBuiltin(descriptor.SKU(), builder), nil
}

// MARK: - SchemaDescriptor

type SchemaDescriptor struct {
	UUID        uuid.UUID
	Fields      *starlark.Dict
	Validations []starlark.Callable
}

func (descriptor SchemaDescriptor) SKU() string {
	return fmt.Sprintf("starfig::descriptor:schema:%s", descriptor.UUID)
}

func (descriptor SchemaDescriptor) Default() starlark.Value {
	result := new(starlark.Dict)
	for _, tuple := range descriptor.Fields.Items() {
		key := tuple.Index(0).(starlark.String)
		value := tuple.Index(1).(Descriptor)
		result.SetKey(key, value.Default())
	}
	return result
}

func (descriptor SchemaDescriptor) IsRequired() starlark.Bool {
	return false
}

func (descriptor SchemaDescriptor) Evaluate(
	thread *starlark.Thread, value starlark.Value) (starlark.Value, error) {
	contextManager := thread.Local(SchemaContextManagerThreadKey).(SchemaContextManager)
	expectedSchemaName, ok := contextManager.GetSchemaName(descriptor)
	if !ok {
		return starlark.None, fmt.Errorf(
			"Unable to find %s in schema evaluation.", descriptor.SKU())
	}

	providedValue, ok := value.(SchemaResult)
	if !ok {
		return starlark.None, fmt.Errorf("Expected %s type but got %s.", expectedSchemaName, value)
	}

	if !contextManager.EqualDescriptor(descriptor, providedValue.SchemaDescriptor) {
		providedSchemaName, ok := contextManager.GetSchemaName(providedValue.SchemaDescriptor)
		if !ok {
			return starlark.None, fmt.Errorf(
				"Unable to find %s in schema evaluation.", providedValue.SchemaDescriptor.SKU())
		}
		return starlark.None, fmt.Errorf(
			"Expected %s but got %s.", expectedSchemaName, providedSchemaName)
	}

	args := starlark.Tuple{providedValue.Evaluated}
	kwargs := []starlark.Tuple{}
	err := runValidations(thread, args, kwargs, descriptor.Validations)

	return providedValue.Evaluated, err
}

func (descriptor SchemaDescriptor) String() string {
	return jsonify(descriptor)
}

func (descriptor SchemaDescriptor) Type() string {
	return "SchemaDescriptor"
}

func (descriptor SchemaDescriptor) Freeze() {
	// no-op for now
}

func (descriptor SchemaDescriptor) Truth() starlark.Bool {
	return true
}

func (descriptor SchemaDescriptor) Hash() (uint32, error) {
	return hashify(descriptor)
}
