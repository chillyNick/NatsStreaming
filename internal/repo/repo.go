package repo

import dbmodel "github.com/chillyNick/NatsStreaming/internal/db_model"

type Repo interface {
	InsertOrder(order []byte) error
	GetOrders() ([]dbmodel.Order, error)
}
