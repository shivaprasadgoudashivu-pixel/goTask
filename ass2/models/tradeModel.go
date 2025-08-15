package models

import (
	"ass2/utils"
	"encoding/json"
)

type TradeModel struct {
	SYM string `json:"sym" `
	Act string `json:"act"`
	Qty int    `json:"qty"`
	Prc int    `json:"prc"`
}

func (t *TradeModel) Validate() error {
	if t.SYM == "" {
		return utils.ErrInvalidSym
	}
	if t.Qty == 0 {
		return utils.ErrInvalidQty
	}
	if t.Act == "" {
		return utils.ErrInvalidAct
	}
	if t.Prc == 0 {
		return utils.ErrInvalidPrc
	}
	return nil
}

func (t *TradeModel) ToBytes() []byte {
	val, _ := json.Marshal(t)
	return val
}
