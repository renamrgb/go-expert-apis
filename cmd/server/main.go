// Package main
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	// "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/renamrgb/go-expert-apis/configs"
	_ "github.com/renamrgb/go-expert-apis/docs"
	"github.com/renamrgb/go-expert-apis/internal/entity"
	"github.com/renamrgb/go-expert-apis/internal/infra/database"
	"github.com/renamrgb/go-expert-apis/internal/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @Title Go Expert API Exemple
// @Version 1.0
// @Description Product API with JWT Authentication
// @termsOfService http://swagger.io/terms/

//@Contact.name Renam Bulhoes
//@Contact.url
//@Contact.email renamgustavo@live.com

//@license.name MIT
//@license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.User{}, &entity.Product{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, conf.TokenAuth, conf.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(conf.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetAllProducts)
		r.Get("/{id}", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})
	r.Route("/user", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
