package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sqars/brewdiary/app/models"
)

func TestBrewIngridientHandler_AddBrewIngridient(t *testing.T) {
	type args struct {
		payload []byte
		url     string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Should add new brewingridient to existing brew",
			args: args{
				payload: []byte(`{"quantity":100,"ingridientId":2}`),
				url:     "/brew/2/ingr",
			},
			wantCode: http.StatusCreated,
		}, {
			name: "Should return 404 if no brew to add ingridient",
			args: args{
				payload: []byte(`{"quantity":100,"ingridientId":2}`),
				url:     "/brew/20/ingr",
			},
			wantCode: http.StatusNotFound,
		}, {
			name: "Should return 400 if try to create with wrong payload",
			args: args{
				payload: []byte(`{"quantity":"will break","ingridientId":2}`),
				url:     "/brew/2/ingr",
			},
			wantCode: http.StatusBadRequest,
		}, {
			name: "Should return 400 if try to create with wrong payload",
			args: args{
				payload: []byte(`{"quantity":100,"ingridientId":"will break}`),
				url:     "/brew/2/ingr",
			},
			wantCode: http.StatusBadRequest,
		}, {
			name: "Should return 400 if try to create with wrong payload",
			args: args{
				payload: []byte(`{"quantity":100,"ingridientId":"will break}`),
				url:     "/brew/2/ingr",
			},
			wantCode: http.StatusBadRequest,
		}, {
			name: "Should return 400 if try to create with without ingridientId",
			args: args{
				payload: []byte(`{"quantity":100}`),
				url:     "/brew/2/ingr",
			},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearBrewTable()
			clearBrewIngridientTable()
			addTestBrew(5)
			addTestIngridient(3)
			req, err := http.NewRequest("POST", tt.args.url, bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusCreated {
				i := models.BrewIngridient{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&i)
				if err != nil {
					t.Errorf("Cannot decode api response")
				}
				if i.BrewID != 2 {
					t.Errorf("Created BrewIngridient should have BrewID 2. Got: %v", i.BrewID)
				}
				if i.Quantity != 100 {
					t.Errorf("Created BrewIngridient should have quantity 200. Got %v", i.Quantity)
				}
				if i.Ingridient.ID != 2 {
					t.Errorf("Created BrewIngridient should have Ingridient with ID 2. Got %v", i.Ingridient.ID)
				}
				if i.Ingridient.Comments != "Comments2" {
					t.Errorf(`Created BrewIngridient should have Ingridient with "Comments2". Got %v`, i.Ingridient.Comments)
				}
			}
		})
	}
}

func TestBrewIngridientHandler_GetBrewIngridient(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name               string
		args               args
		wantCode           int
		wantBrewIngridient models.BrewIngridient
	}{
		{
			name: "Should get existing brew ingridient for brew",
			args: args{
				url: "/brew/3/ingr/2",
			},
			wantCode: http.StatusOK,
			wantBrewIngridient: models.BrewIngridient{
				IngridientID: 2,
				BrewID:       3,
				Quantity:     200,
				Ingridient: models.Ingridient{
					Name: "Ingridient2",
				},
			},
		}, {
			name: "Should return 404 if no brew",
			args: args{
				url: "/brew/20/ingr/2",
			},
			wantCode: http.StatusNotFound,
		}, {
			name: "Should return 404 if no ingridient",
			args: args{
				url: "/brew/3/ingr/220",
			},
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTables()
			addTestBrew(5)
			addTestBrewIngridient(3, 3)
			req, err := http.NewRequest("GET", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusOK {
				i := models.BrewIngridient{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&i)
				if i.IngridientID != tt.wantBrewIngridient.IngridientID {
					t.Errorf("BrewIngridient should have IngridientID equals to %v. Got: %v", tt.wantBrewIngridient.IngridientID, i.IngridientID)
				}
				if i.BrewID != tt.wantBrewIngridient.BrewID {
					t.Errorf("BrewIngridient should have BrewID equals to %v. Got: %v", tt.wantBrewIngridient.BrewID, i.BrewID)
				}
				if i.Quantity != tt.wantBrewIngridient.Quantity {
					t.Errorf("BrewIngridient should have quantity equal to %v. Got: %v", tt.wantBrewIngridient.Quantity, i.Quantity)
				}
				if i.Ingridient.Name != tt.wantBrewIngridient.Ingridient.Name {
					t.Errorf(`BrewIngridient should include Ingridient with name "%v". Got: %v`, tt.wantBrewIngridient.Ingridient.Name, i.Ingridient.Name)
				}
			}
		})
	}
}

func TestBrewIngridientHandler_DeleteBrewIngridient(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Should return 200 if removed",
			args: args{
				url: "/brew/3/ingr/2",
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTables()
			addTestBrew(5)
			addTestBrewIngridient(3, 3)
			req, err := http.NewRequest("DELETE", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusOK {
				req, err = http.NewRequest("GET", tt.args.url, nil)
				if err != nil {
					t.Errorf("Cannot create request after removing brew")
				}
				response = executeRequest(req)
				checkResponseCode(t, http.StatusNotFound, response.Code)
			}
		})
	}
}

func TestBrewIngridientHandler_UpdateBrewIngridient(t *testing.T) {
	type args struct {
		url     string
		payload []byte
	}
	tests := []struct {
		name               string
		args               args
		wantCode           int
		wantBrewIngridient models.BrewIngridient
	}{
		{
			name: "Should update existing brew ingridient for brew",
			args: args{
				url:     "/brew/3/ingr/2",
				payload: []byte(`{"ingridientId":3,"quantity":666}`),
			},
			wantCode: http.StatusOK,
			wantBrewIngridient: models.BrewIngridient{
				IngridientID: 3,
				BrewID:       3,
				Quantity:     666,
				Ingridient: models.Ingridient{
					Name: "Ingridient3",
				},
			},
		}, {
			name: "Should get 404 if no brew",
			args: args{
				url:     "/brew/555/ingr/2",
				payload: []byte(`{"ingridientId":3,"quantity":666}`),
			},
			wantCode: http.StatusNotFound,
		}, {
			name: "Should get 404 if no brew ingridient to update",
			args: args{
				url:     "/brew/2/ingr/555",
				payload: []byte(`{"ingridientId":3,"quantity":666}`),
			},
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTables()
			addTestBrew(5)
			addTestBrewIngridient(3, 3)
			req, err := http.NewRequest("PATCH", tt.args.url, bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusOK {
				i := models.BrewIngridient{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&i)
				if i.IngridientID != tt.wantBrewIngridient.IngridientID {
					t.Errorf("BrewIngridient should have IngridientID equals to %v. Got: %v", tt.wantBrewIngridient.IngridientID, i.IngridientID)
				}
				if i.BrewID != tt.wantBrewIngridient.BrewID {
					t.Errorf("BrewIngridient should have BrewID equals to %v. Got: %v", tt.wantBrewIngridient.BrewID, i.BrewID)
				}
				if i.Quantity != tt.wantBrewIngridient.Quantity {
					t.Errorf("BrewIngridient should have quantity equal to %v. Got: %v", tt.wantBrewIngridient.Quantity, i.Quantity)
				}
				if i.Ingridient.Name != tt.wantBrewIngridient.Ingridient.Name {
					t.Errorf(`BrewIngridient should include Ingridient with name "%v". Got: %v`, tt.wantBrewIngridient.Ingridient.Name, i.Ingridient.Name)
				}
			}
		})
	}
}
