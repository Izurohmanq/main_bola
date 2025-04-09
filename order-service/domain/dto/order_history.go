package dto

import "order-service/constants"

// kenapa tidak menggunakan `json:"dkkngenttot"`
// karena kita tidak melakukan request melalui Front-End/User/Postman
type OrderHistoryRequest struct {
	OrderID uint
	Status  constants.OrderStatusString
}
