package main

import (
	bootstrap "github.com/Prompiriya084/go-mq/InventoryService/Internal/Bootstrap"
	database "github.com/Prompiriya084/go-mq/InventoryService/Internal/Infrastructure/Database"
	rabbitMQ "github.com/Prompiriya084/go-mq/InventoryService/Internal/Infrastructure/Eventbus"
)

// @title Inventory Service API
// @version 1.0
// @description API Example
// @BasePath /
// @schemes http
func main() {
	// app := fiber.New()
	db := database.InitDb()

	ch, err := rabbitMQ.MQConnectionInit()
	if err != nil {
		panic(err.Error())
	}

	// 👇 topology อยู่ตรงนี้
	if err := rabbitMQ.DeclareTopology(ch); err != nil {
		panic(err.Error())
	}

	bootstrap.StartInventoryAPI(db, ch)
	bootstrap.StartInventoryConsumer(db, ch)
	// go inventoryHandler.CheckStock()
	// go inventoryHandler.ReverseStock()

	// app.Listen(":8081")
}
