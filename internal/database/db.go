package database

import (
	"fmt"

	"github.com/chillyNick/NatsStreaming/internal/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func NewPostgresql(cfg config.Database) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"pgx",
		fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.Name,
		),
	)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create database connection")

		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msgf("failed to ping postgresql")

		return nil, err
	}

	return db, nil
}
