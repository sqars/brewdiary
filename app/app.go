package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/app/api"
	"github.com/sqars/brewdiary/config"

	// GORM drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App struct is entry point of application
type App struct {
	Config config.Config
	DB     *gorm.DB
	Router *mux.Router
}

// NewApp is function constructor for main Application Object.
// Method returns created application object with passed config.
func NewApp(conf config.Config) *App {
	app := &App{}
	app.Config = conf
	app.Init()
	return app
}

// Init method initialize application using config.
// Method establish connection with database and
// register app routes.
func (a *App) Init() {
	a.connectDB()
	a.initAPI()
}

func (a *App) connectDB() {
	connDetails := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		a.Config.DbHost,
		a.Config.DbPort,
		a.Config.DbUser,
		a.Config.DbName,
		a.Config.DbPassword,
	)
	db, err := gorm.Open("postgres", connDetails)
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
}

func (a *App) initAPI() {
	r := api.Init(a.DB)
	a.Router = r
}

// Run starts application server
func (a *App) Run() error {
	allowedOrigins := handlers.AllowedOrigins(a.Config.Cors)
	allowedMethods := handlers.AllowedMethods([]string{
		"POST", "PUT", "GET", "DELETE", "PATCH",
	})
	allowedHeaders := handlers.AllowedHeaders([]string{
		"Content-Type",
	})
	log.Print("Starting web server on addres: ", a.Config.Host)
	return http.ListenAndServe(
		a.Config.Host,
		handlers.CORS(allowedMethods, allowedOrigins, allowedHeaders)(a.Router),
	)
}
