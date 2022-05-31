package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

// MARK: - ObjectProvider

func TestObjectProvider(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	descriptor := SchemaDescriptor{UUID: uuid.New(), Fields: new(starlark.Dict)}
	manager.QueueSeenDescriptor(descriptor)
	validations := starlark.NewList([]starlark.Value{
		tester.MockBuiltin(),
		tester.MockBuiltin(),
		tester.MockBuiltin(),
	})
	value, err := ObjectProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{tester.MockBuiltinWithName(descriptor.SKU())},
		[]starlark.Tuple{
			{starlark.String("required"), starlark.Bool(true)},
			{starlark.String("validations"), validations},
		},
	)

	assert.Nil(t, err)
	provider := value.(ObjectDescriptor)
	assert.Equal(t, descriptor.Default(), provider.Default())
	assert.Equal(t, starlark.Bool(true), provider.IsRequired())
	tester.AssertSameValidations(t, validations, provider.Validations)
}

func TestObjectProviderNoArguments(t *testing.T) {
	_, err := ObjectProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, "Object requires a schema type. i.e. Object(Foo).")
}

func TestObjectProviderTooManyArguments(t *testing.T) {
	_, err := ObjectProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.MakeInt(416), starlark.MakeInt(905)},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, "Object can only have one type. i.e. Object(Foo).")
}

func TestObjectProviderNonSchemaWrapped(t *testing.T) {
	_, err := ObjectProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.MakeInt(416)},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, "Object can only be another schema, not 416.")
}

func TestObjectProviderUnknownSchema(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	_, err := ObjectProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{tester.MockBuiltinWithName("unknown-builtin")},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, "Unable to find schema unknown-builtin.")
}

func TestObjectProviderWithInvalidRequiredType(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(descriptor)
	_, err := ObjectProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{tester.MockBuiltinWithName(descriptor.SKU())},
		[]starlark.Tuple{
			{starlark.String("required"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Expected required value to be bool, but got 416.`)
}

func TestObjectProviderWithInvalidValidationsType(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(descriptor)
	_, err := ObjectProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{tester.MockBuiltinWithName(descriptor.SKU())},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err,
		`Expected validations value to be a list of functions, but got 416.`)
}

func TestObjectProviderWithInvalidValidationsElementType(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(descriptor)
	_, err := ObjectProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{tester.MockBuiltinWithName(descriptor.SKU())},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.NewList([]starlark.Value{starlark.MakeInt(416)})},
		},
	)

	assert.ErrorContains(t, err, `Expected validation to be a functions, but got 416.`)
}

func TestObjectProviderWithUnknownKeyword(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.QueueSeenDescriptor(descriptor)
	_, err := ObjectProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{tester.MockBuiltinWithName(descriptor.SKU())},
		[]starlark.Tuple{
			{starlark.String("supreme"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Unknown keyword supreme in Object().`)
}

// MARK: - ObjectDescriptor

func TestObjectDescriptorSKU(t *testing.T) {
	id := uuid.New()
	descriptor := ObjectDescriptor{UUID: id}

	assert.Equal(t, fmt.Sprintf("starfig::descriptor:object:%s", id), descriptor.SKU())
}

func TestObjectDescriptorDefault(t *testing.T) {
	descriptor := ObjectDescriptor{WrappedDescriptor: IntDescriptor{
		DefaultValue: starlark.MakeInt(416),
	}}

	assert.Equal(t, starlark.MakeInt(416), descriptor.Default())
}

func TestObjectDescriptorIsRequired(t *testing.T) {
	descriptor := ObjectDescriptor{Required: true}

	assert.Equal(t, starlark.Bool(true), descriptor.IsRequired())
}

func TestObjectDescriptorEvaluate(t *testing.T) {
	thread := starlark.Thread{}
	descriptor := ObjectDescriptor{
		WrappedDescriptor: IntDescriptor{
			Validations: []starlark.Callable{
				tester.MockBuiltinWithCallback(func(args starlark.Tuple, kwargs []starlark.Tuple) {
					assert.Equal(t, starlark.Tuple{starlark.MakeInt(416)}, args)
					assert.ElementsMatch(t, []starlark.Tuple{}, kwargs)
				}),
			},
		},
	}
	evaluatedValue, err := descriptor.Evaluate(&thread, starlark.MakeInt(416))

	assert.Nil(t, err)
	assert.Equal(t, starlark.MakeInt(416), evaluatedValue)
}

func TestObjectDescriptorEvaluateWrappedEvaluateError(t *testing.T) {
	thread := starlark.Thread{}
	descriptor := ObjectDescriptor{
		WrappedDescriptor: BoolDescriptor{},
	}
	_, err := descriptor.Evaluate(&thread, starlark.MakeInt(416))

	assert.ErrorContains(t, err, "Expected bool type but got 416.")
}

func TestObjectDescriptorEvaluateValidationError(t *testing.T) {
	descriptor := ObjectDescriptor{
		WrappedDescriptor: IntDescriptor{},
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingFunction("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.MakeInt(416))

	assert.ErrorContains(t, err, "yikes!")
}

func TestObjectDescriptorEvaluateValidationUserError(t *testing.T) {
	descriptor := ObjectDescriptor{
		WrappedDescriptor: IntDescriptor{},
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingBuiltin("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.MakeInt(416))

	assert.ErrorContains(t, err, "yikes!")
}

func TestObjectDescriptorString(t *testing.T) {
	id := uuid.New()
	childId := uuid.New()
	descriptor := ObjectDescriptor{
		UUID:              id,
		WrappedDescriptor: IntDescriptor{UUID: childId},
		Required:          true,
	}
	expected := fmt.Sprintf(`{"Type":"ObjectDescriptor","Descriptor":{"UUID":"%s","WrappedDescriptor":{"UUID":"%s","DefaultValue":{},"Required":false,"Validations":null},"Required":true,"Validations":null}}`, id, childId)

	assert.Equal(t, expected, descriptor.String())
}

func TestObjectDescriptorType(t *testing.T) {
	descriptor := ObjectDescriptor{}

	assert.Equal(t, "ObjectDescriptor", descriptor.Type())
}

func TestObjectDescriptorFreeze(t *testing.T) {
	descriptor := ObjectDescriptor{}
	descriptor.Freeze() // no-op
}

func TestObjectDescriptorTruth(t *testing.T) {
	descriptor := ObjectDescriptor{}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestObjectDescriptorHash(t *testing.T) {
	hash, err := ObjectDescriptor{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(2807002309), hash)
}
