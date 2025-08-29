package database

import (
	"errors"
	"food/model"

	"gorm.io/gorm"
)

type FoodDB struct {
	DB *gorm.DB
}

type IFoodDB interface {
	CreateOrder(order *model.ORDER) (*model.ORDER, error)
	GetOrderByID(id string) (*model.ORDER, error)
	GetOrderByStatus(status string, limit int) ([]model.ORDER, error)
	UpdateOrderEvent(OrderEve *model.OrderEvents) error
	UpdateOrder(order *model.OrderEvents) error
}

func NewFoodDB(db *gorm.DB) IFoodDB {
	return &FoodDB{db}
}

func (f *FoodDB) CreateOrder(order *model.ORDER) (*model.ORDER, error) {
	result := f.DB.Create(order)
	if result.Error != nil {
		return nil, errors.New("unable to create order: ")
	}
	return order, nil
}

func (f *FoodDB) GetOrderByID(id string) (*model.ORDER, error) {
	order := new(model.ORDER)
	result := f.DB.Preload("OrderEvents").First(order, id)
	if result.Error != nil {
		return nil, errors.New("unable to find order by ID: ")
	}
	return order, nil
}

func (f *FoodDB) GetOrderByStatus(status string, limit int) ([]model.ORDER, error) {
	var orders []model.ORDER
	result := f.DB.Preload("OrderEvents").Where("status = ?", status).Order("created_at DESC").
		Limit(limit).Find(&orders)
	if result.Error != nil {
		return nil, errors.New("unable to find orders by status: ")
	}
	return orders, nil
}

func (f *FoodDB) UpdateOrderEvent(OrderEve *model.OrderEvents) error {
	result := f.DB.Create(OrderEve)
	if result.Error != nil {
		return errors.New("unable to update order event: ")
	}
	return nil
}

func (f *FoodDB) UpdateOrder(orderEve *model.OrderEvents) error {
	order := new(model.ORDER)
	result := f.DB.Model(order).Where("order_id = ?", orderEve.Order_Id).Update("status", orderEve.Status)
	if result.Error != nil {
		return errors.New("unable to update order event: ")
	}
	return nil
}
