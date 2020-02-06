package model

type Repository struct {
	ID              uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name            string
	Description     string
	RootFolder      *CmisObject `gorm:"foreignkey:RootFolderID"`
	RootFolderID    int
	TypeDefinitions []*TypeDefinition
}

type TypeDefinition struct {
	ID                  uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name                string
	Description         string
	RepositoryID        uint
	PropertyDefinitions []*PropertyDefinition
}

type PropertyDefinition struct {
	ID               uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name             string
	Description      string
	Type             string
	TypeDefinition   *TypeDefinition
	TypeDefinitionID uint
}

type CmisObject struct {
	ID               uint `gorm:"primary_key;AUTO_INCREMENT"`
	RepositoryID     uint
	TypeDefinition   *TypeDefinition
	TypeDefinitionID uint
	Properties       []*CmisProperty
	Children         []*CmisObject `gorm:"many2many:filing;association_jointable_foreignkey:object_id;jointable_foreignkey:parent_id;"`
	Parents          []*CmisObject `gorm:"many2many:filing;association_jointable_foreignkey:parent_id;jointable_foreignkey:object_id;"`
}

type CmisProperty struct {
	ID                   uint `gorm:"primary_key;AUTO_INCREMENT"`
	CmisObjectID         uint
	PropertyDefinition   *PropertyDefinition
	PropertyDefinitionID uint
	Value                string
}
