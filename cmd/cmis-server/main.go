package main

import (
	"docserverclient"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config := docserverclient.NewDefaultConfig()
	router := mux.NewRouter()
	browserRouter := router.PathPrefix("/browser").Subrouter()
	browserRouter.HandleFunc("/{repositoryID}", browserRepository)
	browserRouter.HandleFunc("/{repositoryID}/root", browserServices)
	log.Fatalf("Error running server -> %s", http.ListenAndServe(config.CmisAppPort, router))
}

func browserRepository(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repositoryID, ok := vars["repositoryID"]
	if ok {
		// Return repository information
	} else {
		// Return array of repository information
	}
}

func browserServices(w http.ResponseWriter, r *http.Request) {

}
