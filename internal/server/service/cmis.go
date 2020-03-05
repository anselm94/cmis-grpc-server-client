package service

import (
	"context"
	"docserverclient/internal/server"
	"docserverclient/internal/server/model"
	cmisproto "docserverclient/proto"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/dchest/uniuri"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

// Cmis is a gRPC server implementation
type Cmis struct {
	cmisproto.UnimplementedCmisServiceServer
	DB *gorm.DB
}

// GetRepository callback for Unary call to get the default repository
func (c *Cmis) GetRepository(ctx context.Context, req *empty.Empty) (*cmisproto.Repository, error) {
	log.Printf("Request to fetch the default repository")
	repository := model.Repository{
		ID:              1, // default Repository ID
		RootFolder:      &model.CmisObject{},
		TypeDefinitions: []*model.TypeDefinition{},
	}
	// load the Repository data
	if err := c.DB.First(&repository).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	// Load the root folder data
	if err := c.DB.Model(&repository).Related(&repository.RootFolder, "RootFolder").Error; err != nil {
		log.Print(err)
		return nil, err
	}
	// Load the properties for root folder
	if err := c.DB.Preload("Properties").Find(&repository.RootFolder).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	// Load the type definitions for the repository
	if err := c.DB.Model(&repository).Related(&repository.TypeDefinitions, "TypeDefinitions").Error; err != nil {
		log.Print(err)
		return nil, err
	}
	// Load the corresponding property definitions for type definitions loaded
	if err := c.DB.Preload("PropertyDefinitions").Find(&repository.TypeDefinitions).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	repositoryProto := server.ConvertRepositoryDaoToProto(&repository)
	return repositoryProto, nil
}

// GetObject callback for Unary call (to support CMIS Server's getObject) to fetch the Object
func (c *Cmis) GetObject(ctx context.Context, objectID *cmisproto.CmisObjectId) (*cmisproto.CmisObject, error) {
	log.Printf("Request to get the Object with ID -> %d", objectID.Id)
	cmisObject, err := c.getObject(objectID)
	if err != nil {
		return nil, err
	}
	return cmisObject, nil
}

// SubscribeObject callback for Bidirectional RPC streaming of data between client and server
// Holds the state of the client i.e. the ObjectID of the folder it is in
// * Client continuously streams the ObjectID of the folder it navigates in
// * Server sends back the list of folder/documents in the event of ObjectID requisition or
//   any folder/document created/deleted
// TODO - Notify the client only if any of the folder/document creation/deletion happens
//        in the folder the client is listening to
func (c *Cmis) SubscribeObject(srv cmisproto.CmisService_SubscribeObjectServer) error {
	var cmisObjectID *cmisproto.CmisObjectId // ObjectID of the folder, the client is in
	dbCallback := &DBCallback{
		c:            c,
		cmisObjectID: cmisObjectID,
		srv:          srv,
	}

	// Create a soft-trigger (not a DB trigger) for both Create/Delete operation
	// Will be automatically released, if the bidirectional streaming is end from the client
	createCallbackID := uniuri.New()
	deleteCallbackID := uniuri.New()
	c.DB.Callback().Create().Register(createCallbackID, dbCallback.OnTableUpdated)
	c.DB.Callback().Delete().Register(deleteCallbackID, dbCallback.OnTableUpdated)
	defer c.DB.Callback().Create().Remove(createCallbackID)
	defer c.DB.Callback().Delete().Remove(deleteCallbackID)

	for {
		objectID, err := srv.Recv() // Client streams the ObjectID continuosly
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		log.Printf("Request for a CmisObject with ID - \"%d\"", objectID.Id)
		// Hold the state of the client (i.e. the folder's ObjectID), until the connection is terminated
		cmisObjectID = objectID
		dbCallback.cmisObjectID = objectID
		if cmisObject, err := c.getObject(cmisObjectID); err == nil {
			srv.Send(cmisObject)
		} else {
			log.Print(err)
			return err
		}
	}
}

// getObject retrieves the CmisObject from the DB and converts to Proto
func (c *Cmis) getObject(cmisObjectID *cmisproto.CmisObjectId) (*cmisproto.CmisObject, error) {
	objectID := uint(cmisObjectID.GetId())
	cmisObject := model.CmisObject{
		ID:             objectID,
		Properties:     []*model.CmisProperty{},
		TypeDefinition: &model.TypeDefinition{},
		Children:       []*model.CmisObject{},
		Parents:        []*model.CmisObject{},
	}
	if err := c.DB.Find(&cmisObject).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	if err := c.DB.Model(&cmisObject).Related(&cmisObject.TypeDefinition, "TypeDefinition").Error; err != nil {
		log.Print(err)
		return nil, err
	}
	if err := c.DB.Preload("Properties").Preload("Properties.PropertyDefinition").Find(&cmisObject).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	if err := c.DB.Preload("Parents").Preload("Parents.Properties").Preload("Parents.Properties.PropertyDefinition").Preload("Children").Preload("Children.TypeDefinition").Preload("Children.Properties").Preload("Children.Properties.PropertyDefinition").Find(&cmisObject).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	objectProto := server.ConvertCmisObjectDaoToProto(&cmisObject, true)
	return objectProto, nil
}

// CreateObject callback for Unary call to create a document/folder/any type based on the typeID and name
func (c *Cmis) CreateObject(ctx context.Context, createReq *cmisproto.CreateObjectReq) (*cmisproto.CmisObject, error) {
	log.Printf("Request to create an object with name \"%s\" and type \"%s\"", createReq.Name, createReq.Type)
	typeDef := model.TypeDefinition{
		PropertyDefinitions: []*model.PropertyDefinition{},
	}
	parentCmisObject := model.CmisObject{}

	// Load the Type Definitions and the property definitions
	if err := c.DB.Preload("PropertyDefinitions").Where("name = ?", createReq.Type).First(&typeDef).Error; err != nil {
		return nil, err
	}
	// Load the parent object, so that the new object can be attached as a 'filing' relationship
	if err := c.DB.Preload("Properties").Preload("Properties.PropertyDefinition").Find(&parentCmisObject, uint(createReq.ParentId.Id)).Error; err != nil {
		return nil, err
	}

	// Assemble the CMIS object to be created
	cmisObject := model.CmisObject{
		TypeDefinitionID: typeDef.ID,
		Properties:       make([]*model.CmisProperty, 0),
		Parents: []*model.CmisObject{
			&parentCmisObject,
		},
		RepositoryID: uint(createReq.RepositoryId),
	}

	var objectIDPropDefID uint // temp variable to reference the ID of Property Definition of ObjectId in Table
	for _, propertyDefinition := range typeDef.PropertyDefinitions {
		cmisObjectProperty := model.CmisProperty{
			PropertyDefinitionID: propertyDefinition.ID,
		}
		switch propertyDefinition.Name {
		case "cmis:name":
			cmisObjectProperty.Value = createReq.Name
		case "cmis:parentId":
			cmisObjectProperty.Value = strconv.Itoa(int(createReq.ParentId.Id))
		case "cmis:objectId":
			objectIDPropDefID = propertyDefinition.ID
			continue // Don't add this property, and will be added later
		case "cmis:baseTypeId":
			cmisObjectProperty.Value = createReq.Type
		case "cmis:objectTypeId":
			cmisObjectProperty.Value = createReq.Type
		case "cmis:createdBy":
			cmisObjectProperty.Value = "default"
		case "cmis:lastModifiedBy":
			cmisObjectProperty.Value = "default"
		case "cmis:creationDate":
			cmisObjectProperty.Value = strconv.Itoa(int(time.Now().Unix()))
		case "cmis:lastModificationDate":
			cmisObjectProperty.Value = strconv.Itoa(int(time.Now().Unix()))
		case "cmis:path":
			parentCmisPath := server.GetProperty(&parentCmisObject, "cmis:path")
			cmisObjectProperty.Value = fmt.Sprintf("%s/%s", parentCmisPath, createReq.Name)
		}
		cmisObject.Properties = append(cmisObject.Properties, &cmisObjectProperty)
	}

	// Create the object in DB
	if err := c.DB.Create(&cmisObject).Error; err != nil {
		return nil, err
	}

	// Add the CMIS ObjectId property, as the object Id is not known earlier
	// TODO: Don't use automatic ID creation in table in future design
	if err := c.DB.Model(&cmisObject).Association("Properties").Append(model.CmisProperty{
		PropertyDefinitionID: objectIDPropDefID,
		Value:                strconv.Itoa(int(cmisObject.ID)),
	}).Error; err != nil {
		return nil, err
	}

	// Load the current object with properties and property definitions
	if err := c.DB.Preload("Properties").Preload("Properties.PropertyDefinition").Find(&cmisObject).Error; err != nil {
		return nil, err
	}
	cmisObjectProto := server.ConvertCmisObjectDaoToProto(&cmisObject, false)
	return cmisObjectProto, nil
}

// DeleteObject callback for Unary RPC request to delete an Object
// Deletes the `CmisObject`, along with its associated `CmisProperties` and delete the `filing` (parent-child) relationship
// TODO - Delete tree to be implemented. Currently deleting a folder, orphans its kids (you know, how much you owe then!? in cleaning the DB)
func (c *Cmis) DeleteObject(ctx context.Context, objectID *cmisproto.CmisObjectId) (*cmisproto.CmisObject, error) {
	log.Printf("Request to delete the object with ID \"%d\"", objectID.Id)
	cmisObject := &model.CmisObject{
		ID:         uint(objectID.Id),
		Properties: []*model.CmisProperty{},
		Parents:    []*model.CmisObject{},
	}
	// Load the object to be deleted
	if err := c.DB.Preload("Properties").Preload("Parents").Find(&cmisObject).Error; err != nil {
		return nil, err
	}
	// Cut ties with the parent by deleting the relationship
	if err := c.DB.Model(&cmisObject).Association("Parents").Clear().Error; err != nil {
		return nil, err
	}
	// Slash the lone object and its properties
	if err := c.DB.Model(&cmisObject).Delete(cmisObject).Error; err != nil {
		return nil, err
	}
	cmisObjectProto := server.ConvertCmisObjectDaoToProto(cmisObject, false)
	return cmisObjectProto, nil
}

// DBCallback holds the callback for each of the bidirectional connection,
// and will notify the client of any creation/deletion of records in CmisObject table
// The reference `cmisObjectID` holds the objectID of the folder, the client is in
type DBCallback struct {
	c            *Cmis
	cmisObjectID *cmisproto.CmisObjectId
	srv          cmisproto.CmisService_SubscribeObjectServer
}

// OnTableUpdated will be called for every operation on DB's Create/Delete queries
func (dc *DBCallback) OnTableUpdated(scope *gorm.Scope) {
	if scope.TableName() == "cmis_objects" { // Reject any updates not happening in CmisObjects table
		log.Printf("Updating the client for the folder ID \"%d\"", dc.cmisObjectID.Id)
		cmisObject, _ := dc.c.getObject(dc.cmisObjectID)
		dc.srv.Send(cmisObject) // Update the client of updates in it's current folder
	}
}
