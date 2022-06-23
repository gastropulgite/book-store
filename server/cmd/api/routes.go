package main

import (
	"book-api/internal/data"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)


func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*","http://*" },
		AllowedMethods: []string{"GET", "POST", "PUT", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	r.Get(`/users/login`, app.Login)
	r.Post(`/users/login`, app.Login)

	r.Get(`/users/`, func(w http.ResponseWriter, r *http.Request) {

		var users data.User
		result, err := users.GetUsers()

		if err != nil {
			app.errorLog.Println(err)
			return
		}
		app.writeJSON(w, http.StatusOK, result)
	})
	
	r.Get(`/users/add`, func(w http.ResponseWriter, r *http.Request) {
		var u = data.User{
			Email: `hieuminh@gmail.com`,
			FirstName:  `Hieu`,
			LastName: `Minh`,
			Password: `password`,
		}

		app.infoLog.Println(`Adding user!`)
		id, err := app.models.User.Insert(u)
		if err != nil {
			app.infoLog.Printf(err.Error())
			app.errorJSON(w, err , http.StatusForbidden)
			return
		}

		app.infoLog.Println(`Got back ID of`, id)
		newUser, err := app.models.User.GetUserById(id)

		if err != nil {
			app.infoLog.Printf(err.Error())
			app.errorJSON(w, err , http.StatusNotFound)
			return
		}

		app.writeJSON(w, http.StatusOK, newUser)
		return
	})
	
	return r
}