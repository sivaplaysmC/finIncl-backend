package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sivaplaysmc/finIncl-backend/config"
	"github.com/sivaplaysmc/finIncl-backend/internal/helpers"
	"github.com/sivaplaysmc/finIncl-backend/internal/routes"
)

func main() {
	godotenv.Load()
	app, err := config.GetApp()
	if err != nil {
		log.Fatalln(err)
	}
	router := mux.NewRouter()
	helpers.Mount(router, "/user", routes.GetUsersRouter(app))

	server := http.Server{
		Addr:    "192.168.170.32:6969",
		Handler: router,
	}
	app.InfoLog.Println("Server Listening at", server.Addr)

	server.ListenAndServe()
}
