package main

import (
	"context"
	cmis "docserverclient/proto"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:9999", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Connection could not be established -> %s", err)
	}
	defer connection.Close()

	cmisClient := cmis.NewCmisServiceClient(connection)

	ctxt, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	repository, err := cmisClient.GetRepository(ctxt, nil)
	if err != nil {
		log.Fatalf("Error getting response -> %s", err)
	}
	log.Printf("Greeting: %s", repository.GetName())
}
