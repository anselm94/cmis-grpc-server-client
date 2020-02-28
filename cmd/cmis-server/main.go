package main

import (
	"context"
	"docserverclient"
	"docserverclient/internal/cmisserver"
	cmismodel "docserverclient/internal/cmisserver/model"
	cmisproto "docserverclient/proto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var (
	cmisClient cmisproto.CmisServiceClient
)

func main() {
	config := docserverclient.NewDefaultConfig()

	grpcConnection, err := grpc.Dial(fmt.Sprintf("%s%s", config.GrpcAppHost, config.GrpcAppPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Connection could not be established with %s%s -> %s", config.GrpcAppHost, config.GrpcAppPort, err)
	}
	defer grpcConnection.Close()
	cmisClient = cmisproto.NewCmisServiceClient(grpcConnection)

	router := mux.NewRouter()
	router.HandleFunc("/browser", browserRepositoryInfos)
	router.HandleFunc("/browser/{repositoryID}", browserRepository)
	router.HandleFunc("/browser/{repositoryID}/{objectID}", browserObject)
	router.NotFoundHandler = http.HandlerFunc(browserNotFound)
	log.Fatalf("Error running server -> %s", http.ListenAndServe(config.CmisAppPort, router))
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	jsonObject, err := json.Marshal(data)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonObject)
}

func writeError(w http.ResponseWriter, err string) {
	exceptionMsg := cmismodel.CmisException{
		Exception: "notSupported",
		Message:   err,
	}
	jsonObject, _ := json.Marshal(exceptionMsg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonObject)
}

func writeNotFound(w http.ResponseWriter, err string) {
	exceptionMsg := cmismodel.CmisException{
		Exception: "objectNotFound",
		Message:   err,
	}
	jsonObject, _ := json.Marshal(exceptionMsg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(jsonObject)
}

func browserRepositoryInfos(w http.ResponseWriter, r *http.Request) {
	repositoryIDs := []string{"1"}
	repositoryInfos := map[string]cmismodel.Repository{}
	for _, repositoryID := range repositoryIDs {
		repositoryInfo, err := getRepository(repositoryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		repositoryInfos[repositoryID] = *repositoryInfo
	}
	writeJSON(w, repositoryInfos)
}

func browserRepository(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repositoryID, ok := vars["repositoryID"]
	if ok {
		cmisSelector := r.URL.Query().Get("cmisselector")
		switch cmisSelector {
		case "typeChildren":
			includePropertyDefinitions := r.URL.Query().Get("includePropertyDefinitions") == "true"
			typeChildren, _ := getTypeChildren(repositoryID, includePropertyDefinitions)
			writeJSON(w, typeChildren)
			return
		case "typeDefinition":
			includePropertyDefinitions := r.URL.Query().Get("includePropertyDefinitions") == "true"
			typeID, _ := url.QueryUnescape(r.URL.Query().Get("typeId"))
			typeDefinition, _ := getTypeDefinition(repositoryID, typeID, includePropertyDefinitions)
			writeJSON(w, typeDefinition)
			return
		default:
			writeError(w, "No selector")
		}
	} else {
		writeError(w, "Repository ID not found")
	}
}

func browserObject(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// repositoryID, _ := vars["repositoryID"]
	// objectID, ok := vars["objectID"]
	ok := true
	if ok {
		cmisSelector := r.URL.Query().Get("cmisselector")
		switch cmisSelector {
		default:
			writeNotFound(w, "Object not found")
		}
	} else {
		writeError(w, "Object ID not found")
	}
}

func browserNotFound(w http.ResponseWriter, r *http.Request) {
	writeNotFound(w, "Not found")
}

func getRepository(repositoryID string) (*cmismodel.Repository, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo, err := cmisClient.GetRepository(ctxt, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return cmisserver.ConvertRepositoryProtoToCmis(repo), nil
}

func getTypeChildren(repositoryID string, includePropertyDefinitions bool) (*cmismodel.TypeChildren, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo, err := cmisClient.GetRepository(ctxt, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	typeChildren := cmismodel.TypeChildren{
		Types:        cmisserver.ConvertTypeDefinitionsProtoToCmis(repo.TypeDefinitions, includePropertyDefinitions),
		HasMoreItems: false,
		NumItems:     len(repo.TypeDefinitions),
	}
	return &typeChildren, nil
}

func getTypeDefinition(repositoryID string, typeID string, includePropertyDefinitions bool) (*cmismodel.TypeDefinition, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo, err := cmisClient.GetRepository(ctxt, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	var typedefinitionProto *cmisproto.TypeDefinition
	for _, typedef := range repo.TypeDefinitions {
		if typedef.Name == typeID {
			typedefinitionProto = typedef
		}
	}
	return cmisserver.ConvertTypeDefinitionProtoToCmis(typedefinitionProto, includePropertyDefinitions), nil
}
