package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sqars/brewdiary/app/models"
)

func TestBrewHandler_GetBrew(t *testing.T) {
	type args struct {
		url                string
		brewID             int
		howManyIngridients int
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBrew models.Brew
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
				url:                "/brew/5",
				brewID:             5,
				howManyIngridients: 2,
			},
			wantCode: http.StatusOK,
			wantBrew: models.Brew{
				ID:       5,
				Name:     "Brew5",
				Comments: "Comments5",
				Location: "Location5",
				Ingridients: []models.BrewIngridient{
					models.BrewIngridient{
						ID:           1,
						Quantity:     100,
						IngridientID: 1,
						Ingridient: models.Ingridient{
							ID:       1,
							Name:     "Ingridient1",
							Comments: "Comments1",
						},
					},
					models.BrewIngridient{
						ID:           2,
						Quantity:     200,
						IngridientID: 2,
						Ingridient: models.Ingridient{
							ID:       2,
							Name:     "Ingridient2",
							Comments: "Comments2",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearBrewTable()
			if tt.wantCode == http.StatusOK {
				addTestBrew(tt.args.brewID + 1)
				addTestBrewIngridient(tt.args.brewID, tt.args.howManyIngridients)
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
				brewEquals(t, tt.wantBrew, b)
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
				ID:       1,
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
				brewEquals(t, tt.wantBrew, b)
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
				ID:          3,
				Name:        "NameUpdated",
				Comments:    "CommentsUpdated",
				Location:    "LocationUpdated",
				Ingridients: []models.BrewIngridient{},
			},
		}, {
			name: "Should return 200 code and brew data with updated only 2 fields",
			args: args{
				url:     "/brew/4",
				payload: []byte(`{"comments": "CommentsUpdated","location":"LocationUpdated"}`),
			},
			wantCode: http.StatusOK,
			brewUpdated: models.Brew{
				ID:          4,
				Name:        "Brew4",
				Comments:    "CommentsUpdated",
				Location:    "LocationUpdated",
				Ingridients: []models.BrewIngridient{},
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
				brewEquals(t, tt.brewUpdated, b)
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
