package handler

import (
	"ass2/database"
	"ass2/models"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type TradeHandler struct {
	database.ITradeDB
}

type TradeHandlerInterface interface {
	AddTrade(c *fiber.Ctx) error
	GetPos(c *fiber.Ctx) error
}

func NewTradeHandler(itreadeDB database.ITradeDB) TradeHandlerInterface {
	return &TradeHandler{itreadeDB}
}

func (t *TradeHandler) AddTrade(c *fiber.Ctx) error {
	trade := new(models.TradeModel)
	err := c.BodyParser(trade)
	if err != nil {
		return err
	}

	err = trade.Validate()
	if err != nil {
		return err
	}

	trade, err = t.AddTrade_Values(trade)
	if err != nil {
		return err
	}
	return c.JSON(trade)

}
func (t *TradeHandler) UpdateTrade(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (t *TradeHandler) GetPos(c *fiber.Ctx) error {
	sym := c.Params("sym")
	if sym == "" {
		return errors.New("symbol is required")
	}
	pos, err := t.GetPositions(sym)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"positions": pos,
	})
}
