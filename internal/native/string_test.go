package native

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

// MARK: - StringProvider

func TestStringProvider(t *testing.T) {
	validations := starlark.NewList([]starlark.Value{
		tester.MockBuiltin(),
		tester.MockBuiltin(),
		tester.MockBuiltin(),
	})

	value, err := StringProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.String("hello")},
			{starlark.String("required"), starlark.Bool(true)},
			{starlark.String("validations"), validations},
		},
	)

	assert.Nil(t, err)

	provider := value.(StringDescriptor)

	assert.Equal(t, starlark.String("hello"), provider.Default())
	assert.Equal(t, starlark.Bool(true), provider.IsRequired())
	tester.AssertSameValidations(t, validations, provider.Validations)
}

func TestStringProviderWithArguments(t *testing.T) {
	_, err := StringProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{starlark.String("ok")},
		[]starlark.Tuple{},
	)

	assert.ErrorContains(t, err, `Invalid positional arguments ("ok",) in String().`)
}

func TestStringProviderWithInvalidDefaultType(t *testing.T) {
	_, err := StringProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("default"), starlark.Bool(true)},
		},
	)

	assert.ErrorContains(t, err, `Expected default value to be string, but got True.`)
}

func TestStringProviderWithInvalidRequiredType(t *testing.T) {
	_, err := StringProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("required"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Expected required value to be bool, but got 416.`)
}

func TestStringProviderWithInvalidValidationsType(t *testing.T) {
	_, err := StringProvider(
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

func TestStringProviderWithInvalidValidationsElementType(t *testing.T) {
	_, err := StringProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("validations"), starlark.NewList([]starlark.Value{starlark.MakeInt(416)})},
		},
	)

	assert.ErrorContains(t, err, `Expected validation to be a functions, but got 416.`)
}

func TestStringProviderWithUnknownKeyword(t *testing.T) {
	_, err := StringProvider(
		&starlark.Thread{},
		tester.MockBuiltin(),
		starlark.Tuple{},
		[]starlark.Tuple{
			{starlark.String("supreme"), starlark.MakeInt(416)},
		},
	)

	assert.ErrorContains(t, err, `Unknown keyword supreme in String().`)
}

// MARK: - StringDescriptor

func TestStringDescriptorSKU(t *testing.T) {
	id := uuid.New()
	descriptor := StringDescriptor{UUID: id}

	assert.Equal(t, fmt.Sprintf("starfig::descriptor:string:%s", id), descriptor.SKU())
}

func TestStringDescriptorDefault(t *testing.T) {
	descriptor := StringDescriptor{DefaultValue: starlark.String("hello")}

	assert.Equal(t, starlark.String("hello"), descriptor.Default())
}

func TestStringDescriptorIsRequired(t *testing.T) {
	descriptor := StringDescriptor{Required: true}

	assert.Equal(t, starlark.Bool(true), descriptor.IsRequired())
}

func TestStringDescriptorEvaluate(t *testing.T) {
	descriptor := StringDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltinWithCallback(func(args starlark.Tuple, kwargs []starlark.Tuple) {
				assert.Equal(t, starlark.Tuple{starlark.String("supreme")}, args)
				assert.ElementsMatch(t, []starlark.Tuple{}, kwargs)
			}),
		},
	}
	value, err := descriptor.Evaluate(&starlark.Thread{}, starlark.String("supreme"))

	assert.Nil(t, err)
	assert.Equal(t, starlark.String("supreme"), value)
}

func TestStringDescriptorEvaluateInvalidType(t *testing.T) {
	descriptor := StringDescriptor{}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.MakeInt(416))

	assert.ErrorContains(t, err, `Expected string type but got 416.`)
}

func TestStringDescriptorEvaluateValidationError(t *testing.T) {
	descriptor := StringDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingFunction("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.String("supreme"))

	assert.ErrorContains(t, err, "yikes!")
}

func TestStringDescriptorEvaluateValidationUserError(t *testing.T) {
	descriptor := StringDescriptor{
		Validations: []starlark.Callable{
			tester.MockBuiltin(),
			tester.MockFailingBuiltin("yikes!"),
		},
	}
	_, err := descriptor.Evaluate(&starlark.Thread{}, starlark.String("supreme"))

	assert.ErrorContains(t, err, "yikes!")
}

func TestStringDescriptorString(t *testing.T) {
	id := uuid.New()
	descriptor := StringDescriptor{
		UUID:         id,
		DefaultValue: starlark.String("hello"),
		Required:     true,
	}
	expected := fmt.Sprintf(`{"Type":"StringDescriptor","Descriptor":{"UUID":"%s","DefaultValue":"hello","Required":true,"Validations":null}}`, id)

	assert.Equal(t, expected, descriptor.String())
}

func TestStringDescriptorType(t *testing.T) {
	descriptor := StringDescriptor{}

	assert.Equal(t, "StringDescriptor", descriptor.Type())
}

func TestStringDescriptorFreeze(t *testing.T) {
	descriptor := StringDescriptor{}
	descriptor.Freeze() // no-op
}

func TestStringDescriptorTruth(t *testing.T) {
	descriptor := StringDescriptor{}

	assert.Equal(t, starlark.Bool(true), descriptor.Truth())
}

func TestStringDescriptorHash(t *testing.T) {
	hash, err := StringDescriptor{}.Hash()

	assert.Nil(t, err)
	assert.Equal(t, uint32(3208191298), hash)
}
