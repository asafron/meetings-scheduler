package helpers

import "errors"

var (
	AuthenticationErrorRegisterNoEmail = MakeError("Email is empty")
	AuthenticationErrorRegisterNoPassword = MakeError("Password is empty")
	AuthenticationErrorRegisterPasswordNotValid = MakeError("Password not acceptable")
	AuthenticationErrorRegisterUserCreationFailed = MakeError("Could not create new user")

	AuthenticationErrorLoginAlreadyAuthenticated = MakeError("User is already authenticated")
	AuthenticationErrorLoginWrongEmailPassword = MakeError("Wrong email and/or password")
	AuthenticationErrorLoginUserNotExists = MakeError("Email doesn't exist")

	AuthenticationErrorAuthorizeNewSession = MakeError("new authorization session")
	AuthenticationErrorAuthorizeUserNotLoggedIn = MakeError("user not logged in")

	AuthenticationErrorConfirmationTokenNotValid = MakeError("Confirmation token is not valid")

	GeneralErrorInternal = MakeError("An error has aoccured")
)

func MakeError(msg string) error {
	return errors.New(msg)
}