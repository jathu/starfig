package native

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/util"
	"go.starlark.net/starlark"
)

// MARK: - SchemaResult

type SchemaResult struct {
	UUID             uuid.UUID
	SchemaDescriptor SchemaDescriptor
	Evaluated        *starlark.Dict
}

func (result SchemaResult) Evaluate(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) error {
	contextManager := thread.Local(SchemaContextManagerThreadKey).(SchemaContextManager)
	schemaName, found := contextManager.GetSchemaName(result.SchemaDescriptor)
	if !found {
		// This should ideally never happen because starlark-go ensures
		// all the symbols are loaded and we populate the manager during
		// schema creation and file loading.
		return fmt.Errorf("Unable to find %s. You might be instantiating a schema in a .star file, which is invalid.", result.SchemaDescriptor.SKU())
	}

	if args.Len() > 0 {
		return fmt.Errorf("Invalid positional arguments %s in %s.", args, schemaName)
	}

	kwargMap := util.KwargsToMap(kwargs)

	for _, tuple := range result.SchemaDescriptor.Fields.Items() {
		fieldDescriptor := tuple.Index(1).(Descriptor)
		if fieldDescriptor.IsRequired() {
			fieldName := tuple.Index(0).(starlark.String).GoString()
			_, found := kwargMap[fieldName]
			if !found {
				return fmt.Errorf("Missing required field %s in %s.", fieldName, schemaName)
			}
		}
	}

	for name, value := range kwargMap {
		fieldDescriptorValue, found, err := result.SchemaDescriptor.Fields.Get(starlark.String(name))
		if err != nil || !found {
			return fmt.Errorf("Unknown keyword %s in %s.", name, schemaName)
		}
		fieldDescriptor := fieldDescriptorValue.(Descriptor)
		evaluatedValue, err := fieldDescriptor.Evaluate(thread, value)
		if err != nil {
			return fmt.Errorf("Invalid field %s in %s: %s", name, schemaName, err)
		}
		result.Evaluated.SetKey(starlark.String(name), evaluatedValue)
	}

	return nil
}

func (result SchemaResult) String() string {
	return jsonify(result)
}

func (result SchemaResult) Type() string {
	return "SchemaResult"
}

func (result SchemaResult) Freeze() {
	// no-op for now
}

func (result SchemaResult) Truth() starlark.Bool {
	return true
}

func (result SchemaResult) Hash() (uint32, error) {
	return hashify(result)
}
