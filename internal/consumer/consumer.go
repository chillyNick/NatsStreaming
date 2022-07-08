package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/chillyNick/NatsStreaming/internal/cache"
	"github.com/chillyNick/NatsStreaming/internal/config"
	"github.com/chillyNick/NatsStreaming/internal/repo"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
)

type consumer struct {
	r    repo.Repo
	c    cache.Cache
	conn stan.Conn
	sub  stan.Subscription
}

func NewConsumer(r repo.Repo, c cache.Cache) *consumer {
	return &consumer{r: r, c: c}
}

func (cons *consumer) Observe(cfg config.NatsStreaming) error {
	if cons.conn != nil {
		return fmt.Errorf("firstly you have to stop observe messages")
	}

	sc, err := stan.Connect(
		cfg.ClusterId,
		cfg.ConsumerId,
		stan.NatsURL(fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port)),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect nats-streaming")

		return err
	}

	cons.conn = sc

	// Subscribe with durable name
	cons.sub, err = cons.conn.Subscribe(cfg.Subject, func(m *stan.Msg) {
		log.Debug().Msg("Received a message")

		var msgMap map[string]interface{}
		err := json.Unmarshal(m.Data, &msgMap)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to unmarshal: %s", m.Data)
			return
		}

		orderUid, ok := msgMap["order_uid"]
		if !ok {
			log.Warn().Err(err).Msgf("Received object doesn't contain order_uid field: %s", m.Data)
			return
		}

		orderUidStr := orderUid.(string)

		cons.c.Set(orderUidStr, m.Data)

		if err = cons.r.InsertOrder(m.Data); err != nil {
			log.Warn().Err(err).Msgf("Failed to insert order into db: %s", m.Data)

			return
		}
	}, stan.DurableName(cfg.DurableName))

	if err != nil {
		log.Error().Err(err).Msg("Failed to subscribe nats-streaming")

		return err
	}

	return nil
}

func (cons *consumer) Stop() error {
	err := cons.sub.Unsubscribe()
	if err != nil {
		return err
	}

	return cons.conn.Close()
}
