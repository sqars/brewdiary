package app

import (
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/config"

	// GORM drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App struct is entry point of application
type App struct {
	Config *config.Config
	DB     *gorm.DB
	Router *mux.Router
}

// NewApp is function constructor for main Application Object.
// Method returns created application object with passed config.
func NewApp(conf *config.Config) (*App, error) {
	app := &App{}
	app.Config = conf
	err := app.Init()
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Init method initialize application using config.
// Method establish connection with database and
// register app routes.
func (a *App) Init() error {
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
		return err
	}
	a.DB = db
	return nil
}

// Run starts application server
func (a *App) Run() error {
	log.Print("Starting server on:", a.Config.Host)
	return nil
}
