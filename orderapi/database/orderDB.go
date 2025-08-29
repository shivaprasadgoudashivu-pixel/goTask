package database

import (
	"errors"
	"fmt"
	"orderapi/model"

	"gorm.io/gorm"
)

type IOrderDB interface {
	CreateOrder(order *model.Order) (*model.Order, error)
	GetOrd(orderId int) (*model.Order, error)
	ConfirmOrder(orderId int, ch chan string) (string, error)
}

type OrderDB struct {
	DB *gorm.DB
}

func NewOrderDB(db *gorm.DB) IOrderDB {
	return &OrderDB{db}
}

func (ODB *OrderDB) CreateOrder(order *model.Order) (*model.Order, error) {
	tx := ODB.DB.Create(order)
	if tx.Error != nil {
		return nil, errors.New("Something Went Wrong Try Again")
	}
	return order, nil
}

func (ODB *OrderDB) GetOrd(orderId int) (*model.Order, error) {
	order := new(model.Order)
	tx := ODB.DB.First(order, orderId)
	if tx.Error != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (ODB *OrderDB) ConfirmOrder(orderId int, ch chan string) (string, error) {

	select {
	case msg := <-ch:

		tx := ODB.DB.Model(&model.Order{}).Where("id = ?", orderId).Update("status", msg)
		if tx.Error != nil {
			return "", errors.New("failed to confirm order")
		}
		fmt.Println(tx)
		return "Request processes succefully", nil

	}

}
