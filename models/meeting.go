package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Meeting struct {
	Id                     bson.ObjectId `json:"id" bson:"_id"`
	DisplayId              string        `json:"display_id" bson:"display_id"`
	StartTime              time.Time     `json:"start_time" bson:"start_time"`
	EndTime                time.Time     `json:"end_time" bson:"end_time"`
	ClientDetails          interface{}   `json:"client_details" bson:"client_details"`
	UserId                 string        `json:"user_id" bson:"user_id"`
	CreatedAt              time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at" bson:"updated_at"`
}