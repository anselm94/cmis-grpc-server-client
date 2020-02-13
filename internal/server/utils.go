package server

import (
	"docserverclient/internal/server/model"
	cmis "docserverclient/proto"
)

// ConvertRepositoryDaoToProto converts Repository Dao to Protobuf format
func ConvertRepositoryDaoToProto(repository *model.Repository) *cmis.Repository {
	if repository == nil || repository.ID == 0 {
		return nil
	}
	repositoryProto := &cmis.Repository{
		Id:              int32(repository.ID),
		Name:            repository.Name,
		Description:     repository.Description,
		RootFolder:      ConvertCmisObjectDaoToProto(repository.RootFolder, false),
		TypeDefinitions: make([]*cmis.TypeDefinition, len(repository.TypeDefinitions)),
	}
	for index, typeDef := range repository.TypeDefinitions {
		repositoryProto.TypeDefinitions[index] = ConvertTypeDefinitionDaoToProto(typeDef)
	}
	return repositoryProto
}

// ConvertTypeDefinitionDaoToProto converts TypeDefinition Dao to Protobuf format
func ConvertTypeDefinitionDaoToProto(typeDefinition *model.TypeDefinition) *cmis.TypeDefinition {
	if typeDefinition == nil || typeDefinition.ID == 0 {
		return nil
	}
	typeDefinitionProto := &cmis.TypeDefinition{
		Name:                typeDefinition.Name,
		Description:         typeDefinition.Description,
		PropertyDefinitions: make([]*cmis.PropertyDefinition, len(typeDefinition.PropertyDefinitions)),
	}
	for index, propertyDef := range typeDefinition.PropertyDefinitions {
		typeDefinitionProto.PropertyDefinitions[index] = ConvertPropertyDefinitionDaoToProto(propertyDef)
	}
	return typeDefinitionProto
}

// ConvertPropertyDefinitionDaoToProto converts PropertyDefinition Dao to Protobuf format
func ConvertPropertyDefinitionDaoToProto(propertyDefinition *model.PropertyDefinition) *cmis.PropertyDefinition {
	if propertyDefinition == nil || propertyDefinition.ID == 0 {
		return nil
	}
	propertyDefinitionProto := &cmis.PropertyDefinition{
		Name:        propertyDefinition.Name,
		Description: propertyDefinition.Description,
		Datatype:    propertyDefinition.Type,
	}
	return propertyDefinitionProto
}

// ConvertCmisObjectDaoToProto converts CmisObject Dao to Protobuf format
func ConvertCmisObjectDaoToProto(cmisObject *model.CmisObject, includeChildren bool) *cmis.CmisObject {
	if cmisObject == nil || cmisObject.ID == 0 {
		return nil
	}

	var parents []*cmis.CmisObject
	if cmisObject.Parents != nil {
		parents = make([]*cmis.CmisObject, len(cmisObject.Parents))
		for index, parent := range cmisObject.Parents {
			parents[index] = ConvertCmisObjectDaoToProto(parent, false)
		}
	}

	var properties []*cmis.CmisProperty
	if cmisObject.Properties != nil {
		properties = make([]*cmis.CmisProperty, len(cmisObject.Properties))
		for index, property := range cmisObject.Properties {
			properties[index] = ConvertCmisPropertyDaoToProto(property)
		}
	}

	var children []*cmis.CmisObject
	if cmisObject.Children != nil {
		children = make([]*cmis.CmisObject, len(cmisObject.Children))
		for index, child := range cmisObject.Children {
			children[index] = ConvertCmisObjectDaoToProto(child, false)
		}
	}

	object := &cmis.CmisObject{
		Id: &cmis.CmisObjectId{
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
func ConvertCmisPropertyDaoToProto(cmisProperty *model.CmisProperty) *cmis.CmisProperty {
	if cmisProperty == nil || cmisProperty.ID == 0 {
		return nil
	}
	propertyProto := &cmis.CmisProperty{
		PropertyDefinition: ConvertPropertyDefinitionDaoToProto(cmisProperty.PropertyDefinition),
		Value:              cmisProperty.Value,
	}
	return propertyProto
}
