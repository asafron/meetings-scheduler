package models

import "time"

type Slot struct {
	StartTime time.Time `json:"start_time" bson:"start_time"`
	EndTime   time.Time `json:"end_time" bson:"end_time"`
	User      *User     `json:"user" bson:"user"`
	Interval  uint      `json:"interval" bson:"interval"`
}
