package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

// MARK: - ListProvider

func TestListProviderWithoutType(t *testing.T) {
	_, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{},
	)
	assert.ErrorContains(t, err, "List requires a type. i.e. List(String), List(Foo).")
}

func TestListProviderWithMultipleArguments(t *testing.T) {
	_, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.None, starlark.None},
		[]starlark.Tuple{},
	)
	assert.ErrorContains(t, err, "List can only have one type. i.e. List(String), List(Foo).")
}

func TestListProviderWithBoolArgument(t *testing.T) {
	provider, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.NewBuiltin("Bool", BoolProvider)},
		[]starlark.Tuple{},
	)
	assert.Nil(t, err)
	descriptor := provider.(ListDescriptor)
	assert.Equal(t, "BoolDescriptor", descriptor.WrappedDescriptor.Type())
}

func TestListProviderWithFloatArgument(t *testing.T) {
	provider, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.NewBuiltin("Float", FloatProvider)},
		[]starlark.Tuple{},
	)
	assert.Nil(t, err)
	descriptor := provider.(ListDescriptor)
	assert.Equal(t, "FloatDescriptor", descriptor.WrappedDescriptor.Type())
}

func TestListProviderWithIntArgument(t *testing.T) {
	provider, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.NewBuiltin("Int", IntProvider)},
		[]starlark.Tuple{},
	)
	assert.Nil(t, err)
	descriptor := provider.(ListDescriptor)
	assert.Equal(t, "IntDescriptor", descriptor.WrappedDescriptor.Type())
}

func TestListProviderWithStringArgument(t *testing.T) {
	provider, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.NewBuiltin("String", StringProvider)},
		[]starlark.Tuple{},
	)
	assert.Nil(t, err)
	descriptor := provider.(ListDescriptor)
	assert.Equal(t, "StringDescriptor", descriptor.WrappedDescriptor.Type())
}

func TestListProviderWithSchemaArgument(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	providerResult, err := SchemaProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("fields"), new(starlark.Dict)},
		},
	)
	assert.Nil(t, err)
	builder := providerResult.(*starlark.Builtin)
	provider, err := ListProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{builder},
		[]starlark.Tuple{},
	)
	descriptor := provider.(ListDescriptor)
	assert.Equal(t, "SchemaDescriptor", descriptor.WrappedDescriptor.Type())
}

func TestListProviderWithMissingSchemaArgument(t *testing.T) {
	manager := NewSchemaContextManager()
	thread := starlark.Thread{}
	thread.SetLocal(SchemaContextManagerThreadKey, manager)
	_, err := ListProvider(
		&thread,
		tester.MockBuiltin(),
		starlark.Tuple{tester.MockBuiltinWithName("unknown-func")},
		[]starlark.Tuple{},
	)
	assert.ErrorContains(t, err, "Unable to find unknown-func.")
}

func TestListProviderWithNonFunctionArgument(t *testing.T) {
	_, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.None},
		[]starlark.Tuple{},
	)
	assert.ErrorContains(t, err, "Invalid list object None.")
}

func TestListProviderWithInvalidValidationsType(t *testing.T) {
	_, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.NewBuiltin("String", StringProvider)},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.MakeInt(416)},
		},
	)
	assert.ErrorContains(t, err,
		`Expected validations value to be a list of functions, but got 416.`)
}

func TestListProviderWithInvalidValidationsElementType(t *testing.T) {
	_, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.NewBuiltin("String", StringProvider)},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.NewList([]starlark.Value{starlark.MakeInt(416)})},
		},
	)
	assert.ErrorContains(t, err,
		`Expected validation to be a functions, but got 416.`)
}

func TestListProviderWithUnknownKeyword(t *testing.T) {
	_, err := ListProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.NewBuiltin("String", StringProvider)},
		[]starlark.Tuple{
			{starlark.String("supreme"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Unknown keyword supreme in List().`)
}

// MARK: - ListDescriptor

func TestListDescriptorSKU(t *testing.T) {
	id := uuid.New()
	descriptor := ListDescriptor{UUID: id}

	assert.Equal(t, fmt.Sprintf("starfig::descriptor:list:%s", id), descriptor.SKU())
}

func TestListDescriptorDefault(t *testing.T) {
	descriptor := ListDescriptor{WrappedDescriptor: IntDescriptor{
		DefaultValue: starlark.MakeInt(416),
	}}

	assert.Equal(t, starlark.NewList([]starlark.Value{}), descriptor.Default())
}

func TestListDescriptorIsRequired(t *testing.T) {
	descriptor := ListDescriptor{}

	assert.Equal(t, starlark.Bool(false), descriptor.IsRequired())
}

func TestListDescriptorEvaluate(t *testing.T) {
	userValues := starlark.NewList([]starlark.Value{
		starlark.String("mock1"),
		starlark.String("mock2"),
		starlark.String("mock3"),
	})
	descriptor := ListDescriptor{
		WrappedDescriptor: StringDescriptor{},
		Validations: []starlark.Callable{
			tester.MockBuiltinWithCallback(func(args starlark.Tuple, kwargs []starlark.Tuple) {
				assert.Equal(t, starlark.Tuple{userValues}, args)
				assert.ElementsMatch(t, []starlark.Tuple{}, kwargs)
			}),
		},
	}
	evaluatedValue, err := descriptor.Evaluate(&starlark.Thread{}, userValues)

	assert.Nil(t, err)
	assert.Equal(t, userValues, evaluatedValue)
}

func TestListDescriptorEvaluateInvalidType(t *testing.T) {
	descriptor := ListDescriptor{}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.String("mock"))

	assert.ErrorContains(t, err, `Expected list type but got "mock".`)
}

func TestListDescriptorEvaluateWrappedEvaluateError(t *testing.T) {
	descriptor := ListDescriptor{
		WrappedDescriptor: BoolDescriptor{},
	}

	userValues := starlark.NewList([]starlark.Value{starlark.String("mock")})
	_, err := descriptor.Evaluate(&starlark.Thread{}, userValues)
	assert.ErrorContains(t, err, `Expected bool type but got "mock".`)
}

func TestListDescriptorEvaluateValidationError(t *testing.T) {
	descriptor := ListDescriptor{
		WrappedDescriptor: StringDescriptor{},
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingFunction("yikes!"),
		},
	}
	userValues := starlark.NewList([]starlark.Value{starlark.String("mock")})
	_, err := descriptor.Evaluate(&starlark.Thread{}, userValues)

	assert.ErrorContains(t, err, "yikes!")
}

func TestListDescriptorEvaluateValidationUserError(t *testing.T) {
	descriptor := ListDescriptor{
		WrappedDescriptor: StringDescriptor{},
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingBuiltin("yikes!"),
		},
	}
	userValues := starlark.NewList([]starlark.Value{starlark.String("mock")})
	_, err := descriptor.Evaluate(&starlark.Thread{}, userValues)

	assert.ErrorContains(t, err, "yikes!")
}

func TestListDescriptorString(t *testing.T) {
	id := uuid.New()
	childId := uuid.New()
	descriptor := ListDescriptor{
		UUID:              id,
		WrappedDescriptor: IntDescriptor{UUID: childId},
	}
	expected := fmt.Sprintf(`{"Type":"ListDescriptor","Descriptor":{"UUID":"%s","WrappedDescriptor":{"UUID":"%s","DefaultValue":{},"Required":false,"Validations":null},"Validations":null}}`, id, childId)

	assert.Equal(t, expected, descriptor.String())
}

func TestListDescriptorType(t *testing.T) {
	descriptor := ListDescriptor{}

	assert.Equal(t, "ListDescriptor", descriptor.Type())
}

func TestListDescriptorFreeze(t *testing.T) {
	descriptor := ListDescriptor{}
	descriptor.Freeze() // no-op
}

func TestListDescriptorTruth(t *testing.T) {
	descriptor := ListDescriptor{}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestListDescriptorHash(t *testing.T) {
	hash, err := ListDescriptor{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(971521322), hash)
}
