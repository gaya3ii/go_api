package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gayathriad/go_api/db"
	"github.com/gayathriad/go_api/handlers"
	"github.com/gayathriad/go_api/middleware"
	"github.com/gayathriad/go_api/repository"
)

func main() {
	// connect to database
	database, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// create repository
	userRepo := repository.NewUserRepo(database)

	// create handler
	userHandler := handlers.NewUserHandler(userRepo)

	// create router
	mux := http.NewServeMux()

	// register routes
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUsers(w, r)
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		}
	})
	mux.HandleFunc("/users/bulk", userHandler.CreateUsersBulk)
	mux.HandleFunc("/users/get", userHandler.GetUserByID)
	mux.HandleFunc("/users/update", userHandler.UpdateUser)
	mux.HandleFunc("/users/delete", userHandler.DeleteUser)

	// wrap with middleware
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", middleware.Logging(mux))
}
