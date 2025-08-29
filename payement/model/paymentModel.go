package model

type PAYMENT struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Status string `json:"status`
}
