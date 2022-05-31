package native

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jathu/starfig/internal/target"
	"github.com/jathu/starfig/internal/tester"
	"github.com/stretchr/testify/assert"
)

func TestContextQueueSeenDescriptor(t *testing.T) {
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager := NewSchemaContextManager()
	manager.QueueSeenDescriptor(descriptor)

	assert.Empty(t, manager.builders)
	assert.Equal(t, 1, len(manager.queue))
	foundDescriptor, found := manager.queue[descriptor.SKU()]
	assert.True(t, found)
	assert.Equal(t, &descriptor, foundDescriptor)
}

func TestContextUpdateRecognizedSchema(t *testing.T) {
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager := NewSchemaContextManager()
	manager.QueueSeenDescriptor(descriptor)

	assert.Empty(t, manager.builders)
	assert.Equal(t, 1, len(manager.queue))

	fileTarget := target.FileTarget{
		StarverseDir: tester.GetTestStarverseDir(t),
		Package:      "example",
		Filename:     "STARFIG",
	}
	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName(descriptor.SKU()),
		"Supreme",
		fileTarget,
	)

	assert.Empty(t, manager.queue)

	foundItem, found := manager.builders[descriptor.SKU()]
	assert.True(t, found)
	assert.Equal(t, descriptor, foundItem.SchemaDescriptor)
	assert.Equal(t, "Supreme", foundItem.SchemaName)
	assert.Equal(t, fileTarget, foundItem.FileTarget)
}

func TestContextUpdateRecognizedSchemaWithMissingSKU(t *testing.T) {
	manager := NewSchemaContextManager()
	assert.Empty(t, manager.builders)
	assert.Empty(t, manager.queue)

	manager.UpdateRecognizedSchema(
		tester.MockBuiltinWithName("FakeSKU"),
		"Supreme",
		target.FileTarget{},
	)

	assert.Empty(t, manager.builders)
	assert.Empty(t, manager.queue)
}

func TestContextGetDescriptor(t *testing.T) {
	manager := NewSchemaContextManager()
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders["SupremeSKU"] = &SchemaContextItem{
		SchemaName:       "Supreme",
		SchemaDescriptor: descriptor,
		FileTarget:       target.FileTarget{},
	}

	foundDescriptor, found := manager.GetDescriptor("SupremeSKU")
	assert.True(t, found)
	assert.Equal(t, descriptor, foundDescriptor)
}

func TestContextGetDescriptorFoundInQueue(t *testing.T) {
	manager := NewSchemaContextManager()
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.queue["SupremeSKU"] = &descriptor

	foundDescriptor, found := manager.GetDescriptor("SupremeSKU")
	assert.True(t, found)
	assert.Equal(t, &descriptor, foundDescriptor)
}

func TestContextGetDescriptorNotFound(t *testing.T) {
	manager := NewSchemaContextManager()
	_, found := manager.GetDescriptor("SupremeSKU")
	assert.False(t, found)
}

func TestContextGetSchemaName(t *testing.T) {
	manager := NewSchemaContextManager()
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders[descriptor.SKU()] = &SchemaContextItem{
		SchemaName:       "Supreme",
		SchemaDescriptor: descriptor,
		FileTarget:       target.FileTarget{},
	}

	foundSchemaName, found := manager.GetSchemaName(descriptor)
	assert.True(t, found)
	assert.Equal(t, "Supreme", foundSchemaName)
}

func TestContextGetSchemaNameNotFound(t *testing.T) {
	manager := NewSchemaContextManager()
	descriptor := SchemaDescriptor{UUID: uuid.New()}
	foundSchemaName, found := manager.GetSchemaName(descriptor)
	assert.False(t, found)
	assert.Equal(t, "", foundSchemaName)
}

func TestContextEqualDescriptorSame(t *testing.T) {
	manager := NewSchemaContextManager()
	fileTarget := target.FileTarget{
		StarverseDir: tester.GetTestStarverseDir(t),
		Package:      "example",
		Filename:     "STARFIG",
	}
	firstDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders[firstDescriptor.SKU()] = &SchemaContextItem{
		SchemaName:       "Supreme",
		SchemaDescriptor: firstDescriptor,
		FileTarget:       fileTarget,
	}
	secondDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders[secondDescriptor.SKU()] = &SchemaContextItem{
		SchemaName:       "Supreme",
		SchemaDescriptor: secondDescriptor,
		FileTarget:       fileTarget,
	}

	assert.True(t, manager.EqualDescriptor(firstDescriptor, secondDescriptor))
}

func TestContextEqualDescriptorDifferentFile(t *testing.T) {
	manager := NewSchemaContextManager()
	firstDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders[firstDescriptor.SKU()] = &SchemaContextItem{
		SchemaName:       "Supreme",
		SchemaDescriptor: firstDescriptor,
		FileTarget: target.FileTarget{
			StarverseDir: tester.GetTestStarverseDir(t),
			Package:      "example",
			Filename:     "A",
		},
	}
	secondDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders[secondDescriptor.SKU()] = &SchemaContextItem{
		SchemaName:       "Supreme",
		SchemaDescriptor: secondDescriptor,
		FileTarget: target.FileTarget{
			StarverseDir: tester.GetTestStarverseDir(t),
			Package:      "example",
			Filename:     "B",
		},
	}

	assert.False(t, manager.EqualDescriptor(firstDescriptor, secondDescriptor))
}

func TestContextEqualDescriptorDifferentSchemaName(t *testing.T) {
	manager := NewSchemaContextManager()
	firstDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders[firstDescriptor.SKU()] = &SchemaContextItem{
		SchemaName:       "Supreme",
		SchemaDescriptor: firstDescriptor,
		FileTarget: target.FileTarget{
			StarverseDir: tester.GetTestStarverseDir(t),
			Package:      "example",
			Filename:     "SAME",
		},
	}
	secondDescriptor := SchemaDescriptor{UUID: uuid.New()}
	manager.builders[secondDescriptor.SKU()] = &SchemaContextItem{
		SchemaName:       "Patagonia",
		SchemaDescriptor: secondDescriptor,
		FileTarget: target.FileTarget{
			StarverseDir: tester.GetTestStarverseDir(t),
			Package:      "example",
			Filename:     "SAME",
		},
	}

	assert.False(t, manager.EqualDescriptor(firstDescriptor, secondDescriptor))
}
