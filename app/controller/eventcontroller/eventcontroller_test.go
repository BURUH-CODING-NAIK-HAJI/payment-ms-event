package eventcontroller_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rizface/golang-api-template/app/controller/eventcontroller"
	"github.com/rizface/golang-api-template/app/entity/requestentity"
	"github.com/rizface/golang-api-template/app/entity/responseentity"
	"github.com/rizface/golang-api-template/app/repository/eventrepository"
	"github.com/rizface/golang-api-template/app/service/eventservice"
	"github.com/rizface/golang-api-template/database"
	"github.com/rizface/golang-api-template/database/postgresql"
	"github.com/rizface/golang-api-template/system/router"
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
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGF0YSI6eyJpZCI6ImYzODJkZDUzLTk0MWUtNGY0Zi1hOTFkLWRkNDI1Y2Y3OTZlNCIsIm5hbWUiOiJmYXJpeiIsInVzZXJuYW1lIjoic2FtdWVsIn0sImlzcyI6IlJJWkZBQ0UiLCJzdWIiOiJmYXJpeiIsImV4cCI6MTY2MzQyOTQzNiwibmJmIjoxNjYzMzQzMDM2LCJpYXQiOjE2NjMzNDMwMzZ9.baLSsHdl14xGaN9Ge8L9ThFR_Os5NKzOpkIi3hnx57A"
	repository = eventrepository.New()
	service = eventservice.New(repository, db)
	controller := eventcontroller.New(service)
	eventcontroller.Setup(r, controller)
	dummyData := &requestentity.Event{
		Name:        "Event For Integration Testing",
		Description: "Event For Integration Testing",
		Deadline:    time.Now().Add(5 * 24 * time.Hour).Format(time.RFC3339),
	}
	result, _ := repository.Create(db, uuid.NewString(), dummyData)
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

		eventResponse := new(responseentity.Event)
		recorderResponseBytes, _ := ioutil.ReadAll(recorder.Result().Body)
		json.Unmarshal(recorderResponseBytes, eventResponse)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.NotNil(t, eventResponse)
		assert.Equal(t, id, eventResponse.Id)
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
