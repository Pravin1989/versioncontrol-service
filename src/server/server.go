package server

import (
	"fmt"
	"log"
	"net/http"

	"versioncontrol-service/src/common"
	"versioncontrol-service/src/config"
	handlers "versioncontrol-service/src/services"

	"github.com/gorilla/mux"
)

var (
	handleHomePage               = handlers.HandleHomePage
	handleCallBackFromGithubAuth = handlers.HandleCallBackFromGithubAuth
)

//This method add routes and start server
func Start() {
	router := loadRoutes()
	log.Println("Starting REST Server")
	log.Printf("REST server listening on port %d", config.ServerConfig.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerConfig.Port), router)
	log.Println("Server Crashed : ", err)
}

//Creates router
func loadRoutes() *mux.Router {
	serviceRouter := mux.NewRouter().PathPrefix(fmt.Sprintf("/%s", common.ContextRoot)).Subrouter()
	registerApiRoutes(serviceRouter)
	return serviceRouter
}

//Regiter API methods inside router
func registerApiRoutes(r *mux.Router) {
	r.HandleFunc("/", handleHomePage).Methods(http.MethodGet)
	r.HandleFunc("/callback", handleCallBackFromGithubAuth).Methods(http.MethodGet)
}
