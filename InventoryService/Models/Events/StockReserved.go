package events

import "time"

type StockReserved struct {
	OrderID   string
	ProductID string
	Qty       uint
	Timestamp time.Time
}
