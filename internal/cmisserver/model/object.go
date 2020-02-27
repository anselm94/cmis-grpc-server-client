package cmismodel

type Repository struct {
	RepositoryID string,
	RepositoryName string,
	RepositoryDescription string,
	VendorName string,
	ProductName string,
	ProductVersion string,
	RootFolderID string,
	Capabilities RepositoryCapabilities,
	AclCapabilities AclCapabilities,
	LatestChangeLogToken string,
	CmisVersionSupported string,
	ThinClientURI string,
	ChangesIncomplete bool,
	ChangesOnType []string,
	ExtendedFeatures []RepositoryExtendedFeature,
	RepositoryUrl string,
	RootFolderUrl string,
}

type RepositoryExtendedFeature struct {
	ID string,
	URL string,
	CommonName string,
	VersionLabel string,
	Description string,
}

type RepositoryCapabilities struct {
	CapabilityContentStreamUpdatability string,
	CapabilityChanges string,
	CapabilityRenditions string,
	CapabilityGetDescendants bool,
	capabilityGetFolderTree bool,
	capabilityMultifiling bool,
	capabilityUnfiling bool,
	capabilityVersionSpecificFiling bool,
	capabilityPWCSearchable bool,
	capabilityPWCUpdatable bool,
	capabilityAllVersionsSearchable bool,
	capabilityOrderBy string,
	capabilityQuery string,
	capabilityJoin string,
	capabilityACL string,
	capabilityCreatablePropertyTypes map[string][]string,
	capabilityNewTypeSettableAttributes: map[string]bool
}

type AclCapabilities struct {
	SupportedPermissions string,
	Propagation string,
	Permissions []Permission
	PermissionMapping []PermissionMap
}

type Permission struct {
	Permission string,
	Description string,
}

type PermissionMap struct {
	Key string,
	Permission []string
}