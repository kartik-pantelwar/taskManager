package userhandler

import (
	"encoding/json"
	"net/http"
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

// todo
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
	// var requestUser user.User
	var loginUser user.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	loginResponse, err := u.userService.LoginUser(loginUser)
	if err != nil {
		errorhandling.HandleError(w, "Unable to Login", http.StatusBadRequest)
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
	}
	w.Header().Set("x-user", loginResponse.FounUser.Username)
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	//user ID is fetched by using Context value. We passed context in the Authenticate middleware, which picks the user value and store it is context value, so we can get user ID in any route using context, after using authenticate middleware
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	registeredUser, err := u.userService.GetUserByID(userId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-user", registeredUser.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(registeredUser)
}

func (u *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sess")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	tokenString, expireTime, err := u.userService.GetJwtFromSession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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

}

func (u *UserHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode((map[string]interface{}{"Error": "user not found in context"}))
		return
	}

	err := u.userService.LogoutUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successful Logout"})
}

func (u *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	allUsers, err := u.userService.GetAllUsers()
	if err != nil {
		errorhandling.HandleError(w, "Unable to Get Users", http.StatusInternalServerError)
	}
	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Here is the data of all users",
		Data:    allUsers,
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}
