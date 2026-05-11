package controllers

import (
	"errors"
	"forum/src/models"
	"forum/src/utils"
	"forum/src/views"
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
)

func userLogin(data models.ResponseStruct) {
	if data.User.LoggedIn {
		Index(*data.SetErrorConsume(models.ErrorAlreadyLoggedIn))
		return
	}
	switch data.Request.Method {
	case http.MethodPost:
		attemptLogin(data)
	case http.MethodGet:
		views.UserLogin(&data)
	default:
		data.SetErrorConsume(models.ErrorMethodNotAllowed).WriteResponse()
	}
}

func userLogout(data models.ResponseStruct) {
	GuestUser := models.GetGuestUser()
	cookie, err := data.Request.Cookie("__Host-FRMSessionID")
	if errors.Is(err, http.ErrNoCookie) {
		data.SetUser(GuestUser).SetErrorConsume(models.ErrorAlreadyLoggedOut)
		views.ErrorView(&data)
		// data.SetView("error_view").WriteResponse()
		return
	}
	user, err := models.GetUserBySession(cookie.Value)
	if err != nil {
		data.SetUser(user)
		data.SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	err = user.SetUserSession("")
	if err != nil {
		data.SetUser(user)
		data.SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	http.SetCookie(data.Response, nullifyCookie(cookie))
	data.SetUser(GuestUser)
	data.Message.Has = true
	data.Message.Type = "Success"
	data.Message.Content = "Logout successful"
	views.UserLogout(&data)
}

func nullifyCookie(cookie *http.Cookie) *http.Cookie {
	cookie = &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
	return cookie
}

func attemptLogin(data models.ResponseStruct) {
	var email string
	var password string
	var err error
	if len(data.Request.Form.Get("email")) != 0 {
		email = data.Request.Form.Get("email")
	} else {
		data.SetErrorConsume(models.ErrorEmailFieldEmpty)
		views.UserLogin(&data)
		return
	}
	if len(data.Request.Form.Get("password")) != 0 {
		password = data.Request.Form.Get("password")
	} else {
		data.SetErrorConsume(models.ErrorPasswordFieldEmpty)
		views.UserLogin(&data)
		return
	}
	err = Auth(email, password)
	if err != nil {
		data.User = models.GetGuestUser()
		if !errors.Is(err, models.ErrorWrongPassword) && !errors.Is(err, models.ErrorNotRegistered) {
			(&models.Error{}).Consume(err).LogError()
			data.SetErrorConsume(models.ErrorInternalServerError)
			views.ErrorView(&data)
			return
		}
		data.SetErrorConsume(err)
		views.UserLogin(&data)
		return
	}
	sessionValue, err := uuid.NewV4()
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	data.User, err = models.GetUserByEmail(email)
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	data.User.LoggedIn = true
	err = data.User.SetUserSession(sessionValue.String())
	if err != nil {
		data.User = models.GetGuestUser()
		data.SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	cookie := &http.Cookie{
		Name:     "__Host-FRMSessionID",
		Value:    sessionValue.String(),
		Path:     "/",
		Expires:  time.Now().Add(24*time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSite(http.SameSiteStrictMode),
	}
	http.SetCookie(data.Response, cookie)
	http.Redirect(data.Response, data.Request, "/", http.StatusSeeOther)
}

func userRegister(data models.ResponseStruct) {
	if data.User.LoggedIn {
		Index(*data.SetErrorConsume(models.ErrorAlreadyLoggedIn))
		return
	}
	switch data.Request.Method {
	case http.MethodGet:
		views.UserRegister(&data)
		return
	case http.MethodPost:
		attemptRegister(data)
		return
	default:
		data.SetErrorConsume(models.ErrorMethodNotAllowed).WriteResponse()
	}
}

func attemptRegister(data models.ResponseStruct) {
	var err error
	if len(data.Request.FormValue("username")) == 0 ||
		len(data.Request.FormValue("email")) == 0 ||
		len(data.Request.FormValue("password1")) == 0 ||
		len(data.Request.FormValue("password2")) == 0 {
		data.SetErrorConsume(models.ErrorBadRequest)
		views.UserRegister(&data)
		return
	}
	data.User.Username = data.Request.FormValue("username")
	data.User.Email = data.Request.FormValue("email")
	if err = data.User.ValidateUser(); err != nil {
		data.SetUser(data.User).SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	if models.IsEmailRegistered(data.User.Email) {
		data.SetUser(data.User).SetErrorConsume(models.ErrorEmailIsRegistered)
		views.UserRegister(&data)
		return
	}
	if !models.IsUniqueUsername(data.User.Username) {
		data.SetUser(data.User).SetErrorConsume(models.ErrorUsernameTaken)
		views.UserRegister(&data)
		return
	}
	if !CompareRegistrationPasswords(data.Request.FormValue("password1"), data.Request.FormValue("password2")) {
		data.SetUser(data.User).SetErrorConsume(models.ErrorPasswordMismatch)
		views.UserRegister(&data)
		return
	}
	password := data.Request.FormValue("password1")
	if err = validatePasswordStrength(password); err != nil {
		data.SetUser(data.User).SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	data.User.Hash, err = utils.HashPassword(password)
	if err != nil {
		data.SetUser(data.User).SetErrorConsume(err)
		views.UserRegister(&data)
		return
	}
	if err = data.User.Add(); err != nil {
		data.SetUser(data.User).SetErrorConsume(models.ErrorInvalidUser)
		views.UserRegister(&data)
		return
	}
	data.SetUser(models.GetGuestUser())
	data.Message.Content = "Registration was successful"
	data.Message.Type = "Success"
	data.Message.Has = true
	views.UserLogin(&data)
}

// Strong password validation. Makes sure the password is in between 10-16
// characters and includes letters, numbers and/or punctuation symbols
func validatePasswordStrength(password string) error {
	unameMask := regexp.MustCompile(`^[[:punct:][:alnum:]]{10,16}$`)
	if !unameMask.MatchString(password) {
		return models.ErrorWeakPassword
	}
	return nil
}

func showUserPosts(data models.ResponseStruct) {
	posts, err := data.User.GetPosts()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = posts[i].GetReactionsByUserId(data.User.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	data.Posts = posts
	views.PostsView(&data)
}

func showUserLikedPosts(data models.ResponseStruct) {
	var err error
	var posts models.Posts
	posts, err = data.User.GetLikedPosts()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = posts[i].GetReactionsByUserId(data.User.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	data.Posts = posts
	views.PostsView(&data)
}

func showUserView(data models.ResponseStruct) {
	views.UserView(&data)
}

func showUserActivity(data models.ResponseStruct) {
	err := data.User.GetActivity()
	if err != nil {
		(&models.Error{}).Consume(models.ErrorInternalServerError).LogAndRespondError(data.Response, data.User)
		return
	}
	views.UserActivity(&data)
}
