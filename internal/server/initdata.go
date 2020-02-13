package server

import (
	"docserverclient/internal/server/model"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

// CreateInitData populates initial data like Repository, Type Definitions, Property Definitions etc.
func CreateInitData(db *gorm.DB) {
	log.Println("Populating some initial data...")

	propDefFolderName := &model.PropertyDefinition{
		Name:        "cmis:name",
		Description: "Name",
		Type:        "string",
	}
	propDefFolderParentID := &model.PropertyDefinition{
		Name:        "cmis:parentId",
		Description: "Parent ID",
		Type:        "integer",
	}
	propDefDocumentName := &model.PropertyDefinition{
		Name:        "cmis:name",
		Description: "Name",
		Type:        "string",
	}
	propDefDocumentParentID := &model.PropertyDefinition{
		Name:        "cmis:parentId",
		Description: "Parent ID",
		Type:        "integer",
	}
	typeDefFolder := &model.TypeDefinition{
		Name:        "cmis:folder",
		Description: "Folder",
		PropertyDefinitions: []*model.PropertyDefinition{
			propDefFolderName,
			propDefFolderParentID,
		},
	}
	typeDefDocument := &model.TypeDefinition{
		Name:        "cmis:document",
		Description: "Document",
		PropertyDefinitions: []*model.PropertyDefinition{
			propDefDocumentName,
			propDefDocumentParentID,
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
		log.Printf("Property Definitions for \"%s\", \"%s\" are created and linked to Type Definition \"%s\"", propDefDocumentName.Name, propDefDocumentParentID.Name, typeDefDocument.Name)
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
		},
	}
	if err := db.Create(&rootFolder).Error; err != nil {
		fmt.Println(err)
	} else {
		log.Printf("A folder with name \"%s\" created", rootFolderName)
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
