package utils

import "errors"

// CustomError hold basic error infomation
type CustomError struct {
	error
	status int
}

// Code return error http code
func (err *CustomError) Code() int {
	return err.status
}

// Err define all custom errors
var (
	NoError     = CustomError{error: nil, status: 200}
	NoCreateErr = CustomError{error: nil, status: 201}
	// post data not correct
	ErrPostDataNotCorrect   = CustomError{error: errors.New("post data not correct"), status: 400}
	ErrResourceAlreadyExist = CustomError{error: errors.New("resource is already exist"), status: 400}
	ErrDataValidateFail     = CustomError{error: errors.New("validation fail"), status: 400}

	// authentication
	ErrForbidden                = CustomError{error: errors.New("operation forbidden"), status: 403}
	ErrAuthorizationRequired    = CustomError{error: errors.New("require authorization"), status: 401}
	ErrUseCredentialsNotCorrect = CustomError{error: errors.New("email or password not correct"), status: 401}
	ErrEmailAlreadyInUse        = CustomError{error: errors.New("email address is already in use"), status: 403}
	ErrUserNotFound             = CustomError{error: errors.New("user not found"), status: 404}

	// resources
	ErrResourceNotFound = CustomError{error: errors.New("resource not found"), status: 404}
	ErrEmptyFolder      = CustomError{error: errors.New("download empty folder is not allowed"), status: 400}

	// sevice
	ErrInternalServerError = CustomError{error: errors.New("internal server error"), status: 500}
	// validation
	ErrValidationForProfileName   = CustomError{error: errors.New("name lenth must greate than 3 and less than 20"), status: 400}
	ErrValidationOverMaxShareDate = CustomError{error: errors.New("share files must less than 90 days"), status: 400}
	ErrTokenIsNotValid            = CustomError{error: errors.New("token is not valid"), status: 400}

	// admin
	ErrDeleteAdminIsNotAllowed = CustomError{error: errors.New("delete admin user is not allowed"), status: 400}
)
