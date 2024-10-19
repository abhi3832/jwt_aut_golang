package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct{
	ID uuid.UUID			`json:"idx"`
	First_name string		`json:"first_name" validate:"required, min=2, max = 100"`
	Last_name string		`json:"last_name" validate:"required, min=2, max = 100"`
	Passward string			`json:"passward" validate:"required, min=8"`
	Email string			`json:"email" validate:"required"`
	Phone string			`json:"phone" validate:"required"`
	User_type string		`json:"user_type" validate:"required, eq=ADMIN | eq=USER"`
	Token string			`json:"token"`
	Refresh_token string	`json:"refresh_token"`
	Created_at time.Time	`json:"created_at"`
	Updated_at time.Time	`json:"updated_at"`
	User_id string			`json:"user_id"`
}