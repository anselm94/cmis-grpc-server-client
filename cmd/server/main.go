package main

import (
	"context"
	"log"
	"net"

	cmis "grpc-trial/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CmisServer struct {
	cmis.UnimplementedCmisServiceServer
}

func (*CmisServer) GetRepository(ctx context.Context, req *empty.Empty) (*cmis.Repository, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepository not implemented")
}

func (*CmisServer) GetChildren(ctx context.Context, req *cmis.CmisObject) (*cmis.CmisChildren, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChildren not implemented")
}

func (*CmisServer) GetObject(ctx context.Context, req *cmis.CmisID) (*cmis.CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}

func (*CmisServer) CreateObject(ctx context.Context, req *cmis.CmisObject) (*cmis.CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateObject not implemented")
}

func (*CmisServer) DeleteObject(ctx context.Context, req *cmis.CmisObject) (*cmis.CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteObject not implemented")
}

func main() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("Failed to listen TCP on 9999 -> %s", err)
	}
	grpcServer := grpc.NewServer()
	cmis.RegisterCmisServiceServer(grpcServer, &CmisServer{})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error while serving gRPC -> %s", err)
	}
}
