package error

import (
	errOrder "order-service/constants/error/order"
)

// jadi intinya dari semua error yang telah kira buat
// kita mapping di sini
func ErrMapping(err error) bool {

	var (
		GeneralErrors = GeneralErrors
		TimeErrors    = errOrder.OrderErrors
	)

	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, TimeErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
