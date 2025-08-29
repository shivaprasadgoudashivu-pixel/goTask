package database


import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	RETRY_COUNT    = 10
	RETRY_DURATION = 5 // int seconds
)

func GetConnection(dsn string) (*gorm.DB, error) {
	count := 1
retry:
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Err(err).Str("layer", "db").Msg("failed to connect to the database")

		if count > RETRY_COUNT {
			return nil, err
		}
		count++
		log.Warn().Str("layer", "db").Msg("Trying to connect to the database ...for the number of times.." + fmt.Sprint(count))
		time.Sleep(time.Second * RETRY_DURATION)
		goto retry
		//return nil, err
	} else {
		log.Info().Str("layer", "db").Msg("successfully connected to the database")
		return db, nil
	}
}