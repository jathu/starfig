package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

// MARK: - FloatProvider

func TestFloatProvider(t *testing.T) {
	validations := starlark.NewList([]starlark.Value{
		tester.MockBuiltin(),
		tester.MockBuiltin(),
		tester.MockBuiltin(),
	})

	value, err := FloatProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.Float(3.14)},
			{starlark.String("required"), starlark.Bool(true)},
			{starlark.String("validations"), validations},
		},
	)

	assert.Nil(t, err)

	provider := value.(FloatDescriptor)

	assert.Equal(t, starlark.Float(3.14), provider.Default())
	assert.Equal(t, starlark.Bool(true), provider.IsRequired())
	tester.AssertSameValidations(t, validations, provider.Validations)
}

func TestFloatProviderWithArguments(t *testing.T) {
	_, err := FloatProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.String("ok")},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, `Invalid positional arguments ("ok",) in Float().`)
}

func TestFloatProviderWithInvalidDefaultType(t *testing.T) {
	_, err := FloatProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.Bool(true)},
		},
	)

	assert.ErrorContains(t, err, `Expected default value to be float, but got True.`)
}

func TestFloatProviderWithInvalidRequiredType(t *testing.T) {
	_, err := FloatProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("required"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Expected required value to be bool, but got 416.`)
}

func TestFloatProviderWithInvalidValidationsType(t *testing.T) {
	_, err := FloatProvider(
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

func TestFloatProviderWithInvalidValidationsElementType(t *testing.T) {
	_, err := FloatProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.NewList([]starlark.Value{starlark.MakeInt(416)})},
		},
	)

	assert.ErrorContains(t, err, `Expected validation to be a functions, but got 416.`)
}

func TestFloatProviderWithUnknownKeyword(t *testing.T) {
	_, err := FloatProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("supreme"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Unknown keyword supreme in Float().`)
}

// MARK: - FloatDescriptor

func TestFloatDescriptorSKU(t *testing.T) {
	id := uuid.New()
	descriptor := FloatDescriptor{UUID: id}

	assert.Equal(t, fmt.Sprintf("starfig::descriptor:float:%s", id), descriptor.SKU())
}

func TestFloatDescriptorDefault(t *testing.T) {
	descriptor := FloatDescriptor{DefaultValue: starlark.Float(3.14)}

	assert.Equal(t, starlark.Float(3.14), descriptor.Default())
}

func TestFloatDescriptorIsRequired(t *testing.T) {
	descriptor := FloatDescriptor{Required: true}

	assert.Equal(t, starlark.Bool(true), descriptor.IsRequired())
}

func TestFloatDescriptorEvaluate(t *testing.T) {
	descriptor := FloatDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltinWithCallback(func(args starlark.Tuple, kwargs []starlark.Tuple) {
				assert.Equal(t, starlark.Tuple{starlark.Float(3.14)}, args)
				assert.ElementsMatch(t, []starlark.Tuple{}, kwargs)
			}),
		},
	}
	value, err := descriptor.Evaluate(&starlark.Thread{}, starlark.Float(3.14))

	assert.Nil(t, err)
	assert.Equal(t, starlark.Float(3.14), value)
}

func TestFloatDescriptorEvaluateInvalidType(t *testing.T) {
	descriptor := FloatDescriptor{}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.String("mock"))

	assert.ErrorContains(t, err, `Expected float type but got "mock".`)
}

func TestFloatDescriptorEvaluateValidationError(t *testing.T) {
	descriptor := FloatDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingFunction("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.Float(3.14))

	assert.ErrorContains(t, err, "yikes!")
}

func TestFloatDescriptorEvaluateValidationUserError(t *testing.T) {
	descriptor := FloatDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingBuiltin("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.Float(3.14))

	assert.ErrorContains(t, err, "yikes!")
}
func TestFloatDescriptorString(t *testing.T) {
	id := uuid.New()
	descriptor := FloatDescriptor{
		UUID:         id,
		DefaultValue: starlark.Float(3.14),
		Required:     true,
	}
	expected := fmt.Sprintf(`{"Type":"FloatDescriptor","Descriptor":{"UUID":"%s","DefaultValue":3.14,"Required":true,"Validations":null}}`, id)

	assert.Equal(t, expected, descriptor.String())
}

func TestFloatDescriptorType(t *testing.T) {
	descriptor := FloatDescriptor{}

	assert.Equal(t, "FloatDescriptor", descriptor.Type())
}

func TestFloatDescriptorFreeze(t *testing.T) {
	descriptor := FloatDescriptor{}
	descriptor.Freeze() // no-op
}

func TestFloatDescriptorTruth(t *testing.T) {
	descriptor := FloatDescriptor{}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestFloatDescriptorHash(t *testing.T) {
	hash, err := FloatDescriptor{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(2649988975), hash)
}
