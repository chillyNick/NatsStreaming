package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chillyNick/NatsStreaming/internal/config"
	"github.com/chillyNick/NatsStreaming/internal/model"
	stan "github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
)

const jsonOrder = `
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`

func main() {
	if err := config.ReadConfigYML("config.yaml"); err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration")
	}

	cfg := config.GetConfigInstance().NatsStreaming

	sc, err := stan.Connect(
		cfg.ClusterId,
		cfg.PublisherId,
		stan.NatsURL(fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port)),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect nats-streaming")
	}

	order := new(model.Order)
	err = json.Unmarshal([]byte(jsonOrder), order)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshalize order")
	}

	for i := 0; true; i++ {
		order.OrderUID = strconv.Itoa(i)
		msg, err := json.Marshal(order)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to marshalize order")
		}

		err = sc.Publish(cfg.Subject, msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send order ino nuts")
		}

		time.Sleep(time.Minute)
	}
}
