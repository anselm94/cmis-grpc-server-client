package model

// Repository Dao model
type Repository struct {
	ID              uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name            string
	Description     string
	RootFolder      *CmisObject `gorm:"foreignkey:RootFolderID"`
	RootFolderID    int
	TypeDefinitions []*TypeDefinition // 1:n - has many
}

// TypeDefinition Dao model
type TypeDefinition struct {
	ID                  uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name                string
	Description         string
	RepositoryID        uint
	PropertyDefinitions []*PropertyDefinition // 1:n - has many
}

// PropertyDefinition Dao model
type PropertyDefinition struct {
	ID               uint `gorm:"primary_key;AUTO_INCREMENT"`
	Name             string
	Description      string
	Type             string
	TypeDefinition   *TypeDefinition // n:1 - belongs to
	TypeDefinitionID uint
}

// CmisObject Dao model
type CmisObject struct {
	ID               uint `gorm:"primary_key;AUTO_INCREMENT"`
	RepositoryID     uint
	TypeDefinition   *TypeDefinition // 1:1 - belongs to
	TypeDefinitionID uint
	Properties       []*CmisProperty // 1:n - has many
	Children         []*CmisObject   `gorm:"many2many:filing;association_jointable_foreignkey:object_id;jointable_foreignkey:parent_id;"` // m:n - has many
	Parents          []*CmisObject   `gorm:"many2many:filing;association_jointable_foreignkey:parent_id;jointable_foreignkey:object_id;"` // m:n - has many
}

// CmisProperty Dao model
type CmisProperty struct {
	ID                   uint `gorm:"primary_key;AUTO_INCREMENT"`
	CmisObjectID         uint
	PropertyDefinition   *PropertyDefinition // 1:1 - belongs to
	PropertyDefinitionID uint
	Value                string
}
