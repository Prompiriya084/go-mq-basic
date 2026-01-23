package bootstrap

import (
	"os"

	adapters_eventbus "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Eventbus"
	adapters_handlers "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Handlers"
	adapters_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Services"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

// type OrderDependencies struct {
// 	Config       config.Config
// 	OrderHandler *http.OrderHandler
// 	MQConsumer   *mq.Consumer
// }
// func NewOrderService() (*OrderDependencies, error) {
// 	cfg := config.Load()

// 	db, err := database.Connect(cfg.Database)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// MQ
// 	mqConn, mqCh := mq.Connect(cfg.RabbitMQ)
// 	mq.DeclareTopology(mqCh)

// 	publisher := mq.NewPublisher(mqCh)
// 	consumer  := mq.NewConsumer(mqCh)

// 	// Application
// 	repo := application.NewOrderRepository(db)
// 	usecase := application.NewOrderUseCase(repo, publisher)

// 	// HTTP
// 	handler := http.NewOrderHandler(usecase)

// 	return &OrderDependencies{
// 		Config:       cfg,
// 		OrderHandler: handler,
// 		MQConsumer:   consumer,
// 	}, nil
// }

func StartInventoryConsumer(db *gorm.DB, ch *amqp091.Channel) error {
	// ch, err := rabbitMQ.MQConnectionInit(os.Getenv("RABBITMQ_URL"))
	// if err != nil {
	// 	panic(err.Error())
	// }

	// // 👇 topology อยู่ตรงนี้
	// if err := rabbitMQ.DeclareTopology(ch); err != nil {
	// 	return err
	// }

	eventPublisher := adapters_eventbus.NewRabbitPublisher(ch, os.Getenv("MQ_Inventory_Core_Exchange"))
	repo := adapters_repositories.NewInventoryRepository(db)
	service := services.NewInventoryEventService(repo, eventPublisher)
	subscriber := adapters_handlers.NewInventorySubscriber(ch, service, 3)

	return subscriber.Start()
}
