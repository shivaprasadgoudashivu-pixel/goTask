package database

import (
	"orderapi/model"

	"gorm.io/gorm"
)

type IUserDB interface {
	Create(user *model.User) (*model.User, error)
	GetBy(id uint) (*model.User, error)
}
type UserDb struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) IUserDB {
	return &UserDb{db}
}

func (udb *UserDb) Create(user *model.User) (*model.User, error) {
	tx := udb.DB.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (udb *UserDb) GetBy(id uint) (*model.User, error) {
	user := new(model.User)
	tx := udb.DB.First(user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
