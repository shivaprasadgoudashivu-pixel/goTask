package handler

import (
	"errors"
	"orderapi/database"
	"orderapi/model"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var wg *sync.WaitGroup = new(sync.WaitGroup)

type OrderHandler struct {
	database.IOrderDB
}

type IOrderHandler interface {
	CreateOrders(o *fiber.Ctx) error
	GetOrders(o *fiber.Ctx) error
	ConfirmOrders(o *fiber.Ctx) error
}

func NewOrderHandler(iorderDb database.IOrderDB) IOrderHandler {
	return &OrderHandler{iorderDb}
}

func (ord *OrderHandler) CreateOrders(o *fiber.Ctx) error {
	order := new(model.Order)
	err := o.BodyParser(order)

	if err != nil {
		return errors.New("Error while parsing order payload")
	}
	err = order.ValidateOrder()
	if err != nil {
		return err
	}
	order.TIMESTAMP = time.Now().Unix()
	order.Status = "pending" // Simulating ID generation
	val, err := ord.CreateOrder(order)
	if err != nil {
		return errors.New("Error while creating order")
	}
	return o.JSON(val)
}

func (ord *OrderHandler) GetOrders(o *fiber.Ctx) error {
	id := o.Params("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id")
	}
	order, err := ord.GetOrd(int(Id))
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}

	return o.JSON(order)

}

func (ord *OrderHandler) ConfirmOrders(o *fiber.Ctx) error {

	ch := make(chan string, 10)
	id := o.Params("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id")
	}
	go func() {
		wg.Add(1)
		defer wg.Done()
		ch <- "Confirmed"
		log.Info().Msg("Order confirmation started")
		time.Sleep(2 * time.Second)
		close(ch)
		log.Info().Msg("Order confirmation completed")
	}()
	order, err := ord.ConfirmOrder(int(Id), ch)
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}

	return o.JSON(fiber.Map{
		"message": order,
	})
}
