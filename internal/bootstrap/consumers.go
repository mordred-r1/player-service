package bootstrap

import (
	"fmt"

	"github.com/mordred-r1/player-service/config"
	playercommandconsumer "github.com/mordred-r1/player-service/internal/consumer/player_command_consumer"
	roomeventconsumer "github.com/mordred-r1/player-service/internal/consumer/room_event_consumer"
	playercommandprocessor "github.com/mordred-r1/player-service/internal/services/processors/player_command_processor"
	roomeventprocessor "github.com/mordred-r1/player-service/internal/services/processors/room_event_processor"
)

func InitRoomEventConsumer(cfg *config.Config, roomEventProcessor *roomeventprocessor.RoomEventProcessor) *roomeventconsumer.RoomEventConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return roomeventconsumer.NewRoomEventConsumer(roomEventProcessor, kafkaBrokers, cfg.Kafka.RoomEventsTopic)
}

func InitPlayerCommandConsumer(cfg *config.Config, playerCommandProcessor *playercommandprocessor.PlayerCommandProcessor) *playercommandconsumer.PlayerCommandConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return playercommandconsumer.NewPlayerCommandConsumer(playerCommandProcessor, kafkaBrokers, cfg.Kafka.PlayerCommandTopic)
}
