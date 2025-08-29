package model

import "encoding/json"

type ORDER struct {
	OrderId      int    `json:"order_id" gorm:"primaryKey;autoIncrement"`
	CustomerName string `json:"customer_name"`
	Item         string `json:"item"`
	Address      string `json:"address"`
	COMMONMODEL
	OrderEvents []OrderEvents `json:"orders" gorm:"foreignKey:Order_Id;references:OrderId"`
}

func (o *ORDER) ToBytes() []byte {
	bytes, _ := json.Marshal(o)
	return bytes
}
