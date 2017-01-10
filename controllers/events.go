package controllers

import (
	"github.com/asafron/meetings-scheduler/db"
	"net/http"
	"github.com/asafron/meetings-scheduler/helpers"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/asafron/meetings-scheduler/models"
)

type (
	EventsController struct {
		dal *db.DAL
	}
)

type AddEventRequest struct {
	Name string `json:"name"`
}

type RemoveEventRequest struct {
	DisplayId string        `json:"display_id"`
}

func NewEventsController(dal *db.DAL) *EventsController {
	return &EventsController{dal : dal}
}

func (ec EventsController) GetEventsForUser(writer http.ResponseWriter, req *http.Request) {
	events := ec.dal.GetEventsForUser(helpers.GetCurrentUser(req).DisplayId)
	m := make(map[string]interface{})
	m["events"] = events
	helpers.JsonResponse(writer, http.StatusOK, &helpers.GeneralResponse{
		Success: true,
		Data: m,
	})
}

func (ec EventsController) AddEventForUser(writer http.ResponseWriter, req *http.Request) {
	var request AddEventRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := ec.dal.InsertEvent(request.Name, helpers.GetCurrentUser(req).DisplayId, []models.Slot{}, []models.Meeting{})
	if err != nil {
		log.Fatal(err)
		helpers.JsonResponse(writer, http.StatusOK, &helpers.GeneralResponse{
			Success: false,
			Message: helpers.GeneralErrorInternal.Error(),
		})
		return
	}

	helpers.JsonResponse(writer, http.StatusOK, &helpers.GeneralResponse{
		Success: true,
	})
}

func (ec EventsController) RemoveEvent(writer http.ResponseWriter, req *http.Request) {
	var request RemoveEventRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := ec.dal.RemoveEvent(request.DisplayId)
	if err != nil {
		helpers.JsonResponse(writer, http.StatusOK, &helpers.GeneralResponse{
			Success: false,
			Message: helpers.GeneralErrorInternal.Error(),
		})
		return
	}

	helpers.JsonResponse(writer, http.StatusOK, &helpers.GeneralResponse{
		Success: true,
	})
}