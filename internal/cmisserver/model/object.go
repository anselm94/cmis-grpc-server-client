package cmismodel

// Exception

// CmisException holds all the CMIS exceptions
type CmisException struct {
	Message   string `json:"message"`
	Exception string `json:"exception"`
}

// ******** Repository ********

// Repository holds the CMIS Repository info
type Repository struct {
	RepositoryID          string                       `json:"repositoryId"`
	RepositoryName        string                       `json:"repositoryName"`
	RepositoryDescription string                       `json:"repositoryDescription"`
	VendorName            string                       `json:"vendorName"`
	ProductName           string                       `json:"productName"`
	ProductVersion        string                       `json:"productVersion"`
	RootFolderID          string                       `json:"rootFolderId"`
	Capabilities          *RepositoryCapabilities      `json:"capabilities"`
	ACLCapabilities       *ACLCapabilities             `json:"aclCapabilities"`
	LatestChangeLogToken  *string                      `json:"latestChangeLogToken"`
	CmisVersionSupported  string                       `json:"cmisVersionSupported"`
	ThinClientURI         string                       `json:"thinClientURI"`
	ChangesIncomplete     bool                         `json:"changesIncomplete"`
	ChangesOnType         []string                     `json:"changesOnType"`
	ExtendedFeatures      []*RepositoryExtendedFeature `json:"extendedFeatures"`
	RepositoryURL         string                       `json:"repositoryUrl"`
	RootFolderURL         string                       `json:"rootFolderUrl"`
}

// RepositoryCapabilities holds the capabilities of one CMIS Repository
type RepositoryCapabilities struct {
	CapabilityContentStreamUpdatability string              `json:"capabilityContentStreamUpdatability"`
	CapabilityChanges                   string              `json:"capabilityChanges"`
	CapabilityRenditions                string              `json:"capabilityRenditions"`
	CapabilityGetDescendants            bool                `json:"capabilityGetDescendants"`
	CapabilityGetFolderTree             bool                `json:"capabilityGetFolderTree"`
	CapabilityMultifiling               bool                `json:"capabilityMultifiling"`
	CapabilityUnfiling                  bool                `json:"capabilityUnfiling"`
	CapabilityVersionSpecificFiling     bool                `json:"capabilityVersionSpecificFiling"`
	CapabilityPWCSearchable             bool                `json:"capabilityPWCSearchable"`
	CapabilityPWCUpdatable              bool                `json:"capabilityPWCUpdatable"`
	CapabilityAllVersionsSearchable     bool                `json:"capabilityAllVersionsSearchable"`
	CapabilityOrderBy                   string              `json:"capabilityOrderBy"`
	CapabilityQuery                     string              `json:"capabilityQuery"`
	CapabilityJoin                      string              `json:"capabilityJoin"`
	CapabilityACL                       string              `json:"capabilityACL"`
	CapabilityCreatablePropertyTypes    map[string][]string `json:"capabilityCreatablePropertyTypes"`
	CapabilityNewTypeSettableAttributes map[string]bool     `json:"capabilityNewTypeSettableAttributes"`
}

// ACLCapabilities holds the list of ACLs & Permissions associated with it
type ACLCapabilities struct {
	SupportedPermissions string           `json:"supportedPermissions"`
	Propagation          string           `json:"propagation"`
	Permissions          []*Permission    `json:"permissions"`
	PermissionMapping    []*PermissionMap `json:"permissionMapping"`
}

// Permission is a representation of CMIS Permission
type Permission struct {
	Permission  string `json:"permission"`
	Description string `json:"description"`
}

// PermissionMap holds the map of permission and actions they can do
type PermissionMap struct {
	Key        string   `json:"key"`
	Permission []string `json:"permission"`
}

// RepositoryExtendedFeature holds the extended feature not supported by CMIS natively
type RepositoryExtendedFeature struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	CommonName   string `json:"commonName"`
	VersionLabel string `json:"versionLabel"`
	Description  string `json:"description"`
}

// ******** Types ********

// TypeDefinition holds a CMIS Type Definition. A TypeDefinition can hold multiple PropertyDefinitions
type TypeDefinition struct {
	ID                       string                         `json:"id"`
	LocalName                string                         `json:"localName"`
	LocalNamespace           string                         `json:"localNamespace"`
	DisplayName              string                         `json:"displayName"`
	QueryName                string                         `json:"queryName"`
	Description              string                         `json:"description"`
	BaseID                   string                         `json:"baseId"`
	Creatable                bool                           `json:"creatable"`
	Fileable                 bool                           `json:"fileable"`
	Queryable                bool                           `json:"queryable"`
	FulltextIndexed          bool                           `json:"fulltextIndexed"`
	IncludedInSupertypeQuery bool                           `json:"includedInSupertypeQuery"`
	ControllablePolicy       bool                           `json:"controllablePolicy"`
	ControllableACL          bool                           `json:"controllableACL"`
	TypeMutability           map[string]bool                `json:"typeMutability"`
	Versionable              *bool                          `json:"versionable,omitempty"`
	ContentStreamAllowed     *string                        `json:"contentStreamAllowed,omitempty"`
	PropertyDefinitions      map[string]*PropertyDefinition `json:"propertyDefinitions"`
}

// PropertyDefinition holds a Property Definition
type PropertyDefinition struct {
	ID            string `json:"id"`
	LocalName     string `json:"localName"`
	DisplayName   string `json:"displayName"`
	QueryName     string `json:"queryName"`
	Description   string `json:"description"`
	PropertyType  string `json:"propertyType"`
	Cardinality   string `json:"cardinality"`
	Updateability string `json:"updatability"`
	Inherited     bool   `json:"inherited"`
	Required      bool   `json:"required"`
	Queryable     bool   `json:"queryable"`
	Orderable     bool   `json:"orderable"`
}

// TypeChildren is a representation of all Types for Browser binding
type TypeChildren struct {
	Types        []*TypeDefinition `json:"types"`
	HasMoreItems bool              `json:"hasMoreItems"`
	NumItems     int               `json:"numItems"`
}

// TypeDescendant is a representation of a TypeDefinition as a descendant
type TypeDescendant struct {
	Type *TypeDefinition `json:"type"`
}

// ******** Objects ********

// AllowableActions holds all the actions allowable for a CMIS Object
type AllowableActions struct {
	CanDeleteObject           bool `json:"canDeleteObject"`
	CanUpdateProperties       bool `json:"canUpdateProperties"`
	CanGetFolderTree          bool `json:"canGetFolderTree"`
	CanGetProperties          bool `json:"canGetProperties"`
	CanGetObjectRelationships bool `json:"canGetObjectRelationships"`
	CanGetObjectParents       bool `json:"canGetObjectParents"`
	CanGetFolderParent        bool `json:"canGetFolderParent"`
	CanGetDescendants         bool `json:"canGetDescendants"`
	CanMoveObject             bool `json:"canMoveObject"`
	CanDeleteContentStream    bool `json:"canDeleteContentStream"`
	CanCheckOut               bool `json:"canCheckOut"`
	CanCancelCheckOut         bool `json:"canCancelCheckOut"`
	CanCheckIn                bool `json:"canCheckIn"`
	CanSetContentStream       bool `json:"canSetContentStream"`
	CanGetAllVersions         bool `json:"canGetAllVersions"`
	CanAddObjectToFolder      bool `json:"canAddObjectToFolder"`
	CanRemoveObjectFromFolder bool `json:"canRemoveObjectFromFolder"`
	CanGetContentStream       bool `json:"canGetContentStream"`
	CanApplyPolicy            bool `json:"canApplyPolicy"`
	CanGetAppliedPolicies     bool `json:"canGetAppliedPolicies"`
	CanRemovePolicy           bool `json:"canRemovePolicy"`
	CanGetChildren            bool `json:"canGetChildren"`
	CanCreateDocument         bool `json:"canCreateDocument"`
	CanCreateFolder           bool `json:"canCreateFolder"`
	CanCreateRelationship     bool `json:"canCreateRelationship"`
	CanCreateItem             bool `json:"canCreateItem"`
	CanDeleteTree             bool `json:"canDeleteTree"`
	CanGetRenditions          bool `json:"canGetRenditions"`
	CanGetACL                 bool `json:"canGetACL"`
	CanApplyACL               bool `json:"canApplyACL"`
}

// ACL is the Access Control List containing ACEs
type ACL struct {
	ACEs []*ACE `json:"aces"`
}

// Principal holds the Principal's ID, which can be direct/indirect
type Principal struct {
	PrincipalID string `json:"principal"`
}

// ACE is the CMIS Access Control Entry, defining what Permissions a Principal have
type ACE struct {
	Principal   *Principal `json:"principal"`
	Permissions []string   `json:"permissions"`
	IsDirect    bool       `json:"isDirect"`
}

// CmisProperty holds the representation for a CMIS Property
type CmisProperty struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Cardinality string      `json:"cardinality"`
	Value       interface{} `json:"value"`
}

// CmisObject holds the representation for a CMIS Object
type CmisObject struct {
	Properties         *[]*CmisProperty       `json:"properties,omitempty"`
	SuccinctProperties map[string]interface{} `json:"succinctProperties,omitempty"`
	AllowableActions   *AllowableActions      `json:"allowableActions"`
	ACL                *ACL                   `json:"acl"`
	ExactACL           *bool                  `json:"exactACL"`
}

// CmisObjectRelated holds a CmisObject, as is required only for getChildren calls
type CmisObjectRelated struct {
	Object *CmisObject `json:"object"`
}

// CmisChildren holds the list of Child CmisObjects
type CmisChildren struct {
	Objects      []*CmisObjectRelated `json:"objects"`
	HasMoreItems bool                 `json:"hasMoreItems"`
	NumItems     int                  `json:"numItems"`
}
