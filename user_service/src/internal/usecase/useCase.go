package userservice

import (
	"fmt"
	"time"
	"user_service/src/internal/adaptors/persistance"
	"user_service/src/internal/core/session"
	"user_service/src/internal/core/user"
	"user_service/src/pkg/utilities"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    persistance.UserRepo
	sessionRepo persistance.SessionRepo
}

func NewUserService(userRepo persistance.UserRepo, sessionRepo persistance.SessionRepo) UserService {
	return UserService{userRepo: userRepo, sessionRepo: sessionRepo}
}

// registration function definition
func (u *UserService) RegisterUser(user user.UserRegister) (user.UserResponse, error) {
	//^ checking if user is already registered
	newUser, err := u.userRepo.CreateUser(user)
	return newUser, err
}

type LoginResponse struct {
	FounUser    user.User
	TokenString string
	TokenExpire time.Time
	Session     session.Session
}

func (u *UserService) LoginUser(requestUser user.UserLogin) (LoginResponse, error) {
	loginResponse := LoginResponse{}

	foundUser, err := u.userRepo.GetUser(requestUser.Username)
	if err != nil {
		return loginResponse, fmt.Errorf("invalid username")
	}

	loginResponse.FounUser = foundUser
	if err := matchPassword(requestUser, foundUser.Password); err != nil {
		return loginResponse, fmt.Errorf("invalid password")
	}
	tokenString, tokenExpire, err := utilities.GenerateJWT(foundUser.Uid)
	loginResponse.TokenString = tokenString
	loginResponse.TokenExpire = tokenExpire

	if err != nil {
		return loginResponse, fmt.Errorf("failed to generate jwt")
	}

	session, err := utilities.GenerateSession(foundUser.Uid)
	loginResponse.Session = session
	if err != nil {
		return loginResponse, fmt.Errorf("failed to generate session")
	}

	err = u.sessionRepo.CreateSession(session)
	if err != nil {
		return loginResponse, fmt.Errorf("failed to create session")
	}

	return loginResponse, nil
}

func (u *UserService) GetJwtFromSession(sess string) (string, time.Time, error) {
	var tokenString string
	var tokenExpire time.Time
	session, err := u.sessionRepo.GetSession(sess)
	if err != nil {
		return tokenString, tokenExpire, err
	}

	err = matchSessionToken(sess, session.TokenHash)
	if err != nil {
		return tokenString, tokenExpire, err
	}

	tokenString, tokenExpire, err = utilities.GenerateJWT(session.Uid)
	if err != nil {
		return tokenString, tokenExpire, err
	}

	return tokenString, tokenExpire, nil
}

func (u *UserService) GetUserByID(id int) (user.UserProfile, error) {
	newUser, err := u.userRepo.GetUserByID(id)
	return newUser, err
}

func (u *UserService) LogoutUser(id int) error {
	err := u.sessionRepo.DeleteSession(id)
	return err
}

func matchPassword(user user.UserLogin, password string) error {
	// !error here
	err := utilities.CheckPassword(password, user.Password)
	if err != nil {
		return fmt.Errorf("unable to match password: %v", err)
	}

	return nil
}

func matchSessionToken(id string, tokenHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(tokenHash), []byte(id))
	if err != nil {
		fmt.Println(err, "Unable to Match Password")
	}
	return nil
}

func (u *UserService) GetAllUsers() ([]user.GetUserResponse, error) {
	allUsers, err := u.userRepo.GetUsers()
	if err != nil {
		return []user.GetUserResponse{}, err
	}
	return allUsers, nil
}
