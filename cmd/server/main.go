package main

import (
	"docserverclient"
	"fmt"
	"log"
	"net"

	"docserverclient/internal/server/model"
	"docserverclient/internal/server/service"
	cmis "docserverclient/proto"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"google.golang.org/grpc"
)

func main() {
	appConfig := docserverclient.NewDefaultConfig()
	listener, err := net.Listen("tcp", appConfig.AppPort)
	if err != nil {
		log.Fatalf("Failed to listen TCP on %s -> %s", appConfig.AppPort, err)
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", appConfig.DBHost, appConfig.DBPort, appConfig.DBUser, appConfig.DBName, appConfig.DBPassword, appConfig.DBSSLMode))
	if err != nil {
		log.Fatalf("Error setting up DB -> %s", err)
	}
	defer db.Close()
	db.AutoMigrate(&model.Repository{}, &model.TypeDefinition{}, &model.PropertyDefinition{}, &model.CmisObject{}, &model.CmisProperty{})

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	cmis.RegisterCmisServiceServer(grpcServer, &service.Cmis{
		DB: db,
	})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error while serving gRPC -> %s", err)
	}
}
