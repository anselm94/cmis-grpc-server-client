package main

import (
	"docserverclient"
	"flag"
	"fmt"
	"log"
	"net"

	"docserverclient/internal/server"
	"docserverclient/internal/server/model"
	"docserverclient/internal/server/service"
	cmisproto "docserverclient/proto"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
)

func main() {
	shouldCreateInitData := flag.Bool("populate", false, "Populates some data for initial run. Creates a repository, typedefinitions for folder & document, name & parentId propertydefinitions for each of typedefinitions and a folder & a document in the root folder")
	flag.Parse()

	appConfig := docserverclient.NewDefaultConfig()
	listener, err := net.Listen("tcp", appConfig.GrpcAppPort)
	if err != nil {
		log.Fatalf("Failed to listen TCP on %s -> %s", appConfig.GrpcAppPort, err)
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", appConfig.DBHost, appConfig.DBPort, appConfig.DBUser, appConfig.DBName, appConfig.DBPassword, appConfig.DBSSLMode))
	if err != nil {
		log.Fatalf("Error setting up DB -> %s", err)
	}
	defer db.Close()

	if *shouldCreateInitData {
		log.Println("Data population requested")
		log.Println("Dropping tables if already exists...")
		server.DropTables(db)
	}
	log.Println("Migrating tables...")
	db.AutoMigrate(&model.Repository{}, &model.TypeDefinition{}, &model.PropertyDefinition{}, &model.CmisObject{}, &model.CmisProperty{})
	log.Println("Migration complete")
	if *shouldCreateInitData {
		log.Println("Creating DB constraints...")
		db.Model(&model.TypeDefinition{}).AddForeignKey("repository_id", "repositories(id)", "CASCADE", "CASCADE")
		db.Model(&model.PropertyDefinition{}).AddForeignKey("type_definition_id", "type_definitions(id)", "CASCADE", "CASCADE")
		db.Model(&model.CmisObject{}).AddForeignKey("repository_id", "repositories(id)", "CASCADE", "CASCADE")
		db.Model(&model.CmisObject{}).AddForeignKey("type_definition_id", "type_definitions(id)", "RESTRICT", "CASCADE")
		db.Model(&model.CmisProperty{}).AddForeignKey("cmis_object_id", "cmis_objects(id)", "CASCADE", "CASCADE")
		db.Model(&model.CmisProperty{}).AddForeignKey("property_definition_id", "property_definitions(id)", "RESTRICT", "CASCADE")
		server.CreateInitData(db)
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	cmisproto.RegisterCmisServiceServer(grpcServer, &service.Cmis{
		DB: db,
	})
	log.Printf("Listening to gRPC requests at %s:%s", appConfig.GrpcAppHost, appConfig.GrpcAppPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error while serving gRPC -> %s", err)
	}
}
