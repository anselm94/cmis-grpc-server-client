package server

import (
	"docserverclient/internal/server/model"
	"fmt"

	"github.com/jinzhu/gorm"
)

func CreateInitData(db *gorm.DB) {
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

	//Create repository
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
	}

	//Create Root folder
	rootFolder := &model.CmisObject{
		TypeDefinitionID: typeDefFolder.ID,
		RepositoryID:     repository.ID,
		Properties: []*model.CmisProperty{
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderName.ID,
				Value:                "My Documents",
			},
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderParentID.ID,
				Value:                "",
			},
		},
	}
	if err := db.Create(&rootFolder).Error; err != nil {
		fmt.Println(err)
	}

	//Connect newly created root folder to repository
	if err := db.Model(&repository).Update("root_folder_id", rootFolder.ID).Error; err != nil {
		fmt.Println(err)
	}

	//Create a folder and a document
	folderA := &model.CmisObject{
		TypeDefinitionID: typeDefFolder.ID,
		RepositoryID:     repository.ID,
		Properties: []*model.CmisProperty{
			&model.CmisProperty{
				PropertyDefinitionID: propDefFolderName.ID,
				Value:                "Sales Receipts",
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
	}

	documentA := &model.CmisObject{
		TypeDefinitionID: typeDefDocument.ID,
		RepositoryID:     repository.ID,
		Properties: []*model.CmisProperty{
			&model.CmisProperty{
				PropertyDefinitionID: propDefDocumentName.ID,
				Value:                "Customer Feedback.xlsx",
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
	}

	fmt.Println(repository.Name)
}
