package helpers

import (
	"github.com/gorilla/context"
	"net/http"
	"github.com/asafron/meetings-scheduler/models"
)

type key int

const currentUserKey key = 1

// GetMyKey returns a value for this package from the request values.
func GetCurrentUser(r *http.Request) models.User {
	rv , ok := context.GetOk(r, currentUserKey)
	if !ok {
		return rv.(models.User)
	}
	return rv.(models.User)
}

// SetMyKey sets a value for this package in the request values.
func SetCurrentUser(r *http.Request, val models.User) {
	context.Set(r, currentUserKey, val)
}