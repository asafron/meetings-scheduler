package controllers

import (
	"github.com/asafron/meetings-scheduler/db"
	"github.com/asafron/meetings-scheduler/auth"
	"encoding/json"
	"strings"
	"net/http"
	"github.com/asafron/meetings-scheduler/helpers"
	"github.com/asafron/meetings-scheduler/config"
	"github.com/asafron/meetings-scheduler/mailer"
	log "github.com/Sirupsen/logrus"
)

type (
	UserController struct {
		dal        *db.DAL
		authorizer *auth.Authenticator
	}
)

func NewUserController(dal *db.DAL, auth *auth.Authenticator) *UserController {
	return &UserController{dal : dal, authorizer : auth}
}

type CreateUserRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type RecoverPasswordRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	Token                string `json:"token"`
}

func (uc UserController) CreateUser(writer http.ResponseWriter, req *http.Request) {
	var createRequest CreateUserRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&createRequest)
	if decodeErr != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//validate request
	if createRequest.Email == "" || createRequest.Password == "" || createRequest.PasswordConfirmation == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "Email, password or password confirmation are missing"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	//password validation
	if createRequest.Password != createRequest.PasswordConfirmation || len(createRequest.Password) < 6 {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "Password is to short or doesn't match the password confirmation"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	//see if this user already exists
	email := strings.ToLower(createRequest.Email)
	// Validate username
	_, err := uc.dal.FindAnyUserByEmail(email)
	if err == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "A user with this email already exists, please try another email"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	} else if err != helpers.AuthenticationErrorLoginUserNotExists {
		if err != nil {
			log.Fatal(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Fatal(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//create the user
	confirmationToken, err := uc.authorizer.Register(email, createRequest.Password, createRequest.FirstName, createRequest.LastName)
	if err != nil {
		log.Fatal(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//send the token
	configWrapper :=config.GetConfigWrapper().GetCurrent()
	subject := "Greetings from PushApps"
	body := "To Get started, please confirm your email address by clicking on the following link:\n"
	body += config.GetConfigWrapper().GetCurrent().DashboardBaseUrl + "/users/confirm?email=" + email + "&token=" + confirmationToken + "\n"
	body+= "Once your registration is completeted, you can login at https://my.pushapps.mobi, you should probably want to visit our documentation at https://docs.pushapps.mobi to see how to start the integration with one of our SDK's.\nRegards,\nThe PushApps team"
	to := []string{email }
	err = mailer.SendMail(to, subject, body, configWrapper.EmailServerFrom,configWrapper.EmailServerUsername,configWrapper.EmailServerPassword,configWrapper.EmailServerAddress,configWrapper.EmailServerPort,configWrapper.EmailServerBcc)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "The user has been registered but the confirmation email sending has failed , please try again or contect the system administrator with this message"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	message := "A confirmation email was sent to " + email + ". Please check your mail and follow the instructions to finish the registration process."
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(helpers.GeneralResponse{Success: true, Message: message})
	return
}

func (uc UserController) Login(writer http.ResponseWriter, req *http.Request) {
	var loginRequest LoginRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&loginRequest)
	if decodeErr != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//validate request
	if loginRequest.Email == "" || loginRequest.Password == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "Email or password are missing"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	email := strings.ToLower(loginRequest.Email)
	err := uc.authorizer.Login(writer, req, email, loginRequest.Password)
	if err != nil {
		switch err {
		case helpers.AuthenticationErrorLoginAlreadyAuthenticated:
			currentUser := helpers.GetCurrentUser(req)
			json.NewEncoder(writer).Encode(&currentUser)
			return
		case helpers.AuthenticationErrorLoginWrongEmailPassword, helpers.AuthenticationErrorLoginUserNotExists:
			writer.WriteHeader(http.StatusBadRequest)
			errorResponse := helpers.GeneralResponse{Message: err.Error()}
			json.NewEncoder(writer).Encode(errorResponse)
			return
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			errorResponse := helpers.GeneralResponse{Message: helpers.GeneralErrorInternal.Error()}
			json.NewEncoder(writer).Encode(errorResponse)
			return
		}
	}
	currentUser := helpers.GetCurrentUser(req)
	json.NewEncoder(writer).Encode(&currentUser)
	return
}

func (uc UserController) Logout(writer http.ResponseWriter, req *http.Request) {
	err := uc.authorizer.Logout(writer, req)
	if err != nil {
		//shouldn't happen
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: err.Error()}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	http.Redirect(writer, req, config.GetConfigWrapper().GetCurrent().DashboardBaseUrl + "#/pages/signin", http.StatusSeeOther)
}
/**
If we got here after the auth middleware then we are authenticated...
 */
func (uc UserController) CheckSession(writer http.ResponseWriter, req *http.Request) {
	currentUser := helpers.GetCurrentUser(req)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(&currentUser)
	return
}

/**
Validate confirmation token and redirect to login page
 */
func (uc UserController) ConfirmUser(writer http.ResponseWriter, req *http.Request) {
	email := req.URL.Query().Get("email")
	token := req.URL.Query().Get("token")
	if email == "" || token == "" {
		http.Error(writer, "Confirmation link is invalid", http.StatusBadRequest)
		return
	}
	email = strings.ToLower(email)
	err := uc.authorizer.ConfirmUser(email, token)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	//redirect to login page
	http.Redirect(writer, req, config.GetConfigWrapper().GetCurrent().DashboardBaseUrl + "#/pages/signin", http.StatusSeeOther)
}

/**
Validate confirmation token and redirect to login page
 */
func (uc UserController) ForgotPassword(writer http.ResponseWriter, req *http.Request) {
	var forgotPasswordRequest ForgotPasswordRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&forgotPasswordRequest)
	if decodeErr != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	email := strings.ToLower(forgotPasswordRequest.Email)
	token, err := uc.authorizer.CreatePasswordRecovery(email)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	//send email
	configWrapper :=config.GetConfigWrapper().GetCurrent()
	subject := "Password recovery"
	body := "We all forget our passwords sometimes... Please follow this link to reset your password:\n"
	body += config.GetConfigWrapper().GetCurrent().DashboardBaseUrl + "/users/recover?email=" + email + "&token=" + token + "\n"
	body += "If you didn't request a new password please contact us as soon as possible."
	to := []string{email }
	err = mailer.SendMail(to, subject, body, configWrapper.EmailServerFrom,configWrapper.EmailServerUsername,configWrapper.EmailServerPassword,configWrapper.EmailServerAddress,configWrapper.EmailServerPort,configWrapper.EmailServerBcc)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "We couldn't send your password recovery email, please try again or contect the system administrator with this message"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	message := "A password recovery email was sent to " + email + ". Please check your mail and follow the instructions to set your password."
	response := helpers.GeneralResponse{Message: message, Success: true}
	json.NewEncoder(writer).Encode(response)
	return
}

func (uc UserController) ValidateRecoverLink(writer http.ResponseWriter, req *http.Request) {
	email := req.URL.Query().Get("email")
	token := req.URL.Query().Get("token")
	if email == "" || token == "" {
		http.Error(writer, "Forgot Password link is invalid", http.StatusBadRequest)
		return
	}
	_, err := uc.dal.FindUserByRecoveryToken(token, email)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	//if link is valid, redirect to forgot password pages
	http.Redirect(writer, req, config.GetConfigWrapper().GetCurrent().DashboardBaseUrl + "#/pages/forgot_password?token=" + token, http.StatusSeeOther)
}

/**
validates the recover password token, saves the new password and redirects to login page
 */
func (uc UserController) RecoverUser(writer http.ResponseWriter, req *http.Request) {
	var recoverPasswordRequest RecoverPasswordRequest
	decoder := json.NewDecoder(req.Body)
	decodeErr := decoder.Decode(&recoverPasswordRequest)
	if decodeErr != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//validate request
	if recoverPasswordRequest.Email == "" || recoverPasswordRequest.Password == "" || recoverPasswordRequest.PasswordConfirmation == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "Email, password or password confirmation are missing"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	//password validation
	if recoverPasswordRequest.Password != recoverPasswordRequest.PasswordConfirmation || len(recoverPasswordRequest.Password) < 6 {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		errorResponse := helpers.GeneralResponse{Message: "Password is to short or doesn't match the password confirmation"}
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	err := uc.authorizer.UpdateUserPasswordFromRecovery(recoverPasswordRequest.Email, recoverPasswordRequest.Token, recoverPasswordRequest.Password)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	//redirect to login page
	http.Redirect(writer, req, config.GetConfigWrapper().GetCurrent().DashboardBaseUrl + "#/pages/signin", http.StatusSeeOther)
}