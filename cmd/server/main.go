package main

import (
	"docserverclient"
	"fmt"
	"log"
	"net"

	"docserverclient/internal/server"
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

	//Uncomment to drop and reinitialise data
	server.DropTables(db)

	isDataInitRequired := !db.HasTable("repositories")
	db.AutoMigrate(&model.Repository{}, &model.TypeDefinition{}, &model.PropertyDefinition{}, &model.CmisObject{}, &model.CmisProperty{})
	db.Model(&model.TypeDefinition{}).AddForeignKey("repository_id", "repositories(id)", "CASCADE", "CASCADE")
	db.Model(&model.PropertyDefinition{}).AddForeignKey("type_definition_id", "type_definitions(id)", "CASCADE", "CASCADE")
	db.Model(&model.CmisObject{}).AddForeignKey("repository_id", "repositories(id)", "CASCADE", "CASCADE")
	db.Model(&model.CmisObject{}).AddForeignKey("type_definition_id", "type_definitions(id)", "RESTRICT", "CASCADE")
	db.Model(&model.CmisProperty{}).AddForeignKey("cmis_object_id", "cmis_objects(id)", "CASCADE", "CASCADE")
	db.Model(&model.CmisProperty{}).AddForeignKey("property_definition_id", "property_definitions(id)", "RESTRICT", "CASCADE")
	if isDataInitRequired {
		server.CreateInitData(db)
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	cmis.RegisterCmisServiceServer(grpcServer, &service.Cmis{
		DB: db,
	})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error while serving gRPC -> %s", err)
	}
}
