package userhandler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"user_service/src/internal/core/user"
	userservice "user_service/src/internal/usecase"
	errorhandling "user_service/src/pkg/error_handling"
	pkgresponse "user_service/src/pkg/response"
)

type UserHandler struct {
	userService userservice.UserService
}

func NewUserHandler(usecase userservice.UserService) UserHandler {
	return UserHandler{
		userService: usecase,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var newUser user.UserRegister
	var createdUser user.UserResponse
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		errorhandling.HandleError(w, "Wrong Format Data", http.StatusBadRequest)
		return
	}

	createdUser, err := u.userService.RegisterUser(newUser)
	if err != nil {
		errorhandling.HandleError(w, "Unable to Register User", http.StatusInternalServerError)
		return
	}
	// createdUser = registeredUser
	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "User Registered Successfully",
		Data:    createdUser,
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser user.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		errorhandling.HandleError(w, "Wrong Format Data", http.StatusBadRequest)
		return
	}

	loginResponse, err := u.userService.LoginUser(loginUser)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	atCookie := http.Cookie{
		Name:     "at",
		Value:    loginResponse.TokenString,
		Expires:  loginResponse.TokenExpire,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	sessCookie := http.Cookie{
		Name:     "sess",
		Value:    loginResponse.Session.Id.String(),
		Expires:  loginResponse.Session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &sessCookie)

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Successful Login",
		Data: map[string]interface{}{
			"username": loginResponse.FounUser.Username,
			"user_id":  loginResponse.FounUser.Uid,
		},
	}
	w.Header().Set("x-user", loginResponse.FounUser.Username)
	w.Header().Set("x-userId", strconv.Itoa(loginResponse.FounUser.Uid))
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		errorhandling.HandleError(w, "User Not Found in Context", http.StatusUnauthorized)
		return
	}

	registeredUser, err := u.userService.GetUserByID(userId)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "User Profile Retrieved Successfully",
		Data:    registeredUser,
	}
	w.Header().Set("x-user", registeredUser.Username)
	w.Header().Set("x-userId", strconv.Itoa(registeredUser.Uid))
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sess")
	if err != nil {
		errorhandling.HandleError(w, "Session Cookie Not Found", http.StatusUnauthorized)
		return
	}

	tokenString, expireTime, err := u.userService.GetJwtFromSession(cookie.Value)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	atCookie := http.Cookie{
		Name:     "at",
		Value:    tokenString,
		Expires:  expireTime,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &atCookie)

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Token Refreshed Successfully",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		errorhandling.HandleError(w, "User Not Found in Context", http.StatusUnauthorized)
		return
	}

	err := u.userService.LogoutUser(userId)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	atCookie := http.Cookie{
		Name:     "at",
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &atCookie)

	sessCookie := http.Cookie{
		Name:     "sess",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(w, &sessCookie)

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Successful Logout",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	allUsers, err := u.userService.GetAllUsers()
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Here is the data of all users",
		Data:    allUsers,
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}
