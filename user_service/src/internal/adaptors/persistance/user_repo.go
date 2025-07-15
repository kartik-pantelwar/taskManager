package persistance

import (
	user "user_service/src/internal/core/user"
	"user_service/src/pkg/utilities"

	// "TaskManager/pkg/utilities"
	"fmt"
	// "github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(d *Database) UserRepo {
	return UserRepo{db: d}
}

// *done
func (u *UserRepo) CreateUser(newUser user.UserRegister) (user.UserResponse, error) {
	var uid int
	var createdUser user.UserResponse
	// query := "insert into users(username, email, password, work_location, balance) values($1, $2, $3, $4, $5) returning uid"
	hashPass, err := utilities.HashPassword(newUser.Password)
	if err != nil {
		fmt.Println(err, "unable to hash password")
	}

	query := "insert into users(username, email, password) values($1, $2, $3) returning uid, created_at"
	err = u.db.db.QueryRow(query, newUser.Username, newUser.Email, hashPass).Scan(&uid,
		&createdUser.CreatedAt)

	if err != nil {
		return user.UserResponse{}, err
	}
	createdUser.Uid = uid
	return createdUser, nil
}

func (u *UserRepo) GetUser(username string) (user.User, error) {
	var newUser user.User
	query := "select uid, username, email, created_at, password from users where username = $1"
	err := u.db.db.QueryRow(query, username).Scan(&newUser.Uid, &newUser.Username, &newUser.Email, &newUser.CreatedAt, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}

func (u *UserRepo) GetUserByID(id int) (user.UserProfile, error) {
	var newUser user.UserProfile
	query := "select uid, username, email,created_at from users where uid = $1"
	err := u.db.db.QueryRow(query, id).Scan(&newUser.Uid, &newUser.Username, &newUser.Email, &newUser.CreatedAt)
	if err != nil {
		return user.UserProfile{}, err
	}
	return newUser, nil
}

func (u *UserRepo) GetUsers() ([]user.GetUserResponse, error) {
	var allUsers []user.GetUserResponse
	query := `select uid, username from users`
	rows, err := u.db.db.Query(query)
	if err != nil {
		return allUsers, err
	}
	defer rows.Close()
	for rows.Next() {
		var currentUser user.GetUserResponse
		err = rows.Scan(&currentUser.Uid, &currentUser.Username)
		if err != nil {
			return allUsers, err
		}
		allUsers = append(allUsers, currentUser)
	}
	return allUsers, nil
}
