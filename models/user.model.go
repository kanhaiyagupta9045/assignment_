package models

import "time"

type User struct {
	User_Id    uint32    `json:"user_id"`
	First_Name string    `json:"first_name" validate:"required"`
	Last_Name  string    `json:"last_name" validate:"required"`
	Email      string    `json:"email" validate:"email,required"`
	Password   string    `json:"password" validate:"required"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Deleted_at time.Time `json:"deleted_at"`
}

type LoginInfo struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}
