package main

import (
	"demo/database"
	"demo/model"
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

	db, err := database.GetConnection(DSN)
	service := "Order-payement"
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}
	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)

	app := fiber.New()
	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&model.PAYMENT{})
}
