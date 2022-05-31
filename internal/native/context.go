package native

import (
	"github.com/jathu/starfig/internal/target"
	"go.starlark.net/starlark"
)

var SchemaContextManagerThreadKey string = "starfig-schema-context-manager"

// MARK: - SchemaContextItem

type SchemaContextItem struct {
	SchemaName       string
	SchemaDescriptor SchemaDescriptor
	FileTarget       target.FileTarget
}

// MARK: - SchemaContextManager

type SchemaContextManager struct {
	builders map[string]*SchemaContextItem
	queue    map[string]*SchemaDescriptor
}

func NewSchemaContextManager() SchemaContextManager {
	return SchemaContextManager{
		builders: map[string]*SchemaContextItem{},
		queue:    map[string]*SchemaDescriptor{},
	}
}

func (manager SchemaContextManager) QueueSeenDescriptor(descriptor SchemaDescriptor) {
	manager.queue[descriptor.SKU()] = &descriptor
}

func (manager SchemaContextManager) UpdateRecognizedSchema(
	schemaBuilder *starlark.Builtin, schemaName string, fileTarget target.FileTarget) {
	recognizedSchemaSKU := schemaBuilder.Name()

	descriptor, found := manager.queue[recognizedSchemaSKU]
	if found {
		manager.builders[recognizedSchemaSKU] = &SchemaContextItem{
			SchemaName:       schemaName,
			SchemaDescriptor: *descriptor,
			FileTarget:       fileTarget,
		}

		delete(manager.queue, recognizedSchemaSKU)
	}
}

func (manager SchemaContextManager) GetDescriptor(descriptorSKU string) (Descriptor, bool) {
	item, ok := manager.builders[descriptorSKU]
	if ok {
		return item.SchemaDescriptor, true
	} else {
		// If the schema has not been loaded using "load", it might be still defined in the current
		// file, thus we don't currently know it's name or file target. So loop through the queue
		// to find the SKU.
		descriptor, found := manager.queue[descriptorSKU]
		if found {
			return descriptor, true
		}
	}

	return BoolDescriptor{}, false
}

func (manager SchemaContextManager) GetSchemaName(descriptor Descriptor) (string, bool) {
	item, ok := manager.builders[descriptor.SKU()]
	if ok {
		return item.SchemaName, true
	} else {
		return "", false
	}
}

// Starlark loads modules multiple times, this causes our schema builder to run
// multiple times. Which causes the schema to be defined multiple times. However,
// the schema builder doesn't have context of what file it's in — so we determine
// if schemas are the same if they are from the same file and have the same name.
func (manager SchemaContextManager) EqualDescriptor(
	left SchemaDescriptor, right SchemaDescriptor) bool {
	leftDesc, leftOk := manager.builders[left.SKU()]
	rightDesc, rightOk := manager.builders[right.SKU()]
	if !leftOk || !rightOk {
		return false
	}

	isSameFile := leftDesc.FileTarget.Target() == rightDesc.FileTarget.Target()
	isSameSchemaName := leftDesc.SchemaName == rightDesc.SchemaName

	return isSameFile && isSameSchemaName
}
