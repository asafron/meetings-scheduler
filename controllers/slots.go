package controllers

import (
	"github.com/asafron/meetings-scheduler/db"
	"log"
	"time"
	"encoding/json"
	"net/http"
	"github.com/asafron/meetings-scheduler/helpers"
	"github.com/asafron/meetings-scheduler/models"
)

type (
	SlotsController struct {
		dal *db.DAL
	}
)

type AddSlotsToEventRequest struct {
	DisplayId string         `json:"display_id"`
	Slots     []SlotsRequest `json:"slots"`
}

type SlotsRequest struct {
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	User      string `json:"user"`
	Interval  uint   `json:"interval"`
}

type RemoveSlotFromEventRequest struct {
	EventDisplayId string `json:"event_display_id"`
	DisplayId      string `json:"display_id"`
}

func NewSlotsController(dal *db.DAL) *SlotsController {
	return &SlotsController{dal : dal}
}

func (sc SlotsController) AddSlotsToEvent(writer http.ResponseWriter, req *http.Request) {
	var request AddSlotsToEventRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		log.Fatal(decodeErr)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	event, err := sc.dal.GetEventByDisplayId(request.DisplayId)
	if err != nil {
		log.Fatal(err)
		helpers.JsonResponse(writer, http.StatusOK, &helpers.GeneralResponse{
			Success: false,
			Message: helpers.EventsErrorNotFound.Error(),
		})
		return
	}

	slots := event.Slots
	for _, element := range request.Slots {
		sl := &models.Slot{}
		sl.User = element.User
		sl.Interval = element.Interval
		sl.StartTime = time.Unix(element.StartTime, 0)
		sl.EndTime = time.Unix(element.EndTime, 0)
		slots = append(slots, *sl)
	}

	err = sc.dal.UpdateEvent(event.DisplayId, event.Name, event.AdminUser, slots, event.Meetings)
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

func (sc SlotsController) RemoveSlotFromEvent(writer http.ResponseWriter, req *http.Request) {
	var request RemoveSlotFromEventRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		log.Fatal(decodeErr)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := sc.dal.RemoveSlotFromEvent(request.EventDisplayId, request.DisplayId)
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