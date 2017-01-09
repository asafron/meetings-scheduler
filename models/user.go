package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id                      bson.ObjectId               `json:"id" bson:"_id"`
	DisplayId               string                      `json:"display_id" bson:"display_id"`
	FirstName               string                      `json:"first_name" bson:"first_name"`
	LastName                string                      `json:"last_name" bson:"last_name"`
	Email                   string                      `json:"email" bson:"email"`
	Hash                    []byte                      `json:"-" bson:"hash"`
	ConfirmationToken       string                      `json:"-" bson:"confirmation_token"`
	ConfirmationTokenStatus ConfirmationTokenStatusType `json:"-" bson:"confirmation_token_status"`

	Confirmed               bool                        `json:"-" bson:"confirmed"`
	Status                  UserStatusType              `json:"-" bson:"status"`

	RecoverToken		string                      `json:"-" bson:"recovery_token"`
	RecoverTokenExpiry      time.Time                   `json:"-" bson:"recovery_token_expiry"`
	RecoverTokenStatus	RecoverTokenStatusType      `json:"-" bson:"recovery_token_status"`

}

type UserStatusType string
type ConfirmationTokenStatusType string
type RecoverTokenStatusType string

const (
	USER_NOT_CONFIRMED UserStatusType = "not_confirmed"
	USER_CONFIRMED UserStatusType = "confirmed"
)

const (
	CONFIRMATION_TOKEN_VALID ConfirmationTokenStatusType = "valid"
	CONFIRMATION_TOKEN_INVALID ConfirmationTokenStatusType = "invalid"
)

const (
	RECOVER_TOKEN_VALID RecoverTokenStatusType = "valid"
	RECOVER_TOKEN_INVALID RecoverTokenStatusType = "invalid"
)
