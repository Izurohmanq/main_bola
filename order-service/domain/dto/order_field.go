package dto

// kenapa tidak menggunakan `json:"dkkngenttot"`
// karena kita tidak melakukan request melalui Front-End/User/Postman
type OrderFieldRequest struct {
	OrderID         uint
	FieldScheduleID string
}
