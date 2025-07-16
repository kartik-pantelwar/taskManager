package userservice

import (
	"errors"
	"fmt"
	"log"
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
	newUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		log.Printf("Error: %v", err)
		return newUser, errors.New("Something Went Wrong!")
	}
	return newUser, nil
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
		log.Printf("Error: %v", err)
		return loginResponse, errors.New("Invalid Credentials")
	}

	loginResponse.FounUser = foundUser
	if err := matchPassword(requestUser, foundUser.Password); err != nil {
		log.Printf("Error: %v", err)
		return loginResponse, errors.New("Invalid Credentials")
	}
	tokenString, tokenExpire, err := utilities.GenerateJWT(foundUser.Uid)
	loginResponse.TokenString = tokenString
	loginResponse.TokenExpire = tokenExpire

	if err != nil {
		log.Printf("Error: %v", err)
		return loginResponse, errors.New("Failed to Generate Token")
	}

	session, err := utilities.GenerateSession(foundUser.Uid)
	loginResponse.Session = session
	if err != nil {
		log.Printf("Error: %v", err)
		return loginResponse, errors.New("Failed to Generate Session")
	}

	err = u.sessionRepo.CreateSession(session)
	if err != nil {
		log.Printf("Error: %v", err)
		return loginResponse, errors.New("Failed to Create Session")
	}

	return loginResponse, nil
}

func (u *UserService) GetJwtFromSession(sess string) (string, time.Time, error) {
	var tokenString string
	var tokenExpire time.Time
	session, err := u.sessionRepo.GetSession(sess)
	if err != nil {
		log.Printf("Error: %v", err)
		return tokenString, tokenExpire, errors.New("Invalid Session")
	}

	err = matchSessionToken(sess, session.TokenHash)
	if err != nil {
		log.Printf("Error: %v", err)
		return tokenString, tokenExpire, errors.New("Session Token Mismatch")
	}

	tokenString, tokenExpire, err = utilities.GenerateJWT(session.Uid)
	if err != nil {
		log.Printf("Error: %v", err)
		return tokenString, tokenExpire, errors.New("Failed to Generate Token")
	}

	return tokenString, tokenExpire, nil
}

func (u *UserService) GetUserByID(id int) (user.UserProfile, error) {
	newUser, err := u.userRepo.GetUserByID(id)
	if err != nil {
		log.Println("error", err)
		return user.UserProfile{}, errors.New("User Not Found")
	}
	return newUser, nil
}

func (u *UserService) LogoutUser(id int) error {
	err := u.sessionRepo.DeleteSession(id)
	if err != nil {
		log.Printf("Error: %v", err)
		return errors.New("Failed to Logout User")
	}
	return nil
}

func matchPassword(user user.UserLogin, password string) error {
	// !error here
	err := utilities.CheckPassword(password, user.Password)
	if err != nil {
		log.Printf("Error: %v", err)
		return fmt.Errorf("unable to match password: %v", err)
	}

	return nil
}

func matchSessionToken(id string, tokenHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(tokenHash), []byte(id))
	if err != nil {
		log.Printf("Error: %v", err)
		fmt.Println(err, "Unable to Match Password")
	}
	return nil
}

func (u *UserService) GetAllUsers() ([]user.GetUserResponse, error) {
	allUsers, err := u.userRepo.GetUsers()
	if err != nil {
		log.Printf("Error: %v", err)
		return []user.GetUserResponse{}, errors.New("Unable to Fetch Users")
	}
	return allUsers, nil
}
