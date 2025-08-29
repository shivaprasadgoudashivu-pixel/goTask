package main

import (
	"food/database"
	"food/handler"
	mesagging "food/kafka/messaging"
	"food/model"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN           string
	PORT          string
	SEEDS         string
	Topic         string
	ConsumerGroup string
)

func main() {

	service := "Food-Delivery"
	if SEEDS == "" {
		SEEDS = "kafka1, kafka2, kafka3"
	}

	DSN = os.Getenv("DSN")
	if DSN == "" {
		DSN = `host=localhost user=app password=app123 dbname=ordersdb port=5432 sslmode=disable`
		log.Info().Msg(DSN)
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	db, err := database.GetConnection(DSN)

	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}
	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)

	msgUsersCreated := mesagging.NewMessaging("ordersv1", strings.Split(SEEDS, ","))
	go msgUsersCreated.ProduceRecords()

	app := fiber.New()
	FDHandler := handler.NewFDHandler(database.NewFoodDB(db))

	FD_group := app.Group("/api/v1/orders")
	FD_group.Post("/", FDHandler.Create_Order(msgUsersCreated))
	FD_group.Get("/:order_id", FDHandler.Get_Order_By_ID)
	FD_group.Get("/status/:status/limit/:limit", FDHandler.Get_Order_By_Status)

	app.Listen(":" + PORT)
	err = http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		println(err.Error())
		runtime.Goexit()
	}

	// flag.StringVar(&Topic, "topic", "ordersv1", "ordersv1")
	// // flag.StringVar(&ConsumerGroup, "cg", " ", "--cg=demo-consumer-group")
	// flag.Parse()

	// seeds := []string{"localhost:19092", "localhost:29092", "localhost:39092"}
	// // One client can both produce and consume!
	// // Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	// cl, err := kgo.NewClient(
	// 	kgo.SeedBrokers(seeds...),
	// 	kgo.ConsumerGroup(ConsumerGroup),
	// 	kgo.ConsumeTopics(Topic),
	// 	//kgo.RequiredAcks(kgo.AllISRAcks()), // or kgo.RequireOneAck(), kgo.RequireNoAck()
	// 	//kgo.DisableIdempotentWrite()
	// 	//kgo.RetryTimeout()
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// defer cl.Close()

	// ctx := context.Background()
	// time.Sleep(time.Second * 5)
	// for {
	// 	fetches := cl.PollFetches(ctx)
	// 	if errs := fetches.Errors(); len(errs) > 0 {

	// 		panic(fmt.Sprint(errs))
	// 	}

	// 	// We can iterate through a record iterator...
	// 	iter := fetches.RecordIter()
	// 	for !iter.Done() {
	// 		record := iter.Next()
	// 		var event model.OrderEvents
	// 		if err := json.Unmarshal(record.Value, &event); err != nil {
	// 			log.Printf("failed to unmarshal message: %v", err)
	// 			continue
	// 		}
	// 		fmt.Println("here updates the order table after order events ", event)
	// 		FDHandler.UpdateOrder(&event)
	// 		fmt.Println("Partition-->", record.Partition, "Topic-->", record.Topic, string(record.Value), "from an iterator!")
	// 	}

	// }

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&model.ORDER{}, &model.OrderEvents{})
}
