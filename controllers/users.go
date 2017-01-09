package controllers

import (
	"github.com/asafron/meetings-scheduler/db"
	"github.com/asafron/meetings-scheduler/auth"
)

type (
	UserController struct {
		dal        *db.DAL
		authorizer *auth.Authenticator
	}
)