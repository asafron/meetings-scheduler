package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Guest struct {
	Id        bson.ObjectId     `json:"id" bson:"_id"`
	DisplayId string            `json:"display_id" bson:"display_id"`
	FirstName string            `json:"first_name" bson:"first_name"`
	LastName  string            `json:"last_name" bson:"last_name"`
	Email     string            `json:"email" bson:"email"`
	Phone     string            `json:"phone" bson:"phone"`
	Details   map[string]string `json:"details" bson:"details"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
}