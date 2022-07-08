package sqlx_repo

import (
	dbmodel "github.com/chillyNick/NatsStreaming/internal/db_model"
	"github.com/chillyNick/NatsStreaming/internal/repo"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) repo.Repo {
	return repository{db: db}
}

func (r repository) InsertOrder(order []byte) error {
	const query = `
		INSERT INTO "order"(info)
		VALUES ($1)
	`
	_, err := r.db.Exec(query, order)
	return err
}

func (r repository) GetOrders() ([]dbmodel.Order, error) {
	const query = `
		SELECT info ->> 'order_uid' as order_uid, info from "order"
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get orders from db")

		return nil, err
	}

	var orders []dbmodel.Order
	for rows.Next() {
		var order dbmodel.Order
		err = rows.Scan(&order.OrderUid, &order.Info)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshalize data from db")
		}

		orders = append(orders, order)
	}

	return orders, nil
}
