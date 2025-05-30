package clients

import "github.com/google/uuid"

type UserResponse struct {
	Code    int      `json:"code"`
	Status  int      `json:"status"`
	Message int      `json:"message"`
	Data    UserData `json:"data"`
}

type UserData struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	PhoneNumber string    `json:"phoneNumber"`
}
