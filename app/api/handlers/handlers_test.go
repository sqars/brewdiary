package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sqars/brewdiary/app"
	"github.com/sqars/brewdiary/app/models"
	"github.com/sqars/brewdiary/config"
)

var a app.App

func TestMain(m *testing.M) {
	a = app.App{
		Config: config.Config{
			DbName:     "brewdiary_test",
			DbUser:     "postgres",
			DbPassword: "postgres",
			Host:       "localhost:4200",
			DbHost:     "localhost",
			DbPort:     "5432",
		},
	}
	a.Init()

	clearBrewTable()
	clearIngridientTable()
	clearBrewIngridientTable()
	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func expectMsg(t *testing.T, expected, actual string) {
	t.Errorf("Expected to have %s. Got %s", expected, actual)
}

func clearBrewTable() {
	a.DB.DropTableIfExists(&models.Brew{})
	a.DB.AutoMigrate(&models.Brew{})
}

func clearIngridientTable() {
	a.DB.DropTableIfExists(&models.Ingridient{})
	a.DB.AutoMigrate(&models.Ingridient{})
}

func clearBrewIngridientTable() {
	a.DB.DropTableIfExists(&models.BrewIngridient{})
	a.DB.AutoMigrate(&models.BrewIngridient{})
}
