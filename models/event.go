package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	DisplayId string        `json:"display_id" bson:"display_id"`
	AdminUser *User         `json:"admin_user" bson:"admin_user"`
	Slots     []Slot        `json:"slots" bson:"slots"`
	Name      string        `json:"name" bson:"name"`
}