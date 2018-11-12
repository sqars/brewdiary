package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/sqars/brewdiary/app/models"
)

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

func TestBrewHandler_GetBrew(t *testing.T) {
	type args struct {
		url    string
		brewID int
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Should return 404 if no Brew in db",
			args: args{
				url: "/brew/666",
			},
			wantCode: http.StatusNotFound,
		}, {
			name: "Should return 200 and Brew data",
			args: args{
				url:    "/brew/5",
				brewID: 5,
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearBrewTable()
			if tt.wantCode == http.StatusOK {
				addTestBrew(tt.args.brewID + 1)
			}
			req, err := http.NewRequest("GET", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusOK {
				b := models.Brew{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&b)
				if err != nil {
					t.Errorf("Cannot decode api response")
				}
				expectedName := "Brew" + strconv.Itoa(tt.args.brewID)
				expectedLocation := "Location" + strconv.Itoa(tt.args.brewID)
				expectedComments := "Comments" + strconv.Itoa(tt.args.brewID)
				if b.Name != expectedName {
					expectMsg(t, expectedName, b.Name)
				}
				if b.Location != expectedLocation {
					expectMsg(t, expectedLocation, b.Location)
				}
				if b.Comments != expectedComments {
					expectMsg(t, expectedComments, b.Comments)
				}
				if len(b.Ingridients) != 0 {
					expectMsg(t, " no ingridients", strconv.Itoa(len(b.Ingridients))+" ingridients")
				}
			}
		})
	}
}

func TestBrewHandler_AddBrew(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBrew models.Brew
	}{
		{
			name: "Should return 422 if wrong payload passed",
			args: args{
				payload: []byte(`"name":"TestBrew","wtf":"should break"`),
			},
			wantCode: http.StatusUnprocessableEntity,
		}, {
			name: "Should return 422 if nothing passed in payload",
			args: args{
				payload: []byte(``),
			},
			wantCode: http.StatusUnprocessableEntity,
		}, {
			name: "Should return 422 if nothing passed wrong data types in fields",
			args: args{
				payload: []byte(`"name":123,"location":555`),
			},
			wantCode: http.StatusUnprocessableEntity,
		}, {
			name: "Should return 201 if Brew was created with proper payload",
			args: args{
				payload: []byte(`{"name":"TestBrew","location":"TestLocation","comments":"TestComments"}`),
			},
			wantCode: http.StatusCreated,
			wantBrew: models.Brew{
				Name:     "TestBrew",
				Location: "TestLocation",
				Comments: "TestComments",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearBrewTable()
			req, err := http.NewRequest("POST", "/brew", bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusCreated {
				b := models.Brew{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&b)
				if err != nil {
					t.Errorf("Cannot decode api response")
				}
				if b.Name != tt.wantBrew.Name {
					expectMsg(t, tt.wantBrew.Name, b.Name)
				}
				if b.Comments != tt.wantBrew.Comments {
					expectMsg(t, tt.wantBrew.Comments, b.Comments)
				}
				if b.Location != tt.wantBrew.Location {
					expectMsg(t, tt.wantBrew.Location, b.Location)
				}
			}
		})
	}
}

func TestBrewHandler_DeleteBrew(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should delete brew from DB",
			args: args{
				url: "/brew/3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearBrewTable()
			addTestBrew(5)
			req, err := http.NewRequest("DELETE", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			req, err = http.NewRequest("GET", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request after removing brew")
			}
			response = executeRequest(req)
			checkResponseCode(t, http.StatusNotFound, response.Code)
		})
	}
}

func TestBrewHandler_UpdateBrew(t *testing.T) {
	type args struct {
		url     string
		payload []byte
	}
	tests := []struct {
		name        string
		args        args
		wantCode    int
		brewUpdated models.Brew
	}{
		{
			name: "Should return 404 code if no brew to update",
			args: args{
				url:     "/brew/20",
				payload: []byte(`{"name": "NameUpdated"}`),
			},
			wantCode: http.StatusNotFound,
		}, {
			name: "Should return 422 code if wrong payload",
			args: args{
				url:     "/brew/3",
				payload: []byte(`{"name": 123}`),
			},
			wantCode: http.StatusUnprocessableEntity,
		}, {
			name: "Should return 200 code and brew data when updated",
			args: args{
				url:     "/brew/3",
				payload: []byte(`{"name": "NameUpdated","comments": "CommentsUpdated","location":"LocationUpdated"}`),
			},
			wantCode: http.StatusOK,
			brewUpdated: models.Brew{
				Name:     "NameUpdated",
				Comments: "CommentsUpdated",
				Location: "LocationUpdated",
			},
		}, {
			name: "Should return 200 code and brew data with updated only 2 fields",
			args: args{
				url:     "/brew/4",
				payload: []byte(`{"comments": "CommentsUpdated","location":"LocationUpdated"}`),
			},
			wantCode: http.StatusOK,
			brewUpdated: models.Brew{
				Name:     "Brew4",
				Comments: "CommentsUpdated",
				Location: "LocationUpdated",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearBrewTable()
			addTestBrew(5)
			req, err := http.NewRequest("PATCH", tt.args.url, bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if response.Code == http.StatusOK {
				b := models.Brew{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&b)
				if err != nil {
					t.Errorf("Cannot decode api response")
				}
				if b.Name != tt.brewUpdated.Name {
					expectMsg(t, tt.brewUpdated.Name, b.Name)
				}
				if b.Comments != tt.brewUpdated.Comments {
					expectMsg(t, tt.brewUpdated.Comments, b.Comments)
				}
				if b.Location != tt.brewUpdated.Location {
					expectMsg(t, tt.brewUpdated.Location, b.Location)
				}
			}
		})
	}
}

func TestBrewHandler_GetBrews(t *testing.T) {
	type args struct {
		createBrews bool
	}
	tests := []struct {
		name          string
		args          args
		expectedItems int
	}{
		{
			name: "Should return empty array",
			args: args{
				createBrews: false,
			},
			expectedItems: 0,
		}, {
			name: "Should return created brews",
			args: args{
				createBrews: true,
			},
			expectedItems: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearBrewTable()
			if tt.args.createBrews {
				addTestBrew(tt.expectedItems)
			}
			req, err := http.NewRequest("GET", "/brew", nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			b := []models.Brew{}
			decoder := json.NewDecoder(response.Body)
			err = decoder.Decode(&b)
			if len(b) != tt.expectedItems {
				t.Errorf("Expected response to have %v items. Got %v", tt.expectedItems, len(b))
			}
		})
	}
}
