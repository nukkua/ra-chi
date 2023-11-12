package router;

import (
	"github.com/nukkua/ra-chi/internal/app/handlers"
	"github.com/nukkua/ra-chi/internal/app/database"
	

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nukkua/ra-chi/internal/app/middlewares"
)

func SetupRouter () *chi.Mux{
	db:= database.SetupDatabase();
	
	r:= chi.NewRouter();
	r.Use(middleware.Logger);
	r.Post("/register", handlers.CreateUser(db))
	r.Post("/login", handlers.LoginUser(db))

	r.Group(func(r chi.Router){
		r.Use(middlewares.JwtAuthentication)
		r.Get("/users", handlers.GetUsers(db))
	})




	return r;

}

