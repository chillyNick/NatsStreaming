package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chillyNick/NatsStreaming/internal/cache"
	"github.com/chillyNick/NatsStreaming/internal/config"
	"github.com/chillyNick/NatsStreaming/internal/consumer"
	"github.com/chillyNick/NatsStreaming/internal/database"
	"github.com/chillyNick/NatsStreaming/internal/http_server"
	"github.com/chillyNick/NatsStreaming/internal/repo"
	"github.com/chillyNick/NatsStreaming/internal/repo/sqlx_repo"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := config.ReadConfigYML("config.yaml"); err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration")
	}

	cfg := config.GetConfigInstance()

	db, err := database.NewPostgresql(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	r := sqlx_repo.New(db)

	c := cache.NewMemoryCache()

	if err = initCache(c, r); err != nil {
		log.Fatal().Err(err).Msg("Failed to init cache")
	}

	cons := consumer.NewConsumer(r, c)
	if err = cons.Observe(cfg.NatsStreaming); err != nil {
		log.Fatal().Err(err).Msg("Failed to observe messages from broker messager")
	}

	server := http_server.New(cfg.Rest, c)
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed running http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	v := <-quit
	log.Info().Msgf("signal.Notify: %v", v)

	if err := server.Close(); err != nil {
		log.Error().Err(err).Msg("httpServer.Close")
	} else {
		log.Info().Msg("httpServer close correctly")
	}

	if err := cons.Stop(); err != nil {
		log.Error().Err(err).Msg("cons.Close")
	} else {
		log.Info().Msg("consumer close correctly")
	}

}

func initCache(c cache.Cache, r repo.Repo) error {
	orders, err := r.GetOrders()
	if err != nil {
		return err
	}

	for _, o := range orders {
		c.Set(o.OrderUid, o.Info)
	}

	return nil
}
