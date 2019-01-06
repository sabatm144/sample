package main

import (
	"fmt"
	"log"
	"net/http"
	"sample/server/dbCon"
	"sample/server/routes"
	"time"

	"github.com/rs/cors"
)

type headerMap map[string]int

func main() {

	dbCon.ConnectMongoDB()
	port := 9000
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type"},
	})

	server := &http.Server{
		Handler: c.Handler(routes.HTTPRouteConfig()),
		Addr:    fmt.Sprintf(":%d", port),
		//Good practice: enforces timeout for the server
		WriteTimeout:   90 * time.Second,
		ReadTimeout:    90 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Server is listening on: %d", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Unable to listen: %v\n", err)
	}

}
