package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/config"
	_ "github.com/rsilraf/pos_goexpert/desafios/multithreading/docs"
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/entity"
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/infra/db"
	"github.com/rsilraf/pos_goexpert/desafios/multithreading/internal/infra/web/handlers"
	swagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title			Go Expert - CEP API
// @version			1.0
// @description		CEP API with multithreading and authentication
// @host			localhost:8000
// @basePath		/
// @securityDefinitions.apikey	ApiKeyAuth
// @in header
// @name Authorization
func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	orm, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	orm.AutoMigrate(&entity.User{})

	userHandler := handlers.NewUserHandler(db.NewUserDAO(orm))
	cepHandler := handlers.NewCepHandler()

	// router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("token", conf.JWT))
	r.Use(middleware.WithValue("TTL", conf.JWTTTL))

	// /cep
	r.Route("/cep", func(r chi.Router) {
		r.Use(jwtauth.Verifier(conf.JWT))
		r.Use(jwtauth.Authenticator)

		r.Get("/{cep}", cepHandler.GetCep)
	})
	// user
	r.Post("/users", userHandler.Create)

	// /token
	r.Post("/token", userHandler.GetToken)

	// /docs
	r.Get("/docs/*", swagger.Handler(
		swagger.URL("http://localhost:8000/docs/doc.json"),
	))

	http.ListenAndServe(":8000", r)
}
