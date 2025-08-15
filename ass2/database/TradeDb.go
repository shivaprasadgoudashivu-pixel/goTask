package database

import (
	"ass2/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ITradeDB interface {
	AddTrade_Values(trade *models.TradeModel) (*models.TradeModel, error)
	GetPositions(sym string) (pos string, err error)
}

type TradeDB struct {
	DB *gorm.DB
}

type Totals struct {
	TotalPrc float64
	TotalQty float64
}

func NewTradeDB(db *gorm.DB) ITradeDB {
	return &TradeDB{db}
}

func (tradeDB *TradeDB) AddTrade_Values(trade *models.TradeModel) (*models.TradeModel, error) {
	t := tradeDB.DB.Create(trade)
	if t.Error != nil {
		return nil, errors.New("unable to add trade: ")
	}
	return trade, nil
}

func (tradeDB *TradeDB) GetPositions(sym string) (pos string, err error) {
	var total Totals
	result1 := tradeDB.DB.
		Model(&models.TradeModel{}).
		Where("sym = ?", sym).
		Select("SUM(prc),SUM(qty) AS total_qty").
		Scan(&total)

	if result1.Error != nil {
		return "", errors.New("unable to GetPositions: ")
	}
	pos = fmt.Sprintf("Total prc: %f, Total qty: %f of Sym %s\n", total.TotalPrc, total.TotalQty, sym)
	return pos, nil
}
