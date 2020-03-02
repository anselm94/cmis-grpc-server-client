package cmisserver

import (
	"docserverclient"
	cmismodel "docserverclient/internal/cmisserver/model"
	cmisproto "docserverclient/proto"
	"fmt"
)

var config *docserverclient.Config = docserverclient.NewDefaultConfig()

func ConvertRepositoryProtoToCmis(repositoryProto *cmisproto.Repository) *cmismodel.Repository {
	cmisRepository := cmismodel.Repository{
		RepositoryID:          fmt.Sprint(repositoryProto.Id),
		RepositoryName:        repositoryProto.Name,
		RepositoryDescription: repositoryProto.Description,
		VendorName:            "SAP",
		ProductName:           repositoryProto.Name,
		ProductVersion:        "1.0",
		RootFolderID:          fmt.Sprint(repositoryProto.RootFolder.Id.Id),
		Capabilities: &cmismodel.RepositoryCapabilities{
			CapabilityContentStreamUpdatability: "anytime",
			CapabilityChanges:                   "none",
			CapabilityRenditions:                "none",
			CapabilityGetDescendants:            false,
			CapabilityGetFolderTree:             false,
			CapabilityMultifiling:               false,
			CapabilityUnfiling:                  false,
			CapabilityVersionSpecificFiling:     false,
			CapabilityPWCSearchable:             false,
			CapabilityPWCUpdatable:              false,
			CapabilityAllVersionsSearchable:     false,
			CapabilityOrderBy:                   "none",
			CapabilityQuery:                     "none",
			CapabilityJoin:                      "none",
			CapabilityACL:                       "discover",
			CapabilityCreatablePropertyTypes: map[string][]string{
				"canCreate": []string{},
			},
			CapabilityNewTypeSettableAttributes: map[string]bool{
				"id":                       false,
				"localName":                false,
				"localNamespace":           false,
				"displayName":              false,
				"queryName":                false,
				"description":              false,
				"creatable":                false,
				"fileable":                 false,
				"queryable":                false,
				"fulltextIndexed":          false,
				"includedInSupertypeQuery": false,
				"controllablePolicy":       false,
				"controllableACL":          false,
			},
		},
		ACLCapabilities: &cmismodel.ACLCapabilities{
			SupportedPermissions: "basic",
			Propagation:          "objectonly",
			Permissions: []*cmismodel.Permission{
				&cmismodel.Permission{
					Permission:  "cmis:all",
					Description: "All",
				},
			},
			PermissionMapping: []*cmismodel.PermissionMap{},
		},
		LatestChangeLogToken: nil,
		CmisVersionSupported: "1.1",
		ThinClientURI:        "",
		ChangesIncomplete:    true,
		ChangesOnType:        []string{},
		ExtendedFeatures:     []*cmismodel.RepositoryExtendedFeature{},
		RepositoryURL:        "http://" + config.CmisAppHost + config.CmisAppPort + "/browser/" + fmt.Sprint(repositoryProto.Id),
		RootFolderURL:        "http://" + config.CmisAppHost + config.CmisAppPort + "/browser/" + fmt.Sprint(repositoryProto.Id) + "/" + fmt.Sprint(repositoryProto.RootFolder.Id.Id),
	}
	return &cmisRepository
}

func ConvertTypeDefinitionsProtoToCmis(typedefinitions []*cmisproto.TypeDefinition, includePropertyDefinitions bool) []*cmismodel.TypeDefinition {
	cmisTypeDefinitions := make([]*cmismodel.TypeDefinition, len(typedefinitions))
	for index, typedefinition := range typedefinitions {
		cmisTypeDefinition := ConvertTypeDefinitionProtoToCmis(typedefinition, includePropertyDefinitions)
		cmisTypeDefinitions[index] = cmisTypeDefinition
	}
	return cmisTypeDefinitions
}

func ConvertTypeDefinitionProtoToCmis(typedefinition *cmisproto.TypeDefinition, includePropertyDefinitions bool) *cmismodel.TypeDefinition {
	cmisTypeDefinition := cmismodel.TypeDefinition{
		ID:                       typedefinition.Name,
		LocalName:                typedefinition.Name,
		LocalNamespace:           "grpc-cmis",
		DisplayName:              typedefinition.Name,
		QueryName:                typedefinition.Name,
		Description:              typedefinition.Name,
		BaseID:                   typedefinition.Name,
		Creatable:                false,
		Fileable:                 false,
		Queryable:                false,
		FulltextIndexed:          false,
		IncludedInSupertypeQuery: false,
		ControllablePolicy:       false,
		ControllableACL:          false,
		TypeMutability: map[string]bool{
			"create": false,
			"update": false,
			"delete": false,
		},
		Versionable:          false,
		ContentStreamAllowed: "notallowed",
	}
	if typedefinition.Name == "cmis:document" {
		cmisTypeDefinition.Fileable = true
		cmisTypeDefinition.ContentStreamAllowed = "allowed"
	}
	if includePropertyDefinitions {
		propertyDefinitions := make([]*cmismodel.PropertyDefinition, len(typedefinition.PropertyDefinitions))
		for indexIn, propertydefinition := range typedefinition.PropertyDefinitions {
			propertyDefinitions[indexIn] = &cmismodel.PropertyDefinition{
				ID:            propertydefinition.Name,
				LocalName:     propertydefinition.Description,
				DisplayName:   propertydefinition.Description,
				QueryName:     propertydefinition.Name,
				Description:   propertydefinition.Description,
				PropertyType:  propertydefinition.Datatype,
				Cardinality:   "single",
				Updateability: "read",
				Inherited:     false,
				Required:      true,
				Queryable:     false,
				Orderable:     false,
			}
		}
		cmisTypeDefinition.PropertyDefinitions = propertyDefinitions
	}
	return &cmisTypeDefinition
}

func ConvertCmisObjectProtoToCmis(cmisobject *cmisproto.CmisObject, isSuccinctProperties bool, includeAllowableActions bool, includeACL bool) *cmismodel.CmisObject {
	cmisObject := cmismodel.CmisObject{}
	if isSuccinctProperties {
		cmisObject.SuccinctProperties = make(map[string]string)
		for _, cmisproperty := range cmisobject.Properties {
			cmisObject.SuccinctProperties[cmisproperty.PropertyDefinition.Name] = cmisproperty.Value
		}
	} else {
		cmisObjectProperties := make([]*cmismodel.CmisProperty, len(cmisobject.Properties))
		for index, cmisproperty := range cmisobject.Properties {
			cmisObjectProperties[index] = &cmismodel.CmisProperty{
				ID:          cmisproperty.PropertyDefinition.Name,
				Type:        cmisproperty.PropertyDefinition.Datatype,
				Cardinality: "single",
				Value:       cmisproperty.Value,
			}
		}
		cmisObject.Properties = &cmisObjectProperties
	}
	if includeAllowableActions {
		cmisObject.AllowableActions = &cmismodel.AllowableActions{
			CanDeleteObject:           true,
			CanUpdateProperties:       false,
			CanGetFolderTree:          false,
			CanGetProperties:          true,
			CanGetObjectRelationships: false,
			CanGetObjectParents:       false,
			CanGetFolderParent:        false,
			CanGetDescendants:         false,
			CanMoveObject:             false,
			CanDeleteContentStream:    false,
			CanCheckOut:               false,
			CanCancelCheckOut:         false,
			CanCheckIn:                false,
			CanSetContentStream:       false,
			CanGetAllVersions:         false,
			CanAddObjectToFolder:      true,
			CanRemoveObjectFromFolder: true,
			CanGetContentStream:       false,
			CanApplyPolicy:            false,
			CanGetAppliedPolicies:     false,
			CanRemovePolicy:           false,
			CanGetChildren:            true,
			CanCreateDocument:         true,
			CanCreateFolder:           true,
			CanCreateRelationship:     false,
			CanCreateItem:             false,
			CanDeleteTree:             false,
			CanGetRenditions:          false,
			CanGetACL:                 true,
			CanApplyACL:               false,
		}
	}
	if includeACL {
		isExactACL := true
		cmisObject.ACL = &cmismodel.ACL{
			ACEs: []*cmismodel.ACE{
				&cmismodel.ACE{
					Principal: &cmismodel.Principal{
						PrincipalID: "default",
					},
					IsDirect: true,
					Permissions: []string{
						"cmis:all",
					},
				},
			},
		}
		cmisObject.ExactACL = &isExactACL
	}
	return &cmisObject
}
