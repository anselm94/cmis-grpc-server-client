package cmismodel

// Exception

type CmisException struct {
	Message   string `json:"message"`
	Exception string `json:"exception"`
}

// ******** Repository ********

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

type ACLCapabilities struct {
	SupportedPermissions string           `json:"supportedPermissions"`
	Propagation          string           `json:"propagation"`
	Permissions          []*Permission    `json:"permissions"`
	PermissionMapping    []*PermissionMap `json:"permissionMapping"`
}

type Permission struct {
	Permission  string `json:"permission"`
	Description string `json:"description"`
}

type PermissionMap struct {
	Key        string   `json:"key"`
	Permission []string `json:"permission"`
}

type RepositoryExtendedFeature struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	CommonName   string `json:"commonName"`
	VersionLabel string `json:"versionLabel"`
	Description  string `json:"description"`
}

// ******** Types ********

type TypeDefinition struct {
	ID                       string                `json:"id"`
	LocalName                string                `json:"localName"`
	LocalNamespace           string                `json:"localNamespace"`
	DisplayName              string                `json:"displayName"`
	QueryName                string                `json:"queryName"`
	Description              string                `json:"description"`
	BaseID                   string                `json:"baseId"`
	Creatable                bool                  `json:"creatable"`
	Fileable                 bool                  `json:"fileable"`
	Queryable                bool                  `json:"queryable"`
	FulltextIndexed          bool                  `json:"fulltextIndexed"`
	IncludedInSupertypeQuery bool                  `json:"includedInSupertypeQuery"`
	ControllablePolicy       bool                  `json:"controllablePolicy"`
	ControllableACL          bool                  `json:"controllableACL"`
	TypeMutability           map[string]bool       `json:"typeMutability"`
	Versionable              bool                  `json:"versionable"`
	ContentStreamAllowed     string                `json:"contentStreamAllowed"`
	PropertyDefinitions      []*PropertyDefinition `json:"propertyDefinitions"`
}

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

type TypeChildren struct {
	Types        []*TypeDefinition `json:"types"`
	HasMoreItems bool              `json:"hasMoreItems"`
	NumItems     int               `json:"numItems"`
}
