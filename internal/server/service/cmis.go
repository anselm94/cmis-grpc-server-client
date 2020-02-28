package service

import (
	"context"
	"docserverclient/internal/server"
	"docserverclient/internal/server/model"
	cmisproto "docserverclient/proto"
	"fmt"
	"io"
	"log"

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

func (c *Cmis) GetObject(ctx context.Context, objectID *cmisproto.CmisObjectId) (*CmisObject, error) {

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
	if err := c.DB.Preload("Parents").Preload("Children").Preload("Children.TypeDefinition").Preload("Children.Properties").Preload("Children.Properties.PropertyDefinition").Find(&cmisObject).Error; err != nil {
		log.Print(err)
		return nil, err
	}
	objectProto := server.ConvertCmisObjectDaoToProto(&cmisObject, true)
	return objectProto, nil
}

// CreateObject callback for Unary call to create a document/folder/any type based on the typeID and name
func (c *Cmis) CreateObject(ctx context.Context, createReq *cmisproto.CreateObjectReq) (*empty.Empty, error) {
	log.Printf("Request to create an object with name \"%s\" and type \"%s\"", createReq.Name, createReq.Type)
	// Reference Name property
	namePropDef := model.PropertyDefinition{
		TypeDefinition: &model.TypeDefinition{
			Name: createReq.Type,
		},
	}
	// Reference PropertyID property
	parentIDPropDef := model.PropertyDefinition{
		TypeDefinition: &model.TypeDefinition{
			Name: createReq.Type,
		},
	}
	typeDef := model.TypeDefinition{}
	parentCmisObject := model.CmisObject{}

	// Load the Type Definitions and the property definitions
	if err := c.DB.Preload("TypeDefinition").Where("name = ?", "cmis:name").First(&namePropDef).Error; err != nil {
		return &empty.Empty{}, err
	}
	if err := c.DB.Preload("TypeDefinition").Where("name = ?", "cmis:parentId").First(&parentIDPropDef).Error; err != nil {
		return &empty.Empty{}, err
	}
	if err := c.DB.Where("name = ?", createReq.Type).First(&typeDef).Error; err != nil {
		return &empty.Empty{}, err
	}
	// Load the parent object, so that the new object can be attached as a 'filing' relationship
	if err := c.DB.Find(&parentCmisObject, uint(createReq.ParentId.Id)).Error; err != nil {
		return &empty.Empty{}, err
	}

	// Assemble the CMIS object to be created
	cmisObject := model.CmisObject{
		Properties: []*model.CmisProperty{
			&model.CmisProperty{
				Value:                createReq.Name,
				PropertyDefinitionID: namePropDef.ID,
			}, &model.CmisProperty{
				Value:                fmt.Sprint(createReq.ParentId.Id),
				PropertyDefinitionID: parentIDPropDef.ID,
			},
		},
		TypeDefinitionID: typeDef.ID,
		Parents: []*model.CmisObject{
			&parentCmisObject,
		},
		RepositoryID: uint(createReq.RepositoryId),
	}
	// Create the object in DB
	if err := c.DB.Create(&cmisObject).Error; err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

// DeleteObject callback for Unary RPC request to delete an Object
// Deletes the `CmisObject`, along with its associated `CmisProperties` and delete the `filing` (parent-child) relationship
// TODO - Delete tree to be implemented. Currently deleting a folder, orphans its kids (you know, how much you owe then!? in cleaning the DB)
func (c *Cmis) DeleteObject(ctx context.Context, objectID *cmisproto.CmisObjectId) (*empty.Empty, error) {
	log.Printf("Request to delete the object with ID \"%d\"", objectID.Id)
	cmisObject := &model.CmisObject{
		ID:         uint(objectID.Id),
		Properties: []*model.CmisProperty{},
		Parents:    []*model.CmisObject{},
	}
	// Load the object to be deleted
	if err := c.DB.Preload("Properties").Preload("Parents").Find(&cmisObject).Error; err != nil {
		return &empty.Empty{}, err
	}
	// Cut ties with the parent by deleting the relationship
	if err := c.DB.Model(&cmisObject).Association("Parents").Clear().Error; err != nil {
		return &empty.Empty{}, err
	}
	// Slash the lone object and its properties
	if err := c.DB.Model(&cmisObject).Delete(cmisObject).Error; err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
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
