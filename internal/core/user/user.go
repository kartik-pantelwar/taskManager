// create user struct type
package user

import "time"

type User struct {
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Uid           int       `json:"uid"`
	Work_location string    `json:"work_location"`
	Is_available  bool      `json:"is_available"`
	Created_at    time.Time `json:"created_at"`
	Balance       float64   `json:"balance"`
	Rating        float64   `json:"rating"`
}
