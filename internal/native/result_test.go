package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/target"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

// MARK: - SchemaResult

func TestSchemaResultEvaluate(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	fields := new(starlark.Dict)
	fields.SetKey(starlark.String("ovo"), StringDescriptor{Required: true})
	descriptor := SchemaDescriptor{UUID: uuid.New(), Fields: fields}
	manager.QueueSeenDescriptor(descriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)
	schemaResult := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: descriptor,
		Evaluated:        new(starlark.Dict),
	}
	err := schemaResult.Evaluate(
		&thread,
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("ovo"), starlark.String("yeezy")},
		},
	)

	assert.Nil(t, err)
	expected := new(starlark.Dict)
	expected.SetKey(starlark.String("ovo"), starlark.String("yeezy"))
	assert.Equal(t, expected, schemaResult.Evaluated)
}

func TestSchemaResultEvaluateUnknownSchemaName(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	schemaResult := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: descriptor,
		Evaluated:        new(starlark.Dict),
	}
	err := schemaResult.Evaluate(&thread, starlark.Tuple{}, []starlark.Tuple{})

	expected := fmt.Sprintf("Unable to find %s. You might be instantiating a schema in a .star file, which is invalid.", descriptor.SKU())
	assert.ErrorContains(t, err, expected)
}

func TestSchemaResultEvaluateWithArguments(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(descriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)
	schemaResult := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: descriptor,
		Evaluated:        new(starlark.Dict),
	}
	err := schemaResult.Evaluate(
		&thread,
		starlark.Tuple{starlark.MakeInt(416)},
		[]starlark.Tuple{},
	)

	expected := fmt.Sprintf("Invalid positional arguments (416,) in Supreme.")
	assert.ErrorContains(t, err, expected)
}

func TestSchemaResultEvaluateMissingRequiredField(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	fields := new(starlark.Dict)
	fields.SetKey(starlark.String("ovo"), StringDescriptor{Required: true})
	descriptor := SchemaDescriptor{UUID: uuid.New(), Fields: fields}
	manager.QueueSeenDescriptor(descriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)
	schemaResult := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: descriptor,
		Evaluated:        new(starlark.Dict),
	}
	err := schemaResult.Evaluate(
		&thread,
		starlark.Tuple{},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, "Missing required field ovo in Supreme.")
}

func TestSchemaResultEvaluateUnknownKeyword(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	fields := new(starlark.Dict)
	descriptor := SchemaDescriptor{UUID: uuid.New(), Fields: fields}
	manager.QueueSeenDescriptor(descriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)
	schemaResult := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: descriptor,
		Evaluated:        new(starlark.Dict),
	}
	err := schemaResult.Evaluate(
		&thread,
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("mock"), starlark.MakeInt(416)},
		},
	)

	expected := fmt.Sprintf("Unknown keyword mock in Supreme.")
	assert.ErrorContains(t, err, expected)
}

func TestSchemaResultEvaluateInvalidField(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	fields := new(starlark.Dict)
	fields.SetKey(starlark.String("ovo"), StringDescriptor{Required: true})
	descriptor := SchemaDescriptor{UUID: uuid.New(), Fields: fields}
	manager.QueueSeenDescriptor(descriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)
	schemaResult := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: descriptor,
		Evaluated:        new(starlark.Dict),
	}
	err := schemaResult.Evaluate(
		&thread,
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("ovo"), starlark.MakeInt(416)},
		},
	)

	expected := fmt.Sprintf("Invalid field ovo in Supreme: Expected string type but got 416.")
	assert.ErrorContains(t, err, expected)
}

func TestSchemaResultString(t *testing.T) {
	id := uuid.New()
	childId := uuid.New()
	descriptor := SchemaResult{
		UUID:             id,
		SchemaDescriptor: SchemaDescriptor{UUID: childId},
		Evaluated:        new(starlark.Dict),
	}
	expected := fmt.Sprintf(`{"Type":"SchemaResult","Descriptor":{"UUID":"%s","SchemaDescriptor":{"UUID":"%s","Fields":null,"Validations":null},"Evaluated":{}}}`, id, childId)

	assert.Equal(t, expected, descriptor.String())
}

func TestSchemaResultType(t *testing.T) {
	descriptor := SchemaResult{}

	assert.Equal(t, "SchemaResult", descriptor.Type())
}

func TestSchemaResultFreeze(t *testing.T) {
	descriptor := SchemaResult{}
	descriptor.Freeze() // no-op
}

func TestSchemaResultTruth(t *testing.T) {
	descriptor := SchemaResult{}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestSchemaResultHash(t *testing.T) {
	hash, err := SchemaResult{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(2090021955), hash)
}
