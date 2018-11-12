package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/sqars/brewdiary/app/models"
)

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

func TestIngridientHandler_GetIngridient(t *testing.T) {
	type args struct {
		url string
		IID int
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "Should return 404 if no Ingridient in db",
			args: args{
				url: "/ingr/666",
			},
			wantCode: http.StatusNotFound,
		}, {
			name: "Should return 200 and Ingridient data",
			args: args{
				url: "/ingr/5",
				IID: 5,
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearIngridientTable()
			if tt.wantCode == http.StatusOK {
				addTestIngridient(tt.args.IID + 1)
			}
			req, err := http.NewRequest("GET", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusOK {
				i := models.Ingridient{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&i)
				if err != nil {
					t.Errorf("Cannot decode api response")
				}
				expectedName := "Ingridient" + strconv.Itoa(tt.args.IID)
				expectedComments := "Comments" + strconv.Itoa(tt.args.IID)
				if i.Name != expectedName {
					expectMsg(t, expectedName, i.Name)
				}
				if i.Comments != expectedComments {
					expectMsg(t, expectedComments, i.Comments)
				}
			}
		})
	}
}

func TestIngridient_AddIngridient(t *testing.T) {
	type args struct {
		payload []byte
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantIngr models.Ingridient
	}{
		{
			name: "Should return 422 if wrong payload passed",
			args: args{
				payload: []byte(`"name":"TestIngridient","wtf":"should break"`),
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
				payload: []byte(`"name":"123"`),
			},
			wantCode: http.StatusUnprocessableEntity,
		}, {
			name: "Should return 201 if Ingridient was created with proper payload",
			args: args{
				payload: []byte(`{"name":"TestIngridient","comments":"TestComments"}`),
			},
			wantCode: http.StatusCreated,
			wantIngr: models.Ingridient{
				Name:     "TestIngridient",
				Comments: "TestComments",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearIngridientTable()
			req, err := http.NewRequest("POST", "/ingr", bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if tt.wantCode == http.StatusCreated {
				i := models.Ingridient{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&i)
				if err != nil {
					t.Errorf("Cannot decode api response")
				}
				if i.Name != tt.wantIngr.Name {
					expectMsg(t, tt.wantIngr.Name, i.Name)
				}
				if i.Comments != tt.wantIngr.Comments {
					expectMsg(t, tt.wantIngr.Comments, i.Comments)
				}
			}
		})
	}
}

func TestIngridient_DeleteIngridient(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should delete ingridient from DB",
			args: args{
				url: "/ingr/2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearIngridientTable()
			addTestIngridient(5)
			req, err := http.NewRequest("DELETE", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			req, err = http.NewRequest("GET", tt.args.url, nil)
			if err != nil {
				t.Errorf("Cannot create request after removing ingridient")
			}
			response = executeRequest(req)
			checkResponseCode(t, http.StatusNotFound, response.Code)
		})
	}
}

func TestIngridient_UpdateIngridient(t *testing.T) {
	type args struct {
		url     string
		payload []byte
	}
	tests := []struct {
		name              string
		args              args
		wantCode          int
		ingridientUpdated models.Ingridient
	}{
		{
			name: "Should return 404 code if no ingridient to update",
			args: args{
				url:     "/ingr/20",
				payload: []byte(`{"name": "NameUpdated"}`),
			},
			wantCode: http.StatusNotFound,
		}, {
			name: "Should return 422 code if wrong payload",
			args: args{
				url:     "/ingr/3",
				payload: []byte(`{"name": 123}`),
			},
			wantCode: http.StatusUnprocessableEntity,
		}, {
			name: "Should return 200 code and User data when updated",
			args: args{
				url:     "/ingr/3",
				payload: []byte(`{"name": "NameUpdated","comments": "CommentsUpdated"}`),
			},
			wantCode: http.StatusOK,
			ingridientUpdated: models.Ingridient{
				Name:     "NameUpdated",
				Comments: "CommentsUpdated",
			},
		}, {
			name: "Should return 200 code and User data with updated only 1 fields",
			args: args{
				url:     "/ingr/4",
				payload: []byte(`{"comments": "CommentsUpdated"}`),
			},
			wantCode: http.StatusOK,
			ingridientUpdated: models.Ingridient{
				Name:     "Ingridient4",
				Comments: "CommentsUpdated",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearIngridientTable()
			addTestIngridient(5)
			req, err := http.NewRequest("PATCH", tt.args.url, bytes.NewBuffer(tt.args.payload))
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, tt.wantCode, response.Code)
			if response.Code == http.StatusOK {
				i := models.Ingridient{}
				decoder := json.NewDecoder(response.Body)
				err = decoder.Decode(&i)
				if err != nil {
					t.Errorf("Cannot decode api response")
				}
				if i.Name != tt.ingridientUpdated.Name {
					expectMsg(t, tt.ingridientUpdated.Name, i.Name)
				}
				if i.Comments != tt.ingridientUpdated.Comments {
					expectMsg(t, tt.ingridientUpdated.Comments, i.Comments)
				}
			}
		})
	}
}

func TestIngridient_GetIngridients(t *testing.T) {
	type args struct {
		createIngridients bool
	}
	tests := []struct {
		name          string
		args          args
		expectedItems int
	}{
		{
			name: "Should return empty array",
			args: args{
				createIngridients: false,
			},
			expectedItems: 0,
		}, {
			name: "Should return created ingridients",
			args: args{
				createIngridients: true,
			},
			expectedItems: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearIngridientTable()
			if tt.args.createIngridients {
				addTestIngridient(tt.expectedItems)
			}
			req, err := http.NewRequest("GET", "/ingr", nil)
			if err != nil {
				t.Errorf("Cannot create request")
			}
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)
			i := []models.Ingridient{}
			decoder := json.NewDecoder(response.Body)
			err = decoder.Decode(&i)
			if len(i) != tt.expectedItems {
				t.Errorf("Expected response to have %v items. Got %v", tt.expectedItems, len(i))
			}
		})
	}
}
