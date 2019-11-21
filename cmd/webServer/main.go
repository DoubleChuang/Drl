package main

import (	
	"log"
	"net/http"
	"os"

	"github.com/DoubleChuang/Drl/api"
	"github.com/DoubleChuang/Drl/pkg/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}	
	r := api.RegisterHandlers()
	m := middleware.NewMiddleWareHandler(r)
	if err := http.ListenAndServe(":"+port, m); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
