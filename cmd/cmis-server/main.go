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
	"strconv"
	"strings"
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
	router.HandleFunc("/browser/{repositoryID}/root", browserObject)
	router.NotFoundHandler = http.HandlerFunc(browserNotFound)
	log.Fatalf("Error running server -> %s", http.ListenAndServe(config.CmisAppPort, router))
}

// ***** HTTP helpers *****

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

// ***** CMIS handlers *****

func browserRepositoryInfos(w http.ResponseWriter, r *http.Request) {
	repositoryIDs := []string{"1"}
	repositoryInfos := map[string]cmismodel.Repository{}
	log.Println("Request: Get Repository Infos")
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
		includePropertyDefinitions := r.URL.Query().Get("includePropertyDefinitions") != "false"
		log.Printf("Request: Repository => Repository ID: %s \t\t\t CMIS Selector: %s", repositoryID, cmisSelector)
		switch cmisSelector {
		case "typeChildren":
			typeChildren, _ := getTypeChildren(repositoryID, includePropertyDefinitions)
			writeJSON(w, typeChildren)
			return
		case "typeDescendants":
			typeID, _ := url.QueryUnescape(r.URL.Query().Get("typeId"))
			typeDescendants, _ := getTypeDescendants(repositoryID, typeID, includePropertyDefinitions)
			writeJSON(w, typeDescendants)
			return
		case "typeDefinition":
			typeID, _ := url.QueryUnescape(r.URL.Query().Get("typeId"))
			typeDefinition, _ := getTypeDefinition(repositoryID, typeID, includePropertyDefinitions)
			writeJSON(w, typeDefinition)
			return
		default:
			writeError(w, "No selector")
		}
	} else {
		log.Printf("Request: Unknown Path -> %s", r.URL.Path)
		writeError(w, "Repository ID not found")
	}
}

func browserObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repositoryID, _ := vars["repositoryID"]
	objectID, _ := url.QueryUnescape(r.URL.Query().Get("objectId"))
	if r.Method == http.MethodGet {
		cmisSelector := r.URL.Query().Get("cmisselector")
		isSuccinctProperties := r.URL.Query().Get("succinct") != "false"
		includeACL := r.URL.Query().Get("includeACL") == "true"
		includeAllowableActions := r.URL.Query().Get("includeAllowableActions") == "true"
		log.Printf("Request: Object     => Repository ID: %s \t Object ID:%s \t CMIS Selector: %s", repositoryID, objectID, cmisSelector)
		switch cmisSelector {
		case "object":
			cmisObject, _ := getObject(repositoryID, objectID, isSuccinctProperties, includeAllowableActions, includeACL)
			writeJSON(w, cmisObject)
			return
		case "children":
			cmisChildren, _ := getChildren(repositoryID, objectID, isSuccinctProperties, includeAllowableActions, includeACL)
			writeJSON(w, cmisChildren)
			return
		case "parents":
			cmisObjectParents, _ := getParentObjects(repositoryID, objectID, isSuccinctProperties, includeAllowableActions, includeACL)
			writeJSON(w, cmisObjectParents)
			return
		case "parent":
			cmisObjectParent, _ := getParentObject(repositoryID, objectID, isSuccinctProperties, includeAllowableActions, includeACL)
			writeJSON(w, cmisObjectParent)
			return
		default:
			log.Printf("Request: Unknown CMIS Selector -> %s", r.URL.Path)
			writeNotFound(w, "Object not found")
		}
	} else if r.Method == http.MethodPost {
		r.ParseMultipartForm(1024) // Load upto 1KB of data
		cmisAction := r.PostFormValue("cmisaction")
		log.Printf("Request: Object     => Repository ID: %s \t Object ID:%s \t CMIS Action  : %s", repositoryID, objectID, cmisAction)
		cmisPropertyMap := make(map[string]string)
		for key, propertyKey := range r.PostForm {
			if strings.Contains(key, "propertyId") {
				propertyIndexStr := strings.FieldsFunc(key, func(c rune) bool {
					return c == '[' || c == ']'
				})[1]
				propertyIndex, _ := strconv.Atoi(propertyIndexStr)
				propertyValue := r.PostFormValue(fmt.Sprintf("propertyValue[%d]", propertyIndex))
				cmisPropertyMap[propertyKey[0]] = propertyValue
			}
		}
		switch cmisAction {
		case "createFolder":
			cmisObject, _ := createObject(repositoryID, objectID, cmisPropertyMap)
			writeJSON(w, cmisObject)
		case "createDocument":
			cmisObject, _ := createObject(repositoryID, objectID, cmisPropertyMap)
			writeJSON(w, cmisObject)
		case "deleteTree":
			// Won't delete the tree. Just deletes the folder object
			cmisObject, _ := deleteObject(repositoryID, objectID)
			writeJSON(w, cmisObject)
		case "delete":
			cmisObject, _ := deleteObject(repositoryID, objectID)
			writeJSON(w, cmisObject)
		default:
			log.Printf("Request: Unknown CMIS Action -> %s", r.URL.Path)
			writeNotFound(w, "Unknown action")
		}
	} else {
		log.Printf("Request: Method - %s not supported -> %s", r.Method, r.URL.Path)
		writeError(w, "Method not supported")
	}
}

func browserNotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: Unknown path -> %s", r.URL.Path)
	writeNotFound(w, "Not found")
}

// ***** CMIS helpers *****
// These helpers indirectly maps to CMIS services

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

func getTypeDescendants(repositoryID string, typeID string, includePropertyDefinitions bool) ([]*cmismodel.TypeDescendant, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo, err := cmisClient.GetRepository(ctxt, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	if typeID == "" {
		typeDescendants := make([]*cmismodel.TypeDescendant, len(repo.TypeDefinitions))
		for index, typeDefinitionProto := range repo.TypeDefinitions {
			typeDescendants[index] = &cmismodel.TypeDescendant{
				Type: cmisserver.ConvertTypeDefinitionProtoToCmis(typeDefinitionProto, includePropertyDefinitions),
			}
		}
		return typeDescendants, nil
	}
	typeDescendants := make([]*cmismodel.TypeDescendant, 0)
	return typeDescendants, nil
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

func getObject(repositoryID string, objectID string, isSuccinctProperties bool, includeAllowableActions bool, includeACL bool) (*cmismodel.CmisObject, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmisObjectID, _ := strconv.Atoi(objectID)
	cmisObjectProto, err := cmisClient.GetObject(ctxt, &cmisproto.CmisObjectId{Id: int32(cmisObjectID)})
	if err != nil {
		return nil, err
	}
	return cmisserver.ConvertCmisObjectProtoToCmis(cmisObjectProto, isSuccinctProperties, includeAllowableActions, includeACL), nil
}

func getChildren(repositoryID string, objectID string, isSuccinctProperties bool, includeAllowableActions bool, includeACL bool) (*cmismodel.CmisChildren, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmisObjectID, _ := strconv.Atoi(objectID)
	cmisObjectProto, err := cmisClient.GetObject(ctxt, &cmisproto.CmisObjectId{Id: int32(cmisObjectID)})
	if err != nil {
		return nil, err
	}
	return cmisserver.ConvertCmisChildrenProtoToCmis(cmisObjectProto.Children, isSuccinctProperties, includeAllowableActions, includeACL), nil
}

func getParentObject(repositoryID string, objectID string, isSuccinctProperties bool, includeAllowableActions bool, includeACL bool) (*cmismodel.CmisObject, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmisObjectID, _ := strconv.Atoi(objectID)
	cmisObjectProto, err := cmisClient.GetObject(ctxt, &cmisproto.CmisObjectId{Id: int32(cmisObjectID)})
	if err != nil {
		return nil, err
	}
	return cmisserver.ConvertCmisObjectProtoToCmis(cmisObjectProto.Parents[0], isSuccinctProperties, includeAllowableActions, includeACL), nil
}

func getParentObjects(repositoryID string, objectID string, isSuccinctProperties bool, includeAllowableActions bool, includeACL bool) ([]*cmismodel.CmisObjectRelated, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmisObjectID, _ := strconv.Atoi(objectID)
	cmisObjectProto, err := cmisClient.GetObject(ctxt, &cmisproto.CmisObjectId{Id: int32(cmisObjectID)})
	if err != nil {
		return nil, err
	}
	return cmisserver.ConvertCmisParentProtoToCmis(cmisObjectProto.Parents, isSuccinctProperties, includeAllowableActions, includeACL), nil
}

func createObject(repositoryID string, objectID string, propertyMap map[string]string) (*cmismodel.CmisObject, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmisObjectID, _ := strconv.Atoi(objectID)
	repoID, _ := strconv.Atoi(repositoryID)
	cmisObjectProto, err := cmisClient.CreateObject(ctxt, &cmisproto.CreateObjectReq{
		Name:         propertyMap["cmis:name"],
		Type:         propertyMap["cmis:objectTypeId"],
		ParentId:     &cmisproto.CmisObjectId{Id: int32(cmisObjectID)},
		RepositoryId: int32(repoID),
	})
	if err != nil {
		return nil, err
	}
	return cmisserver.ConvertCmisObjectProtoToCmis(cmisObjectProto, true, false, false), nil
}

func deleteObject(repositoryID string, objectID string) (*cmismodel.CmisObject, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmisObjectID, _ := strconv.Atoi(objectID)
	cmisObjectProto, err := cmisClient.DeleteObject(ctxt, &cmisproto.CmisObjectId{Id: int32(cmisObjectID)})
	if err != nil {
		return nil, err
	}
	return cmisserver.ConvertCmisObjectProtoToCmis(cmisObjectProto, true, false, false), nil
}
