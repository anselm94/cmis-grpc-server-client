package server

import (
	"docserverclient/internal/server/model"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateInitData populates initial data like Repository, Type Definitions, Property Definitions etc.
func CreateInitData(db *gorm.DB) {
	log.Println("Populating some initial data...")

	// Property Definitions for Folder
	// (Property Defintions for folder & documents are to be created individually,
	// as each property definition is associated to either of the type definitions)
	propDefFolderName := model.PropertyDefinition{
		Name:        "cmis:name",
		Description: "Name",
		Type:        "string",
	}
	propDefFolderParentID := model.PropertyDefinition{
		Name:        "cmis:parentId",
		Description: "Parent ID",
		Type:        "id",
	}
	propDefFolderObjectID := model.PropertyDefinition{
		Name:        "cmis:objectId",
		Description: "Object ID",
		Type:        "id",
	}
	propDefFolderBaseTypeID := model.PropertyDefinition{
		Name:        "cmis:baseTypeId",
		Description: "Base Type ID",
		Type:        "id",
	}
	propDefFolderObjectTypeID := model.PropertyDefinition{
		Name:        "cmis:objectTypeId",
		Description: "Object Type ID",
		Type:        "id",
	}
	propDefFolderCreatedBy := model.PropertyDefinition{
		Name:        "cmis:createdBy",
		Description: "Created By",
		Type:        "string",
	}
	propDefFolderLastModifiedBy := model.PropertyDefinition{
		Name:        "cmis:lastModifiedBy",
		Description: "Last Modified By",
		Type:        "string",
	}
	propDefFolderCreationDate := model.PropertyDefinition{
		Name:        "cmis:creationDate",
		Description: "Creation Date",
		Type:        "datetime",
	}
	propDefFolderLastModificationDate := model.PropertyDefinition{
		Name:        "cmis:lastModificationDate",
		Description: "Last Modification Date",
		Type:        "datetime",
	}
	propDefFolderPath := model.PropertyDefinition{
		Name:        "cmis:path",
		Description: "Path",
		Type:        "string",
	}

	// Property Definitions for Documents
	propDefDocumentName := model.PropertyDefinition{
		Name:        "cmis:name",
		Description: "Name",
		Type:        "string",
	}
	propDefDocumentParentID := model.PropertyDefinition{
		Name:        "cmis:parentId",
		Description: "Parent ID",
		Type:        "id",
	}
	propDefDocumentObjectID := model.PropertyDefinition{
		Name:        "cmis:objectId",
		Description: "Object ID",
		Type:        "id",
	}
	propDefDocumentBaseTypeID := model.PropertyDefinition{
		Name:        "cmis:baseTypeId",
		Description: "Base Type ID",
		Type:        "id",
	}
	propDefDocumentObjectTypeID := model.PropertyDefinition{
		Name:        "cmis:objectTypeId",
		Description: "Object Type ID",
		Type:        "id",
	}
	propDefDocumentCreatedBy := model.PropertyDefinition{
		Name:        "cmis:createdBy",
		Description: "Created By",
		Type:        "string",
	}
	propDefDocumentLastModifiedBy := model.PropertyDefinition{
		Name:        "cmis:lastModifiedBy",
		Description: "Last Modified By",
		Type:        "string",
	}
	propDefDocumentCreationDate := model.PropertyDefinition{
		Name:        "cmis:creationDate",
		Description: "Creation Date",
		Type:        "datetime",
	}
	propDefDocumentLastModificationDate := model.PropertyDefinition{
		Name:        "cmis:lastModificationDate",
		Description: "Last Modification Date",
		Type:        "datetime",
	}
	propDefDocumentPath := model.PropertyDefinition{
		Name:        "cmis:path",
		Description: "Path",
		Type:        "string",
	}

	// Type Definitions
	typeDefFolder := &model.TypeDefinition{
		Name:        "cmis:folder",
		Description: "Folder",
		PropertyDefinitions: []*model.PropertyDefinition{
			&propDefFolderName,
			&propDefFolderParentID,
			&propDefFolderObjectID,
			&propDefFolderBaseTypeID,
			&propDefFolderObjectTypeID,
			&propDefFolderCreatedBy,
			&propDefFolderLastModifiedBy,
			&propDefFolderCreationDate,
			&propDefFolderLastModificationDate,
			&propDefFolderPath,
		},
	}
	typeDefDocument := &model.TypeDefinition{
		Name:        "cmis:document",
		Description: "Document",
		PropertyDefinitions: []*model.PropertyDefinition{
			&propDefDocumentName,
			&propDefDocumentParentID,
			&propDefDocumentObjectID,
			&propDefDocumentBaseTypeID,
			&propDefDocumentObjectTypeID,
			&propDefDocumentCreatedBy,
			&propDefDocumentLastModifiedBy,
			&propDefDocumentCreationDate,
			&propDefDocumentLastModificationDate,
			&propDefDocumentPath,
		},
	}

	//Create repository and type definitions & corresponding property definitions
	repository := &model.Repository{
		Name:        "SAP Drive",
		Description: "Cloud Storage by SAP",
		TypeDefinitions: []*model.TypeDefinition{
			typeDefFolder,
			typeDefDocument,
		},
	}
	if err := db.Create(&repository).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("Repository with name \"%s\" created", repository.Name)
		log.Printf("Type Definitions for \"%s\", \"%s\" are created in the repository \"%s\"", typeDefFolder.Name, typeDefDocument.Name, repository.Name)
		log.Printf("Property Definitions for \"%s\", \"%s\" are created and linked to Type Definition \"%s\"", propDefFolderName.Name, propDefFolderParentID.Name, typeDefFolder.Name)
		log.Printf("Property Definitions for \"%s\", \"%s\" are created and linked to Type Definition \"%s\"", propDefFolderName.Name, propDefFolderParentID.Name, typeDefDocument.Name)
	}

	//Create Root folder
	rootFolderName := "My Documents"
	rootFolder := &model.CmisObject{
		TypeDefinitionID: typeDefFolder.ID,
		RepositoryID:     repository.ID,
		Properties: []*model.CmisProperty{
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderName.ID,
				Value:                rootFolderName,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderParentID.ID,
				Value:                "",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderBaseTypeID.ID,
				Value:                typeDefFolder.Name,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderObjectTypeID.ID,
				Value:                typeDefFolder.Name,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderCreatedBy.ID,
				Value:                "default",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderLastModifiedBy.ID,
				Value:                "default",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderCreationDate.ID,
				Value:                strconv.Itoa(int(time.Now().Unix())),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderLastModificationDate.ID,
				Value:                strconv.Itoa(int(time.Now().Unix())),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderPath.ID,
				Value:                "/",
			},
		},
	}
	if err := db.Create(&rootFolder).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("A folder with name \"%s\" created", rootFolderName)
	}
	// Append CmisObjectID property after creating the object, as ID is dynamically generated
	if err := db.Model(&rootFolder).Association("Properties").Append(&model.CmisProperty{
		PropertyDefinitionID: propDefFolderObjectID.ID,
		Value:                fmt.Sprint(rootFolder.ID),
	}).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("CmisObjectID \"%d\" is updated to object with name \"%s\"", rootFolder.ID, rootFolderName)
	}

	//Connect newly created root folder to repository
	if err := db.Model(&repository).Update("root_folder_id", rootFolder.ID).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("The folder \"%s\" is assigned as the root folder for \"%s\" Repository", rootFolderName, repository.Name)
	}

	//Create a folder and a document
	folderAName := "Sales Receipts"
	folderA := &model.CmisObject{
		TypeDefinitionID: typeDefFolder.ID,
		RepositoryID:     repository.ID,
		Properties: []*model.CmisProperty{
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderName.ID,
				Value:                folderAName,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderParentID.ID,
				Value:                fmt.Sprint(rootFolder.ID),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderBaseTypeID.ID,
				Value:                typeDefFolder.Name,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderObjectTypeID.ID,
				Value:                typeDefFolder.Name,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderCreatedBy.ID,
				Value:                "default",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderLastModifiedBy.ID,
				Value:                "default",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderCreationDate.ID,
				Value:                strconv.Itoa(int(time.Now().Unix())),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderLastModificationDate.ID,
				Value:                strconv.Itoa(int(time.Now().Unix())),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderPath.ID,
				Value:                "/My Documents",
			},
		},
		Parents: []*model.CmisObject{
			rootFolder,
		},
	}
	if err := db.Create(&folderA).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("A folder named \"%s\" is created under the folder \"%s\"", folderAName, rootFolderName)
	}
	// Append CmisObjectID property after creating the object, as ID is dynamically generated
	if err := db.Model(&folderA).Association("Properties").Append(&model.CmisProperty{
		PropertyDefinitionID: propDefFolderObjectID.ID,
		Value:                fmt.Sprint(folderA.ID),
	}).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("CmisObjectID \"%d\" is updated to object with name \"%s\"", folderA.ID, folderAName)
	}

	documentAName := "Customer Feedback.xlsx"
	documentA := &model.CmisObject{
		TypeDefinitionID: typeDefDocument.ID,
		RepositoryID:     repository.ID,
		Properties: []*model.CmisProperty{
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentName.ID,
				Value:                documentAName,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentParentID.ID,
				Value:                fmt.Sprint(rootFolder.ID),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentBaseTypeID.ID,
				Value:                typeDefDocument.Name,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentObjectTypeID.ID,
				Value:                typeDefDocument.Name,
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentCreatedBy.ID,
				Value:                "default",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentLastModifiedBy.ID,
				Value:                "default",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentCreationDate.ID,
				Value:                strconv.Itoa(int(time.Now().Unix())),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentLastModificationDate.ID,
				Value:                strconv.Itoa(int(time.Now().Unix())),
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentPath.ID,
				Value:                "/My Documents",
			},
		},
		Parents: []*model.CmisObject{
			rootFolder,
		},
	}
	if err := db.Create(&documentA).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("A document named \"%s\" is created under the folder \"%s\"", documentAName, rootFolderName)
	}
	// Append CmisObjectID property after creating the object, as ID is dynamically generated
	if err := db.Model(&documentA).Association("Properties").Append(&model.CmisProperty{
		PropertyDefinitionID: propDefDocumentObjectID.ID,
		Value:                fmt.Sprint(documentA.ID),
	}).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("CmisObjectID \"%d\" is updated to object with name \"%s\"", documentA.ID, documentAName)
	}

	log.Println("Data population complete")
}

// DropTables drops all the tables if it already exists
func DropTables(db *gorm.DB) {
	db.DropTableIfExists("cmis_properties")
	db.DropTableIfExists("cmis_objects")
	db.DropTableIfExists("property_definitions")
	db.DropTableIfExists("type_definitions")
	db.DropTableIfExists("repositories")
	db.DropTableIfExists("filing")
}
