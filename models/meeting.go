package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	"math"
)

const MINUTES_IN_HOUR = 60

type Meeting struct {
	Id             bson.ObjectId `json:"id" bson:"_id"`
	Day            int           `json:"day" bson:"day"`
	Month          int           `json:"month" bson:"month"`
	Year           int           `json:"year" bson:"year"`
	StartTime      int           `json:"start_time" bson:"start_time"`
	EndTime        int           `json:"end_time" bson:"end_time"`
	Representative string        `json:"representative" bson:"representative"`
	UserName       string        `json:"user_name" bson:"user_name"`
	UserEmail      string        `json:"user_email" bson:"user_email"`
	UserPhone      string        `json:"user_phone" bson:"user_phone"`
	UserSchool     string        `json:"user_school" bson:"user_school"`
	UserIdNumber   string        `json:"user_id_number" bson:"user_id_number"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
}

func (meeting Meeting) GetMeetingDateAsString() string {
	return GetMeetingDateAsString(meeting.Day, meeting.Month, meeting.Year, meeting.StartTime, meeting.EndTime)
}

func GetMeetingDateAsString (day, month, year, startTime, endTime int) string {
	dayStr := fmt.Sprintf("%d", day)
	if len(dayStr) == 1 {
		dayStr = "0" + dayStr
	}

	monthStr := fmt.Sprintf("%d", month)
	if len(monthStr) == 1 {
		monthStr = "0" + monthStr
	}

	yearStr := fmt.Sprintf("%d", year)

	startHour := fmt.Sprintf("%d", int(math.Floor(float64(startTime)/float64(MINUTES_IN_HOUR))))
	if len(startHour) == 1 {
		startHour = "0" + startHour
	}
	startMin := fmt.Sprintf("%d", int(math.Remainder(float64(startTime), MINUTES_IN_HOUR)))
	if len(startMin) == 1 {
		startMin = "0" + startMin
	}

	endHour := fmt.Sprintf("%d", int(math.Floor(float64(endTime)/float64(MINUTES_IN_HOUR))))
	if len(endHour) == 1 {
		endHour = "0" + endHour
	}
	endMin := fmt.Sprintf("%d", int(math.Remainder(float64(endTime), MINUTES_IN_HOUR)))
	if len(endMin) == 1 {
		endMin = "0" + endMin
	}

	return fmt.Sprintf("%d/%d/%d %d:%d-%d:%d", dayStr, monthStr, yearStr, startHour, startMin, endHour, endMin)
}