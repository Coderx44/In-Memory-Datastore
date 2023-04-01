package server

import (
	"fmt"
	"strconv"

	"github.com/urfave/negroni"
)

func StartApiServer() {
	port := 3000
	server := negroni.Classic()

	dependencies, err := InitDependencies()
	if err != nil {
		panic(err)
	}

	router := InitRouter(dependencies)
	server.UseHandler(router)

	addr := fmt.Sprintf(":%s", strconv.Itoa(port))
	server.Run(addr)

}
