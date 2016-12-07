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
	return fmt.Sprintf("%d/%d/%d %d:%d-%d:%d", meeting.Day, meeting.Month, meeting.Year,
		int(math.Floor(float64(meeting.StartTime)/float64(MINUTES_IN_HOUR))),
		int(math.Remainder(float64(meeting.StartTime), MINUTES_IN_HOUR)),
		int(math.Floor(float64(meeting.EndTime)/float64(MINUTES_IN_HOUR))),
		int(math.Remainder(float64(meeting.EndTime), MINUTES_IN_HOUR)))
}