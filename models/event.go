package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Event struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	DisplayId string        `json:"display_id" bson:"display_id"`
	AdminUser string        `json:"admin_user" bson:"admin_user"`
	Slots     []Slot        `json:"slots" bson:"slots"`
	Name      string        `json:"name" bson:"name"`
	Meetings  []Meeting     `json:"meetings" bson:"meetings"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}