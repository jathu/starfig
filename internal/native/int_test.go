package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

// MARK: - IntProvider

func TestIntProvider(t *testing.T) {
	validations := starlark.NewList([]starlark.Value{
		tester.MockBuiltin(),
		tester.MockBuiltin(),
		tester.MockBuiltin(),
	})

	value, err := IntProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.MakeInt(416)},
			{starlark.String("required"), starlark.Bool(true)},
			{starlark.String("validations"), validations},
		},
	)

	assert.Nil(t, err)

	provider := value.(IntDescriptor)

	assert.Equal(t, starlark.MakeInt(416), provider.Default())
	assert.Equal(t, starlark.Bool(true), provider.IsRequired())
	tester.AssertSameValidations(t, validations, provider.Validations)
}

func TestIntProviderWithArguments(t *testing.T) {
	_, err := IntProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.String("ok")},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, `Invalid positional arguments ("ok",) in Int().`)
}

func TestIntProviderWithInvalidDefaultType(t *testing.T) {
	_, err := IntProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.Bool(true)},
		},
	)

	assert.ErrorContains(t, err, `Expected default value to be int, but got True.`)
}

func TestIntProviderWithInvalidRequiredType(t *testing.T) {
	_, err := IntProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("required"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Expected required value to be bool, but got 416.`)
}

func TestIntProviderWithInvalidValidationsType(t *testing.T) {
	_, err := IntProvider(
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

func TestIntProviderWithInvalidValidationsElementType(t *testing.T) {
	_, err := IntProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.NewList([]starlark.Value{starlark.MakeInt(416)})},
		},
	)

	assert.ErrorContains(t, err, `Expected validation to be a functions, but got 416.`)
}

func TestIntProviderWithUnknownKeyword(t *testing.T) {
	_, err := IntProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("supreme"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Unknown keyword supreme in Int().`)
}

// MARK: - IntDescriptor

func TestIntDescriptorSKU(t *testing.T) {
	id := uuid.New()
	descriptor := IntDescriptor{UUID: id}

	assert.Equal(t, fmt.Sprintf("starfig::descriptor:int:%s", id), descriptor.SKU())
}

func TestIntDescriptorDefault(t *testing.T) {
	descriptor := IntDescriptor{DefaultValue: starlark.MakeInt(416)}

	assert.Equal(t, starlark.MakeInt(416), descriptor.Default())
}

func TestIntDescriptorIsRequired(t *testing.T) {
	descriptor := IntDescriptor{Required: true}

	assert.Equal(t, starlark.Bool(true), descriptor.IsRequired())
}

func TestIntDescriptorEvaluate(t *testing.T) {
	descriptor := IntDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltinWithCallback(func(args starlark.Tuple, kwargs []starlark.Tuple) {
				assert.Equal(t, starlark.Tuple{starlark.MakeInt(416)}, args)
				assert.ElementsMatch(t, []starlark.Tuple{}, kwargs)
			}),
		},
	}
	value, err := descriptor.Evaluate(&starlark.Thread{}, starlark.MakeInt(416))

	assert.Nil(t, err)
	assert.Equal(t, starlark.MakeInt(416), value)
}

func TestIntDescriptorEvaluateInvalidType(t *testing.T) {
	descriptor := IntDescriptor{}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.String("mock"))

	assert.ErrorContains(t, err, `Expected int type but got "mock".`)
}

func TestIntDescriptorEvaluateValidationError(t *testing.T) {
	descriptor := IntDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingFunction("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.MakeInt(416))

	assert.ErrorContains(t, err, "yikes!")
}

func TestIntDescriptorEvaluateValidationUserError(t *testing.T) {
	descriptor := IntDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingBuiltin("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.MakeInt(416))

	assert.ErrorContains(t, err, "yikes!")
}

func TestIntDescriptorString(t *testing.T) {
	id := uuid.New()
	descriptor := IntDescriptor{
		UUID:         id,
		DefaultValue: starlark.MakeInt(416),
		Required:     true,
	}
	expected := fmt.Sprintf(`{"Type":"IntDescriptor","Descriptor":{"UUID":"%s","DefaultValue":{},"Required":true,"Validations":null}}`, id)

	assert.Equal(t, expected, descriptor.String())
}

func TestIntDescriptorType(t *testing.T) {
	descriptor := IntDescriptor{}

	assert.Equal(t, "IntDescriptor", descriptor.Type())
}

func TestIntDescriptorFreeze(t *testing.T) {
	descriptor := IntDescriptor{}
	descriptor.Freeze() // no-op
}

func TestIntDescriptorTruth(t *testing.T) {
	descriptor := IntDescriptor{}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestIntDescriptorHash(t *testing.T) {
	hash, err := IntDescriptor{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(1043499210), hash)
}
