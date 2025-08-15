package main

import (
	"ass2/database"
	"ass2/handler"
	"ass2/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN   string
	PORT  string
	debug bool
)

func main() {

	service := "trade-service"
	app := fiber.New()
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

	if err != nil {
		//log.Fatal().Msg("unable to connect to the database..." + err.Error())
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}
	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)

	tradeHandler := handler.NewTradeHandler(database.NewTradeDB(db))
	trade_group := app.Group("/trade")
	trade_group.Post("/", tradeHandler.AddTrade)
	trade_group.Post("/Getposition/:sym", tradeHandler.GetPos)

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.TradeModel{})
}
