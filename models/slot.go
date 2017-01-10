package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Slot struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	DisplayId string        `json:"display_id" bson:"display_id"`
	StartTime time.Time     `json:"start_time" bson:"start_time"`
	EndTime   time.Time     `json:"end_time" bson:"end_time"`
	User      string        `json:"user" bson:"user"`
	Interval  uint          `json:"interval" bson:"interval"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
