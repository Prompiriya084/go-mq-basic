package bootstrap

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
