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
	Children         []*CmisObject `gorm:"many2many:filing;association_jointable_foreignkey:child_id;jointable_foreignkey:cmis_object_id;"`
	Parents          []*CmisObject `gorm:"many2many:filing;association_jointable_foreignkey:cmis_object_id;jointable_foreignkey:child_id;"`
}

type CmisProperty struct {
	ID                   uint `gorm:"primary_key"`
	CmisObjectID         uint
	PropertyDefinition   PropertyDefinition
	PropertyDefinitionID uint
	Value                string
}