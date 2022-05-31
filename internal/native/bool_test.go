package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

// MARK: - BoolProvider

func TestBoolProvider(t *testing.T) {
	validations := starlark.NewList([]starlark.Value{
		tester.MockBuiltin(),
		tester.MockBuiltin(),
		tester.MockBuiltin(),
	})

	value, err := BoolProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.Bool(true)},
			{starlark.String("required"), starlark.Bool(true)},
			{starlark.String("validations"), validations},
		},
	)

	assert.Nil(t, err)
	provider := value.(BoolDescriptor)
	assert.Equal(t, starlark.Bool(true), provider.Default())
	assert.Equal(t, starlark.Bool(true), provider.IsRequired())
	tester.AssertSameValidations(t, validations, provider.Validations)
}

func TestBoolProviderWithArguments(t *testing.T) {
	_, err := BoolProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.String("ok")},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, `Invalid positional arguments ("ok",) in Bool().`)
}

func TestBoolProviderWithInvalidDefaultType(t *testing.T) {
	_, err := BoolProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Expected default value to be bool, but got 416.`)
}

func TestBoolProviderWithInvalidRequiredType(t *testing.T) {
	_, err := BoolProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("required"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Expected required value to be bool, but got 416.`)
}

func TestBoolProviderWithInvalidValidationsType(t *testing.T) {
	_, err := BoolProvider(
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

func TestBoolProviderWithInvalidValidationsElementType(t *testing.T) {
	_, err := BoolProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.NewList([]starlark.Value{starlark.MakeInt(416)})},
		},
	)

	assert.ErrorContains(t, err, `Expected validation to be a functions, but got 416.`)
}

func TestBoolProviderWithUnknownKeyword(t *testing.T) {
	_, err := BoolProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("supreme"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Unknown keyword supreme in Bool().`)
}

// MARK: - BoolDescriptor

func TestBoolDescriptorSKU(t *testing.T) {
	id := uuid.New()
	descriptor := BoolDescriptor{UUID: id}

	assert.Equal(t, fmt.Sprintf("starfig::descriptor:bool:%s", id), descriptor.SKU())
}

func TestBoolDescriptorDefault(t *testing.T) {
	descriptor := BoolDescriptor{DefaultValue: true}

	assert.Equal(t, starlark.Bool(true), descriptor.Default())
}

func TestBoolDescriptorIsRequired(t *testing.T) {
	descriptor := BoolDescriptor{Required: true}

	assert.Equal(t, starlark.Bool(true), descriptor.IsRequired())
}

func TestBoolDescriptorEvaluate(t *testing.T) {
	descriptor := BoolDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltinWithCallback(func(args starlark.Tuple, kwargs []starlark.Tuple) {
				assert.Equal(t, starlark.Tuple{starlark.Bool(true)}, args)
				assert.ElementsMatch(t, []starlark.Tuple{}, kwargs)
			}),
		},
	}
	value, err := descriptor.Evaluate(&starlark.Thread{}, starlark.Bool(true))

	assert.Nil(t, err)
	assert.Equal(t, starlark.Bool(true), value)
}

func TestBoolDescriptorEvaluateInvalidType(t *testing.T) {
	descriptor := BoolDescriptor{}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.String("mock"))

	assert.ErrorContains(t, err, `Expected bool type but got "mock".`)
}

func TestBoolDescriptorEvaluateValidationError(t *testing.T) {
	descriptor := BoolDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingFunction("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.Bool(false))

	assert.ErrorContains(t, err, "yikes!")
}

func TestBoolDescriptorEvaluateValidationUserError(t *testing.T) {
	descriptor := BoolDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingBuiltin("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.Bool(false))

	assert.ErrorContains(t, err, "yikes!")
}

func TestBoolDescriptorString(t *testing.T) {
	id := uuid.New()
	descriptor := BoolDescriptor{
		UUID:         id,
		DefaultValue: true,
		Required:     true,
	}
	expected := fmt.Sprintf(`{"Type":"BoolDescriptor","Descriptor":{"UUID":"%s","DefaultValue":true,"Required":true,"Validations":null}}`, id)

	assert.Equal(t, expected, descriptor.String())
}

func TestBoolDescriptorType(t *testing.T) {
	descriptor := BoolDescriptor{}

	assert.Equal(t, "BoolDescriptor", descriptor.Type())
}

func TestBoolDescriptorFreeze(t *testing.T) {
	descriptor := BoolDescriptor{}
	descriptor.Freeze() // no-op
}

func TestBoolDescriptorTruth(t *testing.T) {
	descriptor := BoolDescriptor{}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestBoolDescriptorHash(t *testing.T) {
	hash, err := BoolDescriptor{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(349392354), hash)
}
