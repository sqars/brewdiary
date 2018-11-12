package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sqars/brewdiary/app/models"
)

func addTestBrewIngridient(brewID, count int) {
	if count < 1 {
		count = 1
	}
	brew := models.Brew{}
	brew.Get(a.DB)
	ingridients := []models.BrewIngridient{}
	for i := 0; i < count; i++ {
		brewIngridient := models.BrewIngridient{
			Quantity: (count + 1) * 100,
			BrewID:   brewID,
		}
		brewIngridient.Create(a.DB)
		ingridients = append(ingridients, brewIngridient)
	}
	brew.Ingridients = ingridients
	brew.Update(a.DB)
}

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
