package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

func clearTables() {
	clearBrewTable()
	clearIngridientTable()
	clearBrewIngridientTable()
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

func addTestBrewIngridient(brewID, count int) {
	if count < 1 {
		count = 1
	}
	brew := models.Brew{}
	brew.Get(a.DB)
	ingridients := []models.BrewIngridient{}
	addTestIngridient(count)
	for i := 0; i < count; i++ {
		brewIngridient := models.BrewIngridient{
			Quantity:     (i + 1) * 100,
			BrewID:       brewID,
			IngridientID: i + 1,
		}
		brewIngridient.Create(a.DB)
		ingridients = append(ingridients, brewIngridient)
	}
	brew.Ingridients = ingridients
	brew.Update(a.DB)
}

func addTestBrew(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		b := models.Brew{
			Name:        "Brew" + strconv.Itoa(i+1),
			Location:    "Location" + strconv.Itoa(i+1),
			Comments:    "Comments" + strconv.Itoa(i+1),
			Ingridients: []models.BrewIngridient{},
		}
		a.DB.Create(&b)
	}
}

func addTestIngridient(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		in := models.Ingridient{
			Name:     "Ingridient" + strconv.Itoa(i+1),
			Comments: "Comments" + strconv.Itoa(i+1),
		}
		a.DB.Create(&in)
	}
}

func brewEquals(t *testing.T, bExp, bAct models.Brew) {
	// we dont want to compare time fields
	options := cmp.FilterValues(func(x, y interface{}) bool {
		return reflect.TypeOf(x) == reflect.TypeOf(time.Time{}) || reflect.TypeOf(y) == reflect.TypeOf(time.Time{})
	}, cmp.Ignore())
	diff := cmp.Diff(bExp, bAct, options)
	if diff != "" {
		t.Errorf("Brews differ: (-want +got)\n%s", diff)
	}
}
