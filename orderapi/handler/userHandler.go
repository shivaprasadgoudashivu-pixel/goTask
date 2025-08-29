package handler

import (
	"errors"
	"orderapi/database"
	"orderapi/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	database.IUserDB // prmoted field
}

type IUserHandler interface {
	CreateUser(c *fiber.Ctx) error
	GetUserBy(c *fiber.Ctx) error
}

func NewUserHandler(iuserdb database.IUserDB) IUserHandler {
	return &UserHandler{iuserdb}
}

func (uh *UserHandler) CreateUser(c *fiber.Ctx) error {

	user := new(model.User)
	err := c.BodyParser(user)
	if err != nil {
		return err
	}

	err = user.ValidateUser()
	if err != nil {
		return err
	}
	
	user, err = uh.Create(user)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (uh *UserHandler) GetUserBy(c *fiber.Ctx) error {
	id := c.Params("id") // Retrieves the value of ":id"

	_id, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id")
	}

	user, err := uh.GetBy(uint(_id))
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}

	return c.JSON(user)
}
