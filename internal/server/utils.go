package server

import (
	"docserverclient/internal/server/model"
	cmisproto "docserverclient/proto"
)

// ConvertRepositoryDaoToProto converts Repository Dao to Protobuf format
func ConvertRepositoryDaoToProto(repository *model.Repository) *cmisproto.Repository {
	if repository == nil || repository.ID == 0 {
		return nil
	}
	repositoryProto := &cmisproto.Repository{
		Id:              int32(repository.ID),
		Name:            repository.Name,
		Description:     repository.Description,
		RootFolder:      ConvertCmisObjectDaoToProto(repository.RootFolder, false),
		TypeDefinitions: make([]*cmisproto.TypeDefinition, len(repository.TypeDefinitions)),
	}
	for index, typeDef := range repository.TypeDefinitions {
		repositoryProto.TypeDefinitions[index] = ConvertTypeDefinitionDaoToProto(typeDef)
	}
	return repositoryProto
}

// ConvertTypeDefinitionDaoToProto converts TypeDefinition Dao to Protobuf format
func ConvertTypeDefinitionDaoToProto(typeDefinition *model.TypeDefinition) *cmisproto.TypeDefinition {
	if typeDefinition == nil || typeDefinition.ID == 0 {
		return nil
	}
	typeDefinitionProto := &cmisproto.TypeDefinition{
		Name:                typeDefinition.Name,
		Description:         typeDefinition.Description,
		PropertyDefinitions: make([]*cmisproto.PropertyDefinition, len(typeDefinition.PropertyDefinitions)),
	}
	for index, propertyDef := range typeDefinition.PropertyDefinitions {
		typeDefinitionProto.PropertyDefinitions[index] = ConvertPropertyDefinitionDaoToProto(propertyDef)
	}
	return typeDefinitionProto
}

// ConvertPropertyDefinitionDaoToProto converts PropertyDefinition Dao to Protobuf format
func ConvertPropertyDefinitionDaoToProto(propertyDefinition *model.PropertyDefinition) *cmisproto.PropertyDefinition {
	if propertyDefinition == nil || propertyDefinition.ID == 0 {
		return nil
	}
	propertyDefinitionProto := &cmisproto.PropertyDefinition{
		Name:        propertyDefinition.Name,
		Description: propertyDefinition.Description,
		Datatype:    propertyDefinition.Type,
	}
	return propertyDefinitionProto
}

// ConvertCmisObjectDaoToProto converts CmisObject Dao to Protobuf format
func ConvertCmisObjectDaoToProto(cmisObject *model.CmisObject, includeChildren bool) *cmisproto.CmisObject {
	if cmisObject == nil || cmisObject.ID == 0 {
		return nil
	}

	var parents []*cmisproto.CmisObject
	if cmisObject.Parents != nil {
		parents = make([]*cmisproto.CmisObject, len(cmisObject.Parents))
		for index, parent := range cmisObject.Parents {
			parents[index] = ConvertCmisObjectDaoToProto(parent, false)
		}
	}

	var properties []*cmisproto.CmisProperty
	if cmisObject.Properties != nil {
		properties = make([]*cmisproto.CmisProperty, len(cmisObject.Properties))
		for index, property := range cmisObject.Properties {
			properties[index] = ConvertCmisPropertyDaoToProto(property)
		}
	}

	var children []*cmisproto.CmisObject
	if cmisObject.Children != nil {
		children = make([]*cmisproto.CmisObject, len(cmisObject.Children))
		for index, child := range cmisObject.Children {
			children[index] = ConvertCmisObjectDaoToProto(child, false)
		}
	}

	object := &cmisproto.CmisObject{
		Id: &cmisproto.CmisObjectId{
			Id: int32(cmisObject.ID),
		},
		TypeDefinition: ConvertTypeDefinitionDaoToProto(cmisObject.TypeDefinition),
		Properties:     properties,
		Children:       children,
		Parents:        parents,
	}
	return object
}

// ConvertCmisPropertyDaoToProto converts CmisProperty Dao to Protobuf format
func ConvertCmisPropertyDaoToProto(cmisProperty *model.CmisProperty) *cmisproto.CmisProperty {
	if cmisProperty == nil || cmisProperty.ID == 0 {
		return nil
	}
	propertyProto := &cmisproto.CmisProperty{
		PropertyDefinition: ConvertPropertyDefinitionDaoToProto(cmisProperty.PropertyDefinition),
		Value:              cmisProperty.Value,
	}
	return propertyProto
}

func GetProperty(cmisObjectWithPropertiesAndPropertyDefinitions *model.CmisObject, cmisPropertyID string) string {
	if len(cmisObjectWithPropertiesAndPropertyDefinitions.Properties) > 0 {
		for _, cmisProperty := range cmisObjectWithPropertiesAndPropertyDefinitions.Properties {
			if cmisProperty.PropertyDefinition.Name == cmisPropertyID {
				return ""
			}
		}
	}
	return ""
}
