package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/core-go/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"go-service/internal/app"
)

func main() {
	conf := app.Config{}
	config.Load(&conf, "configs/config")

	r := mux.NewRouter()
	er2 := app.Route(r, context.Background(), conf)
	if er2 != nil {
		panic(er2)
	}
	/*
		headersOk := handlers.AllowedHeaders([]string{"*"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS})
		StartServer(handlers.CORS(originsOk, headersOk, methodsOk)(r), conf.Server)
	*/
	handler := cors.AllowAll().Handler(r)
	fmt.Println("Start service " + conf.Server.Name)
	if err := http.ListenAndServe(Addr(conf.Server.Port), handler); err != nil {
		panic(err)
	}
}
func Addr(port *int64) string {
	server := ""
	if port != nil && *port >= 0 {
		server = ":" + strconv.FormatInt(*port, 10)
	}
	return server
}
