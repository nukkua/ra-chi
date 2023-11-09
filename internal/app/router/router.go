package router;

import (
	"github.com/nukkua/ra-chi/internal/app/handlers"
	"github.com/nukkua/ra-chi/internal/app/database"


	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter () *chi.Mux{
	db:= database.SetupDatabase();
	r:= chi.NewRouter();
	
	r.Use(middleware.Logger);
	
	r.Get("/users", handlers.GetUsers(db))
	r.Post("/user", handlers.CreateUser(db))

	return r;

}

