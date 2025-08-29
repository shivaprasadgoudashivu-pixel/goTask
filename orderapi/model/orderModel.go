package model

import "errors"

type Order struct {
	COMMONMODEL
	Order_Id    int `json:"order_id" gorm:"primaryKey;autoIncrement"`
	User_Id     int `json:"user_id" `
	Total_cents int `json:"total_cents"`
}

func (order *Order) ValidateOrder() error {

	if order.User_Id == 0 {
		return errors.New("user ID cannot be empty")
	}
	// if order.Status == "" {
	// 	return errors.New("status cannot be empty")
	// }
	if order.Total_cents < 0 {
		return errors.New("total cents cannot be negative")
	}
	return nil

}
