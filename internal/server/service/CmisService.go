package service

import (
	"context"
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
	repository := &model.Repository{
		ID: 1,
	}
	rootFolder := &model.CmisObject{}
	typeDefinitions := &[]model.TypeDefinition{}
	propertyDefinitions := &[]model.PropertyDefinition{}
	c.DB.Model(repository).Related(rootFolder, "RootFolder")
	c.DB.Model(repository).Related(typeDefinitions, "TypeDefinitions")
	c.DB.Model(typeDefinitions).Related(propertyDefinitions)

	typedefsResp := []*cmis.TypeDefinition{}
	for _, typedef := range *typeDefinitions {
		typeDefResp := &cmis.TypeDefinition{
			Name:                typedef.Name,
			Description:         typedef.Description,
			PropertyDefinitions: []*cmis.PropertyDefinition{},
		}
		for _, propdef := range *propertyDefinitions {
			if propdef.TypeDefinitionID == typedef.ID {
				propDefResp := cmis.PropertyDefinition{
					Name:        propdef.Name,
					Description: propdef.Description,
					Datatype:    propdef.Type,
				}
				typeDefResp.PropertyDefinitions = append(typeDefResp.PropertyDefinitions, &propDefResp)
			}
		}
		typedefsResp = append(typedefsResp, typeDefResp)
	}
	repositoryResp := &cmis.Repository{
		Name:            repository.Name,
		Description:     repository.Description,
		TypeDefinitions: typedefsResp,
	}
	return repositoryResp, nil
}

func (c *Cmis) GetObject(req *cmis.CmisObjectId, srv cmis.CmisService_GetObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}

func (c *Cmis) CreateObject(ctx context.Context, req *cmis.CmisObject) (*cmis.CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateObject not implemented")
}

func (c *Cmis) DeleteObject(ctx context.Context, req *cmis.CmisObject) (*cmis.CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteObject not implemented")
}
