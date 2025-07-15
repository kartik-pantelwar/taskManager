// create user struct type
package user

import "time"

type UserProfile struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
type User struct {
	Uid       int       `json:"uid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserResponse struct {
	Uid       int       `json:"uid"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserResponse struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
}
