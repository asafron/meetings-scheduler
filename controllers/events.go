package controllers

import (
	"github.com/asafron/meetings-scheduler/db"
	"net/http"
	"github.com/asafron/meetings-scheduler/helpers"
)

type (
	EventsController struct {
		dal *db.DAL
	}
)

func NewEventsController(dal *db.DAL) *EventsController {
	return &EventsController{dal : dal}
}

func (mc EventsController) GetEventsForUser(writer http.ResponseWriter, req *http.Request) {
	events := mc.dal.GetEventsForUser(helpers.GetCurrentUser(req).DisplayId)
	m := make(map[string]interface{})
	m["events"] = events
	helpers.JsonResponse(writer, http.StatusOK, &helpers.GeneralResponse{
		Success: true,
		Data: m,
	})
}