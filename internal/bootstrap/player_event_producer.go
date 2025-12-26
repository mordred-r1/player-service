package bootstrap

import (
	"fmt"

	"github.com/mordred-r1/player-service/config"
	playereventproducer "github.com/mordred-r1/player-service/internal/producer/player_event_producer"
)

func InitPlayerEventProducer(cfg *config.Config) *playereventproducer.PlayerEventProducer {
	kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return playereventproducer.NewPlayerEventProducer(kafkaBrokers, cfg.Kafka.PlayerEventsTopic)
}
