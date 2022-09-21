package eventcontroller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/rizface/golang-api-template/app/controller/eventcontroller"
	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/entity/responseentity"
	"github.com/rizface/golang-api-template/app/repository/eventrepository"
	"github.com/rizface/golang-api-template/app/service/eventservice"
	"github.com/rizface/golang-api-template/database"
	"github.com/rizface/golang-api-template/database/postgresql"
	"github.com/rizface/golang-api-template/system/router"
	"github.com/rizface/golang-api-template/system/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var config *database.Database
var db *gorm.DB
var r *chi.Mux
var token string
var repository eventrepository.EventRepositoryInterface
var service eventservice.EventServiceInterface
var controller eventcontroller.EventControllerInterface
var id string

func TestMain(m *testing.M) {
	config = &database.Database{
		Name:     "ms-event",
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "root",
	}
	db = postgresql.New(config)
	r = router.CreateRouter()
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGF0YSI6eyJpZCI6IjdiYzE1YmI2LTU2N2UtNDA2Yi05ZjcxLWJkODNjM2M5MDEyZiIsIm5hbWUiOiJmYXJpeiIsInVzZXJuYW1lIjoibWFyZGlhbnRvIn0sImlzcyI6IlJJWkZBQ0UiLCJzdWIiOiJmYXJpeiIsImV4cCI6MTY2Mzg1OTA1MSwibmJmIjoxNjYzNzcyNjUxLCJpYXQiOjE2NjM3NzI2NTF9.ZaCM2PGeQv6HzMhrBDeyFpfD24NYQcXBWnKISsD9-3E"
	repository = eventrepository.New()
	service = eventservice.New(repository, db)
	controller := eventcontroller.New(service)
	eventcontroller.Setup(r, controller)
	dummyData := &requestentity.Event{
		Name:        "Event For Integration Testing",
		Description: "Event For Integration Testing",
		Deadline:    time.Now().Add(5 * 24 * time.Hour).Format(time.RFC3339),
	}
	claim := security.DecodeToken(token, "Bearer")
	result, _ := repository.Create(db, claim.UserData.Id, dummyData)
	id = result.Id
	m.Run()
}

func TestGetEvent(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetOneEvent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", id), nil)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		response := new(responseentity.Response)
		bytesResponse, _ := ioutil.ReadAll(recorder.Result().Body)
		json.Unmarshal(bytesResponse, response)
		event := response.Data.(map[string]interface{})
		assert.Equal(t, id, event["id"])
	})
	t.Run("Event Not Found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.NewString()), nil)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
	t.Run("Unauthorized", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", id), nil)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}

func TestDeleteEvent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/"+id, nil)
		recorder := httptest.NewRecorder()
		request.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusOK, recorder.Code)

		response := new(responseentity.Response)
		responseBytes, _ := ioutil.ReadAll(recorder.Result().Body)
		json.Unmarshal(responseBytes, response)
		assert.Equal(t, id, response.Data.(map[string]interface{})["id"])
	})
	t.Run("Not Found", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/"+uuid.NewString(), nil)
		recorder := httptest.NewRecorder()
		request.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}

func TestCreateEvent(t *testing.T) {
	t.Run("SUCCESS", func(t *testing.T) {
		payload := requestentity.Event{
			Name:        "Event Testing",
			Description: "Event created from integration testing",
			Deadline:    time.Now().Format(time.RFC3339),
		}

		bytesRequest, _ := json.Marshal(payload)
		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bytesRequest))
		recorder := httptest.NewRecorder()
		request.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, request)

		response := new(responseentity.Response)
		bytesResponse, _ := ioutil.ReadAll(recorder.Result().Body)
		json.Unmarshal(bytesResponse, response)
		err := validation.Validate(
			response.Data.(map[string]interface{})["id"],
			is.UUID,
		)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Nil(t, err)
		id = response.Data.(map[string]interface{})["id"].(string)
	})

	t.Run("BAD REQUEST", func(t *testing.T) {
		payload := requestentity.Event{}

		bytesRequest, _ := json.Marshal(payload)
		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bytesRequest))
		recorder := httptest.NewRecorder()
		request.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("UNAUTHORIZED", func(t *testing.T) {
		payload := requestentity.Event{}

		bytesRequest, _ := json.Marshal(payload)
		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bytesRequest))
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}

func TestUpdateEvent(t *testing.T) {
	t.Run("SUCCESS", func(t *testing.T) {
		payload := requestentity.Event{
			Name:        "Update Event",
			Description: "Event Updated",
			Deadline:    time.Now().Format(time.RFC3339),
		}
		bytesPayload, _ := json.Marshal(payload)
		request := httptest.NewRequest(http.MethodPut, "/"+id, bytes.NewReader(bytesPayload))
		request.Header.Add("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)

		response := new(responseentity.Response)
		bytesResponse, _ := ioutil.ReadAll(recorder.Result().Body)
		json.Unmarshal(bytesResponse, response)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, id, response.Data.(map[string]interface{})["id"].(string))
		assert.Equal(t, payload.Name, response.Data.(map[string]interface{})["name"].(string))
	})

	t.Run("Unauhorized", func(t *testing.T) {
		payload := requestentity.Event{
			Name:        "Update Event",
			Description: "Event Updated",
			Deadline:    time.Now().Format(time.RFC3339),
		}
		bytesPayload, _ := json.Marshal(payload)
		request := httptest.NewRequest(http.MethodPut, "/"+id, bytes.NewReader(bytesPayload))
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("Bad Requuest", func(t *testing.T) {
		payload := requestentity.Event{
			Description: "Event Updated",
			Deadline:    time.Now().Format(time.RFC3339),
		}
		bytesPayload, _ := json.Marshal(payload)
		request := httptest.NewRequest(http.MethodPut, "/"+id, bytes.NewReader(bytesPayload))
		recorder := httptest.NewRecorder()
		request.Header.Add("Authorization", "Bearer "+token)
		r.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
