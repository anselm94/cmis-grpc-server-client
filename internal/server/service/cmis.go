package service

import (
	"context"
	"docserverclient/internal/server"
	"docserverclient/internal/server/model"
	cmis "docserverclient/proto"
	"fmt"
	"io"

	"github.com/dchest/uniuri"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

type Cmis struct {
	cmis.UnimplementedCmisServiceServer
	DB *gorm.DB
}

func (c *Cmis) GetRepository(ctx context.Context, req *empty.Empty) (*cmis.Repository, error) {
	repository := model.Repository{
		ID:              1,
		RootFolder:      &model.CmisObject{},
		TypeDefinitions: []*model.TypeDefinition{},
	}
	c.DB.First(&repository)
	c.DB.Model(&repository).Related(&repository.RootFolder, "RootFolder")
	c.DB.Preload("Properties").Find(&repository.RootFolder)
	c.DB.Model(&repository).Related(&repository.TypeDefinitions, "TypeDefinitions")
	c.DB.Preload("PropertyDefinitions").Find(&repository.TypeDefinitions)

	repositoryProto := server.ConvertRepositoryDaoToProto(&repository)
	return repositoryProto, nil
}

func (c *Cmis) SubscribeObject(srv cmis.CmisService_SubscribeObjectServer) error {
	var cmisObjectID *cmis.CmisObjectId
	dbCallback := &DBCallback{
		c:            c,
		cmisObjectID: cmisObjectID,
		srv:          srv,
	}
	createCallbackID := uniuri.New()
	deleteCallbackID := uniuri.New()
	c.DB.Callback().Create().Register(createCallbackID, dbCallback.onTableUpdated)
	c.DB.Callback().Delete().Register(deleteCallbackID, dbCallback.onTableUpdated)
	defer c.DB.Callback().Create().Remove(createCallbackID)
	defer c.DB.Callback().Delete().Remove(deleteCallbackID)

	for {
		objectID, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		cmisObjectID = objectID
		dbCallback.cmisObjectID = objectID
		cmisObject, err := c.getObject(cmisObjectID)
		if err != nil {
			return nil
		}
		srv.Send(cmisObject)
	}
}

func (c *Cmis) getObject(cmisObjectID *cmis.CmisObjectId) (*cmis.CmisObject, error) {
	objectID := uint(cmisObjectID.GetId())
	cmisObject := model.CmisObject{
		ID:             objectID,
		Properties:     []*model.CmisProperty{},
		TypeDefinition: &model.TypeDefinition{},
		Children:       []*model.CmisObject{},
		Parents:        []*model.CmisObject{},
	}
	c.DB.Find(&cmisObject)
	c.DB.Model(&cmisObject).Related(&cmisObject.TypeDefinition, "TypeDefinition")
	c.DB.Preload("Parents").Preload("Children").Preload("Children.TypeDefinition").Preload("Children.Properties").Preload("Children.Properties.PropertyDefinition").Find(&cmisObject)
	objectProto := server.ConvertCmisObjectDaoToProto(&cmisObject, true)
	return objectProto, nil
}

type DBCallback struct {
	c            *Cmis
	cmisObjectID *cmis.CmisObjectId
	srv          cmis.CmisService_SubscribeObjectServer
}

func (dc *DBCallback) onTableUpdated(scope *gorm.Scope) {
	if scope.TableName() == "cmis_objects" {
		cmisObject, _ := dc.c.getObject(dc.cmisObjectID)
		dc.srv.Send(cmisObject)
	}
}

func (c *Cmis) CreateObject(ctx context.Context, createReq *cmis.CreateObjectReq) (*empty.Empty, error) {
	namePropDef := model.PropertyDefinition{
		TypeDefinition: &model.TypeDefinition{
			Name: createReq.Type,
		},
	}
	parentIDPropDef := model.PropertyDefinition{
		TypeDefinition: &model.TypeDefinition{
			Name: createReq.Type,
		},
	}
	typeDef := model.TypeDefinition{}
	parentCmisObject := model.CmisObject{}

	c.DB.Preload("TypeDefinition").Where("name = ?", "cmis:name").First(&namePropDef)
	c.DB.Preload("TypeDefinition").Where("name = ?", "cmis:parentId").First(&parentIDPropDef)
	c.DB.Where("name = ?", createReq.Type).First(&typeDef)
	c.DB.Find(&parentCmisObject, uint(createReq.ParentId.Id))

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
	if err := c.DB.Create(&cmisObject).Error; err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func (c *Cmis) DeleteObject(ctx context.Context, objectID *cmis.CmisObjectId) (*empty.Empty, error) {
	cmisObject := &model.CmisObject{
		ID:         uint(objectID.Id),
		Properties: []*model.CmisProperty{},
		Parents:    []*model.CmisObject{},
	}
	// Load the object to be deleted
	c.DB.Preload("Properties").Preload("Parents").Find(&cmisObject)
	// Cut ties with the parent by deleting the relationship
	c.DB.Model(&cmisObject).Association("Parents").Clear()
	// Slash the lone object and its properties
	if err := c.DB.Model(&cmisObject).Delete(cmisObject).Error; err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}
