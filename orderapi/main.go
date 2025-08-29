package main

import (
	"orderapi/database"
	"orderapi/handler"
	"orderapi/model"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN  string
	PORT string
)

func main() {

	DSN = os.Getenv("DSN")
	if DSN == "" {
		DSN = `host=localhost user=app password=app123 dbname=usersdb port=5432 sslmode=disable`
		log.Info().Msg(DSN)
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	service := "orderapi"
	app := fiber.New()

	db, err := database.GetConnection(DSN)
	Init(db)

	if err != nil {
		//log.Fatal().Msg("unable to connect to the database..." + err.Error())
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}
	log.Info().Str("service", service).Msg("database connection is established")

	OrderHandler := handler.NewOrderHandler(database.NewOrderDB(db))
	userHandler := handler.NewUserHandler(database.NewUserDB(db))
	user_group := app.Group("/api/v1/users")
	user_group.Post("/", userHandler.CreateUser)
	user_group.Get("/:id", userHandler.GetUserBy)
	orders_group := app.Group("/OrdersAPiCall")
	orders_group.Post("/", OrderHandler.CreateOrders)
	orders_group.Get(":id/", OrderHandler.GetOrders)
	orders_group.Get(":id/confirm", OrderHandler.ConfirmOrders)

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&model.User{})
}
