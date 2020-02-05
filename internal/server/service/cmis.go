package service

import (
	"context"
	"docserverclient/internal/server"
	"docserverclient/internal/server/model"
	cmis "docserverclient/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Cmis struct {
	cmis.UnimplementedCmisServiceServer
	DB *gorm.DB
}

func (c *Cmis) GetRepository(ctx context.Context, req *empty.Empty) (*cmis.Repository, error) {
	repository := model.Repository{
		ID:              1,
		RootFolder:      model.CmisObject{},
		TypeDefinitions: []model.TypeDefinition{},
	}
	c.DB.Find(&repository)
	c.DB.Model(&repository).Related(&repository.RootFolder, "RootFolder")
	c.DB.Preload("Properties").Find(&repository.RootFolder)
	c.DB.Model(&repository).Related(&repository.TypeDefinitions, "TypeDefinitions")
	c.DB.Preload("PropertyDefinitions").Find(&repository.TypeDefinitions)

	repositoryProto := server.ConvertRepositoryDaoToProto(&repository)
	return repositoryProto, nil
}

func (c *Cmis) GetObject(cmisObjectID *cmis.CmisObjectId, stream cmis.CmisService_GetObjectServer) error {
	objectID := uint(cmisObjectID.GetId())
	cmisObject := model.CmisObject{
		ID:             objectID,
		Properties:     []model.CmisProperty{},
		TypeDefinition: model.TypeDefinition{},
		Children:       []*model.CmisObject{},
		Parents:        []*model.CmisObject{},
	}
	c.DB.Find(&cmisObject)
	c.DB.Model(&cmisObject).Related(&cmisObject.TypeDefinition, "TypeDefinition")
	c.DB.Preload("Parents").Preload("Children").Preload("Children.TypeDefinition").Preload("Children.Properties").Preload("Children.Properties.PropertyDefinition").Find(&cmisObject)
	objectProto := server.ConvertCmisObjectDaoToProto(&cmisObject, true)
	stream.Send(objectProto)
	return nil
}

func (c *Cmis) CreateObject(ctx context.Context, cmisObject *cmis.CmisObject) (*cmis.CmisObject, error) {

	return nil, status.Errorf(codes.Unimplemented, "method CreateObject not implemented")
}

func (*Cmis) DeleteObject(ctx context.Context, req *cmis.CmisObjectId) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteObject not implemented")
}
