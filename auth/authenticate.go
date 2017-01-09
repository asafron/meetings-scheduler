package auth

import (
	"golang.org/x/crypto/bcrypt"
	"os/exec"
	"strings"
	"net/http"
	"github.com/gorilla/sessions"
	"time"
	"github.com/asafron/meetings-scheduler/db"
	"github.com/asafron/meetings-scheduler/helpers"
	"github.com/asafron/meetings-scheduler/models"
	"github.com/asafron/meetings-scheduler/config"
)

type Authenticator struct {
	dal         *db.DAL
	cookieJar   *sessions.CookieStore
}

func NewAuthenticator(dal *db.DAL, cookieKey string) *Authenticator {
	a := Authenticator{}
	a.dal = dal
	a.cookieJar = sessions.NewCookieStore([]byte(cookieKey))
	return &a
}


func (a *Authenticator) Register(role string, email string , password string, firstName string , lastName string, company string, website string) (string, error) {
	if email == "" {
		return "", helpers.AuthenticationErrorRegisterNoEmail
	}
	if password == "" {
		return "", helpers.AuthenticationErrorRegisterNoPassword
	}

	// Generate and save hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", helpers.AuthenticationErrorRegisterPasswordNotValid
	}

	// Create confirmation token
	token, tokenErr := exec.Command("uuidgen").Output()
	if tokenErr !=nil {
		return "", tokenErr
	}
	confirmationToken := strings.TrimRight(strings.ToLower(string(token)), "\n")
	err = a.dal.InsertUser(email, hash, firstName, lastName, company, website, confirmationToken)
	if err != nil {
		return "", helpers.AuthenticationErrorRegisterUserCreationFailed
	}
	return confirmationToken, nil
}

func (a *Authenticator) Login(rw http.ResponseWriter, req *http.Request, email string, password string) error {
	session, _ := a.cookieJar.Get(req, "auth")
	if session.Values["email"] != nil {
		// Set the current user
		username:= (session.Values["email"]).(string)
		user, err := a.dal.FindActiveUserByEmail(username)
		if err == nil {
			helpers.SetCurrentUser(req,*user)
		}
		return helpers.AuthenticationErrorLoginAlreadyAuthenticated
	}
	// Try to find the user, to see if it already logged in...
	user, err := a.dal.FindActiveUserByEmail(email)
	if  err == nil {
		verify := bcrypt.CompareHashAndPassword(user.Hash, []byte(password))
		if verify != nil {
			return helpers.AuthenticationErrorLoginWrongEmailPassword
		}
	} else {
		return helpers.AuthenticationErrorLoginUserNotExists
	}
	session.Values["email"] = email
	session.Save(req, rw)
	helpers.SetCurrentUser(req,*user)
	return nil
}

func (a *Authenticator) Authorize(rw http.ResponseWriter, req *http.Request) (*models.User, error) {
	var user *models.User
	authSession, err := a.cookieJar.Get(req, "auth")
	if err != nil {
		return user, helpers.AuthenticationErrorAuthorizeNewSession
	}
	username := authSession.Values["email"]
	if !authSession.IsNew && username != nil {
		user, err = a.dal.FindActiveUserByEmail(username.(string))
		if err == helpers.AuthenticationErrorLoginUserNotExists {
			authSession.Options.MaxAge = -1 // kill the cookie
			authSession.Save(req, rw)
			return user, helpers.AuthenticationErrorLoginUserNotExists
		} else if err != nil {
			return user, helpers.GeneralErrorInternal
		}
	}
	if username == nil {
		return user,helpers.AuthenticationErrorAuthorizeUserNotLoggedIn
	}

	return user,nil
}

func (a *Authenticator) Logout(rw http.ResponseWriter, req *http.Request) error {
	session, _ := a.cookieJar.Get(req, "auth")
	defer session.Save(req, rw)
	session.Options.MaxAge = -1 // kill the cookie
	return nil
}

func (auth *Authenticator) ConfirmUser(email string, confirmationToken string) error {
	// Locate the user
	user, err := auth.dal.FindUserByConfirmationToken(confirmationToken, email)
	if err != nil {
		return helpers.AuthenticationErrorConfirmationTokenNotValid
	}
	// If all is OK, confirm the user
	err = auth.dal.UpdateUserConfirmation(user.Id, models.USER_CONFIRMED, models.CONFIRMATION_TOKEN_INVALID, true)
	if err != nil {
		return helpers.GeneralErrorInternal
	}
	return nil
}

func (auth *Authenticator) CreatePasswordRecovery(email string) (string,error) {
	// Locate the user
	user, err := auth.dal.FindActiveUserByEmail(email)
	if err != nil {
		return "", err
	}
	// If all is OK, create password recovery
	token, err := helpers.CreateToken()
	if err != nil {
		return "", err
	}
	// Update user, set expiry to 24 hours from now
	err = auth.dal.UpdateUserRecovery(user.Id, token, models.RECOVER_TOKEN_VALID, time.Now().UTC().Add(time.Hour * 24) )
	if err!=nil {
		return "", err
	}
	return token, nil
}

func (auth *Authenticator) UpdateUserPasswordFromRecovery(email string, token string, password string ) error {
	//get the user
	user, err := auth.dal.FindUserByRecoveryToken(token,email)
	if err != nil {
		return err
	}
	//generate hash
	// Generate and save hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.AuthenticationErrorRegisterPasswordNotValid
	}
	err = auth.dal.UpdateUserPassword(user.Id,hash, models.RECOVER_TOKEN_INVALID, time.Now().UTC())
	if err!=nil {
		return err
	}
	return nil
}

func (auth *Authenticator) AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		var err error=nil
		user, err = auth.Authorize(w, r)
		if err != nil {
			http.Redirect(w, r, config.GetConfigWrapper().GetCurrent().DashboardLoggedInUrl, http.StatusSeeOther)
			return
		}
		helpers.SetCurrentUser(r,*user)
		h.ServeHTTP(w, r)
	})
}