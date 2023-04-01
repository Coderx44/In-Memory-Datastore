package server

import (
	"net/http"

	"github.com/Coderx44/gg/commands"
	"github.com/Coderx44/gg/middleware"
	"github.com/gorilla/mux"
)

func InitRouter(dep dependencies) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/commands", middleware.ParseCommands(commands.HandleCommands(dep.dt))).Methods(http.MethodPost)

	return
}
