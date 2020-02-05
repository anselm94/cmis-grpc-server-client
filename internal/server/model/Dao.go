package model

type Repository struct {
	ID              uint `gorm:"primary_key"`
	Name            string
	Description     string
	RootFolder      CmisObject `gorm:"foreignkey:RootFolderID"`
	RootFolderID    int
	TypeDefinitions []TypeDefinition
}

type TypeDefinition struct {
	ID                  uint `gorm:"primary_key"`
	Name                string
	Description         string
	RepositoryID        uint
	PropertyDefinitions []PropertyDefinition
}

type PropertyDefinition struct {
	ID               uint `gorm:"primary_key"`
	Name             string
	Description      string
	Type             string
	TypeDefinitionID uint
}

type CmisObject struct {
	ID               uint `gorm:"primary_key"`
	RepositoryID     uint
	TypeDefinition   TypeDefinition
	TypeDefinitionID uint
	Properties       []CmisProperty
	ParentID         uint
	Children         []*CmisObject `gorm:"many2many:filing;association_jointable_foreignkey:parent_id"`
}

type CmisProperty struct {
	ID                   uint `gorm:"primary_key"`
	CmisObjectID         uint
	PropertyDefinition   PropertyDefinition
	PropertyDefinitionID uint
	Value                string
}
