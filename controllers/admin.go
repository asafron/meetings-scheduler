package controllers

import (
	"github.com/asafron/meetings-scheduler/db"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"encoding/json"
	"github.com/asafron/meetings-scheduler/helpers"
	"strings"
	"github.com/asafron/meetings-scheduler/models"
	"github.com/bradfitz/slice"
)

type (
	AdminController struct {
		dal *db.DAL
	}
)

type ManagerGetAllMeetingsRequest struct {
	Password string `json:"auth"`
}

type ManagerMeetingsResponse struct {
	Meetings []models.Meeting `json:"meetings"`
}

func NewAdminController(dal *db.DAL) *AdminController {
	return &AdminController{dal : dal}
}

func (ac AdminController) ManagerGetAllMeetings(writer http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var requestObj ManagerGetAllMeetingsRequest
	err := decoder.Decode(&requestObj)
	if err != nil {
		log.Error(err)
		helpers.JsonResponse(writer, http.StatusBadRequest, helpers.ErrorResponse{Message: helpers.RESPONSE_ERROR_MESSAGE_BAD_REQUEST_INPUT_NOT_VALID})
		return
	}

	log.Info("request object created")

	if (strings.Compare(requestObj.Password, "2901hrsbvnuyh") != 0) {
		log.Error("Auth failed")
		helpers.JsonResponse(writer, http.StatusInternalServerError, helpers.ErrorResponse{Message: helpers.RESPONSE_ERROR_MESSAGE_INTERNAL_SERVER_ERROR})
		return
	}

	meetingResponse := ManagerMeetingsResponse{}

	meetings := ac.dal.GetAllMeetings()
	if meetings != nil {
		meetingResponse.Meetings = meetings
	}

	slice.Sort(meetings[:], func(i, j int) bool {
		if meetings[i].Year != meetings[j].Year {
			return meetings[i].Year < meetings[j].Year
		} else if meetings[i].Month != meetings[j].Month {
			return meetings[i].Month < meetings[j].Month
		} else if meetings[i].Day != meetings[j].Day {
			return meetings[i].Day < meetings[j].Day
		} else if meetings[i].StartTime != meetings[j].StartTime {
			return meetings[i].StartTime < meetings[j].StartTime
		} else if meetings[i].EndTime != meetings[j].EndTime {
			return meetings[i].StartTime < meetings[j].StartTime
		} else {
			return true // arbitrary
		}
	})

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(&meetingResponse)
}