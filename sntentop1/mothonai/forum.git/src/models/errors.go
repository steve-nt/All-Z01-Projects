package models

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	ErrorNotRegistered           = errors.New("Email is not registered")
	ErrorEmailIsRegistered       = errors.New("Email is already registered")
	ErrorEmailFieldEmpty         = errors.New("Email field can't be empty")
	ErrorPasswordFieldEmpty      = errors.New("Password field can't be empty")
	ErrorInvalidUsername         = errors.New("Username is invalid")
	ErrorUsernameTaken           = errors.New("Username is taken")
	ErrorInvalidUser             = errors.New("Invalid user")
	ErrorWrongPassword           = errors.New("Wrong password")
	ErrorWeakPassword            = errors.New("Weak password. Use lower and upper case letters, symbols and number. Length must be between 10-16 characters.")
	ErrorPasswordMismatch        = errors.New("Password mismatch")
	ErrorAlreadyLoggedIn         = errors.New("Already logged in.")
	ErrorAlreadyLoggedOut        = errors.New("Already logged out.")
	ErrorNotFound                = errors.New("Not found")
	ErrorPostEmptyId             = errors.New("Post ID can't be empty.")
	ErrorInvalidPostId           = errors.New("Invalid post ID")
	ErrorInvalidCommentId        = errors.New("Invalid comment ID")
	ErrorInvalidCategoryId       = errors.New("Invalid category ID")
	ErrorPostBodyEmpty           = errors.New("Post body can't be empty.")
	ErrorPostTitleEmpty          = errors.New("Post title can't be empty.")
	ErrorPostHasNoCategory       = errors.New("Post category can't be empty.")
	ErrorPostPermissionDenied    = errors.New("You must be logged in to create a post.")
	ErrorUserPermissionDenied    = errors.New("You must be logged in.")
	ErrorCommentEmpty            = errors.New("Comment can't be empty.")
	ErrorCommentTooLong          = errors.New("Comment is too long.")
	ErrorCommentPermissionDenied = errors.New("You must be logged in to create a comment.")
	ErrorCommentEmptyId          = errors.New("Comment ID can't be empty.")
	ErrorCategoryEmptyId         = errors.New("Category ID can't be empty.")
	ErrorCategoryAlreadyExists   = errors.New("Category already exists.")
	ErrorCategoryNameEmpty       = errors.New("Category name can't be empty.")
	ErrorCategoryNameTooLong     = errors.New("Category name is too long. Use less than 128 characters.")
	ErrorUnauthorizedAction      = errors.New("Unauthorized action.")
	ErrorMethodNotAllowed        = errors.New("Method not allowed.")
	ErrorBadRequest              = errors.New("Bad request.")
	ErrorInternalServerError     = errors.New("Internal server error.")
	ErrorUnknownAction           = errors.New("Unknown action requested.")
	ErrorImageTooBig             = errors.New("Image is too big. Maximum size is 20MB.")
	ErrorInvalidImageType        = errors.New("Invalid image type. Allowed types: JPEG, PNG, GIF.")
	ErrorFailedToGetCaller       = errors.New("Failed to get caller information")
	ErrorNoRows                  = errors.New("No rows")
	ErrorEmailNotFoundForOAuth   = errors.New("Could not associate email with given OAuth provider. Try to login with a password or another provider.")
	ErrorAccessToken			 = errors.New("Failed to retrieve access token")
	ErrorCookieNotFound			 = errors.New("State cookie not found")
	ErrorInvalidOAuthState		 = errors.New("Invalid OAuth state")
	ErrorContentNotFound         = errors.New("Content not found. It doesn't exist or it may have been deleted")
)

type Error struct {
	Has        bool
	StatusCode int
	Message    string
	Error      error
}

type ErrorIface interface {
	LogError()
	RespondError(http.ResponseWriter)
	Consume(error) Error
}

// Converts error type to *Error
func (e *Error) Consume(err error) *Error {
	e.Message = strings.ReplaceAll(err.Error(), "\n", ": ")
	e.Error = err
	e.Has = true
	switch err {
	case ErrorNotFound,
		ErrorContentNotFound:
		e.StatusCode = http.StatusNotFound
	case
		ErrorUnauthorizedAction,
		ErrorUserPermissionDenied,
		ErrorCommentPermissionDenied,
		ErrorPostPermissionDenied:
		e.StatusCode = http.StatusForbidden
	case ErrorMethodNotAllowed:
		e.StatusCode = http.StatusMethodNotAllowed
	case ErrorPostEmptyId,
		ErrorInvalidPostId,
		ErrorInvalidCommentId,
		ErrorInvalidCategoryId,
		ErrorPostBodyEmpty,
		ErrorPostTitleEmpty,
		ErrorPostHasNoCategory,
		ErrorCommentEmpty,
		ErrorCommentTooLong,
		ErrorCommentEmptyId,
		ErrorCategoryEmptyId,
		ErrorCategoryNameEmpty,
		ErrorCategoryNameTooLong,
		ErrorEmailFieldEmpty,
		ErrorPasswordFieldEmpty,
		ErrorBadRequest:
		e.StatusCode = http.StatusBadRequest
	case ErrorInternalServerError:
		e.StatusCode = http.StatusInternalServerError
	}
	return e
}

// Logs *Error to terminal
func (e *Error) LogError() {
	log.Printf("Error: %s", e.Message)
}

// Responds *Error to user with the error_view template
func (e *Error) RespondError(res http.ResponseWriter, user User) {
	data := ResponseStruct{}
	data.Init().SetResponse(res)
	data.SetError(*e).SetUser(user).SetView("error_view").WriteResponse()
}

// Logs *Error to terminal and responds to user with error_view template
func (e *Error) LogAndRespondError(res http.ResponseWriter, user User) {
	e.LogError()
	e.RespondError(res, user)
}
