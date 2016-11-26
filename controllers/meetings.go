package controllers

import (
	"github.com/asafron/meetings-scheduler/db"
	"net/http"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/asafron/meetings-scheduler/helpers"
	"math"
	"fmt"
	"github.com/asafron/meetings-scheduler/models"
)

const DEFAULT_MEETING_INTERVAL = 30

type (
	MeetingsController struct {
		dal *db.DAL
	}
)

type AvailableMeetingTime struct {
	Day       int `json:"day"`
	Month     int `json:"month"`
	Year      int `json:"year"`
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
}

type AvailableMeetingTimeRequest struct {
	Availabilities []AvailableMeetingTime `json:"availabilities"`
}

type ScheduleMeetingRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
}

type MeetingsResponse struct {
	Meetings []models.Meeting `json:"meetings"`
}

func NewMeetingsController(dal *db.DAL) *MeetingsController {
	return &MeetingsController{dal : dal}
}

func (mc MeetingsController) AddAvailableTime(writer http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var requestObj AvailableMeetingTimeRequest
	err := decoder.Decode(&requestObj)
	if err != nil {
		log.Error(err)
		helpers.JsonResponse(writer, http.StatusBadRequest, helpers.ErrorResponse{Message: helpers.RESPONSE_ERROR_MESSAGE_BAD_REQUEST_INPUT_NOT_VALID})
		return
	}

	log.Info("request object created")

	for i := 0; i < len(requestObj.Availabilities); i++ {
		av := requestObj.Availabilities[i]

		// verify availability data
		if av.EndTime < av.StartTime {
			log.Error("end time is before start time, skipping...")
			continue
		}

		if math.Remainder(float64(av.StartTime), DEFAULT_MEETING_INTERVAL) != 0 {
			log.Error("start time is not acceptable due to default meeting interval, skipping...")
			continue
		}

		if math.Remainder(float64(av.EndTime), DEFAULT_MEETING_INTERVAL) != 0 {
			log.Error("start time is not acceptable due to default meeting interval, skipping...")
			continue
		}

		// dividing it to 30 min meetings
		for j := 0; j < (av.EndTime - av.StartTime) / DEFAULT_MEETING_INTERVAL; j++ {
			meetingStartTime := av.StartTime + (j * DEFAULT_MEETING_INTERVAL)
			meetingEndTime := meetingStartTime + DEFAULT_MEETING_INTERVAL
			log.Info(fmt.Sprintf("trying to create a meeting: %d/%d/%d %d-%d", av.Day, av.Month, av.Year, meetingStartTime, meetingEndTime))
			meetingErr := mc.dal.InsertAvailableMeetingTime(av.Day, av.Month, av.Year, meetingStartTime, meetingEndTime)
			if meetingErr != nil {
				log.Error(meetingErr)
				helpers.JsonResponse(writer, http.StatusBadRequest, helpers.ErrorResponse{Message: fmt.Sprintf("%s", meetingErr.Error())})
				return
			}
		}

	}

	helpers.JsonResponse(writer, http.StatusOK, helpers.MinimalResponse{Success:true})
	return
}

func (mc MeetingsController) ScheduleMeeting(writer http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var requestObj ScheduleMeetingRequest
	err := decoder.Decode(&requestObj)
	if err != nil {
		log.Error(err)
		helpers.JsonResponse(writer, http.StatusBadRequest, helpers.ErrorResponse{Message: helpers.RESPONSE_ERROR_MESSAGE_BAD_REQUEST_INPUT_NOT_VALID})
		return
	}

	log.Info("request object created")

	// validation
	if len(requestObj.Name) == 0 || len(requestObj.Email) == 0 || len(requestObj.Phone) == 0 {
		log.Error("some of the user details are missing")
		helpers.JsonResponse(writer, http.StatusBadRequest, helpers.ErrorResponse{Message: "some of the user details are missing"})
		return
	}

	meetingErr := mc.dal.UpdateMeetingDetails(requestObj.Day, requestObj.Month, requestObj.Year, requestObj.StartTime,
							requestObj.EndTime, requestObj.Name, requestObj.Email, requestObj.Phone)
	if meetingErr != nil {
		log.Error(meetingErr)
		helpers.JsonResponse(writer, http.StatusBadRequest, helpers.ErrorResponse{Message: fmt.Sprintf("%s", meetingErr.Error())})
		return
	}

	helpers.JsonResponse(writer, http.StatusOK, helpers.MinimalResponse{Success:true})
	return
}

func (mc MeetingsController) GetAllMeetings(writer http.ResponseWriter, req *http.Request) {
	meetingResponse := MeetingsResponse{}

	meetings := mc.dal.GetAllMeetings()
	if meetings != nil {
		meetingResponse.Meetings = meetings
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(&meetingResponse)
}

func (mc MeetingsController) GetAvailableMeetings(writer http.ResponseWriter, req *http.Request) {
	meetingResponse := MeetingsResponse{}

	meetings := mc.dal.GetAvailableMeetings()
	if meetings != nil {
		meetingResponse.Meetings = meetings
	}

	// TODO: need to order the dates...

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(&meetingResponse)
}

func (mc MeetingsController) GetScheduledMeetings(writer http.ResponseWriter, req *http.Request) {
	meetingResponse := MeetingsResponse{}

	meetings := mc.dal.GetScheduledMeetings()
	if meetings != nil {
		meetingResponse.Meetings = meetings
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(&meetingResponse)
}