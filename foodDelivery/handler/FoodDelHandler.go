//go:generate mockgen -source=FoodDelHandler.go -destination=../internal/mocks/mock_user.go -package=mocks

package handler

import (
	"errors"
	"fmt"
	"food/database"
	mesagging "food/kafka/messaging"
	"food/model"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var status = []string{"Placed", "Accepted", "In-Progress", "Delivered"}

var wg *sync.WaitGroup

type FDHandler struct {
	database.IFoodDB
}

type FDHandlerInterface interface {
	Create_Order(msg *mesagging.Messaging) func(c *fiber.Ctx) error
	Get_Order_By_ID(c *fiber.Ctx) error
	Get_Order_By_Status(c *fiber.Ctx) error
	update_Order_event(o *model.ORDER, msg *mesagging.Messaging) error
	UpdateOrder(o *model.OrderEvents) error
}

func NewFDHandler(ifoodDB database.IFoodDB) FDHandlerInterface {
	return &FDHandler{ifoodDB}
}

func (f *FDHandler) Create_Order(msg *mesagging.Messaging) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		order := new(model.ORDER)
		err := c.BodyParser(order)
		if err != nil {
			return err
		}

		order.Status = "Placed"
		order.CreatedAt = time.Now().Unix()

		order, err = f.CreateOrder(order)
		if err != nil {
			return err
		}
		orderEven := new(model.OrderEvents)
		orderEven.Order_Id = order.OrderId
		f.update_Order_event(order, msg)
		msg.ChMessaging <- order.ToBytes()
		return c.JSON(order)
	}

}

func (f *FDHandler) Get_Order_By_ID(c *fiber.Ctx) error {
	id := c.Params("order_id")
	if id == "" {
		return errors.New("order ID is required")
	}
	order, err := f.GetOrderByID(id)
	if err != nil {
		return errors.New("failed to retrieve order")
	}
	return c.JSON(order)

}

func (f *FDHandler) Get_Order_By_Status(c *fiber.Ctx) error {
	status := c.Params("status")
	limit := c.Params("limit")
	lim, err := strconv.Atoi(limit)

	if status == "" {
		return errors.New("order status is required")
	}
	orders, err := f.GetOrderByStatus(status, lim)
	if err != nil {
		return err
	}
	return c.JSON(orders)

}

func (f *FDHandler) update_Order_event(o *model.ORDER, msg *mesagging.Messaging) error {

	wg = new(sync.WaitGroup)

	if o == nil {
		return errors.New("order cannot be nil")
	}
	go func() {
		for _, val := range status {
			oe := new(model.OrderEvents)
			fmt.Print("Updating order status to: ", val, "\n")

			oe.Order_Id = o.OrderId
			oe.Status = val
			oe.UpDatesAt = time.Now().Unix()

			wg.Add(1)

			time.Sleep(time.Second * 3)
			f.UpdateOrderEvent(oe)
			msg.ChMessaging <- oe.ToBytes()
			wg.Done()
		}

	}()
	return nil
}

func (f *FDHandler) UpdateOrder(oe *model.OrderEvents) error {
	if oe == nil {
		return errors.New("Error Ocurred")
	}
	err := f.UpdateOrder(oe)

	if err != nil {
		return errors.New("Error while updating the table order")
	}
	return nil

}
