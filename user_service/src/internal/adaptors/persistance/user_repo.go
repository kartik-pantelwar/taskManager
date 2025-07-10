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

func (u *UserRepo) CreateUser(newUser user.User) (user.User, error) {
	var uid int
	// query := "insert into users(username, email, password, work_location, balance) values($1, $2, $3, $4, $5) returning uid"
	hashPass, err := utilities.HashPassword(newUser.Password)
	if err != nil {
		fmt.Println(err, "unable to hash password")
	}
	if newUser.Balance == 0 {
		// Don't include balance - let database use default (5.0)
		query := "insert into users(username, email, password, work_location) values($1, $2, $3, $4) returning uid"
		err = u.db.db.QueryRow(query, newUser.Username, newUser.Email, hashPass, newUser.Work_location).Scan(&uid)
	} else {
		// User explicitly provided balance
		query := "insert into users(username, email, password, work_location, balance) values($1, $2, $3, $4, $5) returning uid"
		err = u.db.db.QueryRow(query, newUser.Username, newUser.Email, hashPass, newUser.Work_location, newUser.Balance).Scan(&uid)
	}
	// err = u.db.db.QueryRow(query, newUser.Username, newUser.Email, hashPass, newUser.Work_location, newUser.Balance).Scan(&uid) //scan will check the numbers of rows executed, and will assign that number to uid
	if err != nil {
		return user.User{}, err
	}
	newUser.Uid = uid
	return newUser, nil
}

func (u *UserRepo) GetUser(username string) (user.User, error) {
	var newUser user.User
	query := "select uid, username, email, balance, password from users where username = $1"
	err := u.db.db.QueryRow(query, username).Scan(&newUser.Uid, &newUser.Username, &newUser.Email, &newUser.Balance, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}

func (u *UserRepo) GetUserByID(id int) (user.User, error) {
	var newUser user.User
	query := "select uid, username, email,work_location,balance,created_at from users where uid = $1"
	// query := "select uid, username, email, password from users where uid = $1"
	//! password not entered in above query
	err := u.db.db.QueryRow(query, id).Scan(&newUser.Uid, &newUser.Username, &newUser.Email, &newUser.Work_location, &newUser.Balance, &newUser.Created_at)
	// err := u.db.db.QueryRow(query, id).Scan(&newUser.Uid, &newUser.Username, &newUser.Email, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}
