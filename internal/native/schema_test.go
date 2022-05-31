package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/target"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
	"golang.org/x/exp/maps"
)

// MARK: - SchemaProvider

func TestSchemaProvider(t *testing.T) {
	thread := starlark.Thread{}
	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	fields := new(starlark.Dict)
	fields.SetKey(starlark.String("example_1"), StringDescriptor{})
	fields.SetKey(starlark.String("example_2"), BoolDescriptor{})
	fields.SetKey(starlark.String("example_3"), IntDescriptor{})
	validations := starlark.NewList([]starlark.Value{
		tester.MockBuiltin(),
		tester.MockBuiltin(),
		tester.MockBuiltin(),
	})
	providerResult, err := SchemaProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("fields"), fields},
			{starlark.String("validations"), validations},
		},
	)

	assert.Nil(t, err)
	builder := providerResult.(*starlark.Builtin)
	queuedDescriptor := manager.queue[maps.Keys(manager.queue)[0]]
	assert.Equal(t, queuedDescriptor.SKU(), builder.Name())
	assert.Equal(t, fields, queuedDescriptor.Fields)
	tester.AssertSameValidations(t, validations, queuedDescriptor.Validations)
}

func TestSchemaProviderWithArguments(t *testing.T) {
	_, err := SchemaProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.String("ok")},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, `Invalid positional arguments ("ok",) in Schema().`)
}

func TestSchemaProviderWithInvalidFieldsType(t *testing.T) {
	_, err := SchemaProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("fields"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Expected fields to be a dict, but got 416.`)
}

func TestSchemaProviderWithInvalidValidationsType(t *testing.T) {
	_, err := SchemaProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err,
		`Expected validations value to be a list of functions, but got 416.`)
}

func TestSchemaProviderWithInvalidValidationsElementType(t *testing.T) {
	_, err := SchemaProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.NewList([]starlark.Value{starlark.MakeInt(416)})},
		},
	)

	assert.ErrorContains(t, err, `Expected validation to be a functions, but got 416.`)
}

func TestSchemaProviderWithUnknownKeyword(t *testing.T) {
	_, err := SchemaProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("supreme"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Unknown keyword supreme in Schema().`)
}

func TestSchemaProviderWithInvalidFieldsValueType(t *testing.T) {
	fields := new(starlark.Dict)
	fields.SetKey(starlark.String("example"), starlark.MakeInt(416))
	_, err := SchemaProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("fields"), fields},
		},
	)

	assert.ErrorContains(t, err, `Expected a descriptor for "example", but got 416.`)
}

// MARK: - createSchemaBuilder

func TestCreateSchemaBuilder(t *testing.T) {
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	result, err := createSchemaBuilder(descriptor)

	assert.Nil(t, err)
	assert.Equal(t, descriptor.SKU(), result.Name())
}

// MARK: - SchemaDescriptor

func TestSchemaDescriptorSKU(t *testing.T) {
	id := uuid.New()
	descriptor := SchemaDescriptor{UUID: id}

	assert.Equal(t, fmt.Sprintf("starfig::descriptor:schema:%s", id), descriptor.SKU())
}

func TestSchemaDescriptorDefault(t *testing.T) {
	fields := new(starlark.Dict)
	fields.SetKey(starlark.String("ex-bool"), BoolDescriptor{})
	fields.SetKey(starlark.String("ex-string"), StringDescriptor{DefaultValue: "hello"})
	descriptor := SchemaDescriptor{Fields: fields}

	expected := new(starlark.Dict)
	expected.SetKey(starlark.String("ex-bool"), starlark.Bool(false))
	expected.SetKey(starlark.String("ex-string"), starlark.String("hello"))
	assert.Equal(t, expected, descriptor.Default())
}

func TestSchemaDescriptorIsRequired(t *testing.T) {
	descriptor := SchemaDescriptor{UUID: uuid.New()}

	assert.Equal(t, starlark.Bool(false), descriptor.IsRequired())
}

func TestSchmeaDescriptorEvaluate(t *testing.T) {
	thread := starlark.Thread{}

	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)

	expectedEvaluated := new(starlark.Dict)

	descriptor := SchemaDescriptor{
		UUID: uuid.New(),
		Validations: []starlark.Callable{
			tester.MockBuiltinWithCallback(func(args starlark.Tuple, kwargs []starlark.Tuple) {
				assert.Equal(t, starlark.Tuple{expectedEvaluated}, args)
				assert.ElementsMatch(t, []starlark.Tuple{}, kwargs)
			}),
		},
	}
	manager.QueueSeenDescriptor(descriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)

	userValue := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: descriptor,
		Evaluated:        expectedEvaluated,
	}
	result, err := descriptor.Evaluate(&thread, userValue)
	assert.Nil(t, err)
	assert.Equal(t, expectedEvaluated, result)
}

func TestSchmeaDescriptorEvaluateUnknownSchema(t *testing.T) {
	thread := starlark.Thread{}

	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)

	descriptor := SchemaDescriptor{UUID: uuid.New()}
	_, err := descriptor.Evaluate(&thread, starlark.None)

	expected := fmt.Sprintf("Unable to find %s in schema evaluation.", descriptor.SKU())
	assert.ErrorContains(t, err, expected)
}

func TestSchmeaDescriptorEvaluateIncorrectPrimitiveType(t *testing.T) {
	thread := starlark.Thread{}

	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)

	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(descriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)
	_, err := descriptor.Evaluate(&thread, starlark.MakeInt(416))

	assert.ErrorContains(t, err, "Expected Supreme type but got 416.")
}

func TestSchmeaDescriptorEvaluateUnknownSchemaName(t *testing.T) {
	thread := starlark.Thread{}

	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)

	wantedDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(wantedDescriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(wantedDescriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)

	otherDescriptor := SchemaDescriptor{UUID: uuid.New()}
	userValue := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: otherDescriptor,
		Evaluated:        new(starlark.Dict),
	}
	_, err := wantedDescriptor.Evaluate(&thread, userValue)

	expected := fmt.Sprintf("Unable to find %s in schema evaluation.", otherDescriptor.SKU())
	assert.ErrorContains(t, err, expected)
}

func TestSchmeaDescriptorEvaluateIncorrectSchemaType(t *testing.T) {
	thread := starlark.Thread{}

	manager := NewSchemaContextManager()
	thread.SetLocal(SchemaContextManagerThreadKey, manager)

	wantedDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(wantedDescriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(wantedDescriptor.SKU()),
		"Supreme",
		target.FileTarget{},
	)

	otherDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(otherDescriptor)
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(otherDescriptor.SKU()),
		"Patagonia",
		target.FileTarget{},
	)

	userValue := SchemaResult{
		UUID:             uuid.New(),
		SchemaDescriptor: otherDescriptor,
		Evaluated:        new(starlark.Dict),
	}
	_, err := wantedDescriptor.Evaluate(&thread, userValue)

	assert.ErrorContains(t, err, "Expected Supreme but got Patagonia.")
}

func TestSchemaDescriptorString(t *testing.T) {
	id := uuid.New()
	fields := new(starlark.Dict)
	fields.SetKey(starlark.String("ex-bool"), BoolDescriptor{})
	fields.SetKey(starlark.String("ex-string"), StringDescriptor{DefaultValue: "hello"})
	descriptor := SchemaDescriptor{
		UUID:   id,
		Fields: fields,
	}

	expected := fmt.Sprintf(`{"Type":"SchemaDescriptor","Descriptor":{"UUID":"%s","Fields":{},"Validations":null}}`, id)

	assert.Equal(t, expected, descriptor.String())
}

func TestSchemaDescriptorType(t *testing.T) {
	descriptor := SchemaDescriptor{UUID: uuid.New()}

	assert.Equal(t, "SchemaDescriptor", descriptor.Type())
}

func TestSchemaDescriptorFreeze(t *testing.T) {
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	descriptor.Freeze() // no-op
}

func TestSchemaDescriptorTruth(t *testing.T) {
	descriptor := SchemaDescriptor{UUID: uuid.New()}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestSchemaDescriptorHash(t *testing.T) {
	hash, err := SchemaDescriptor{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(1590550606), hash)
}
