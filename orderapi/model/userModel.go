package model

import "errors"

type User struct {
	ID	 int     `json:"id" gorm:"primaryKey;autoIncrement"`
	NAME   string  `json:"name"`
	EMAIL  string  `json:"email" gorm:"unique"`
	Orders []Order `json:"orders" gorm:"foreignKey:User_Id;references:ID"`
}

func (user *User) ValidateUser() error {

	if user.NAME == "" {
		return errors.New("name cannot be empty")
	}
	if user.EMAIL == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}
