package main

import (
	"fmt"
	"os"

	"github.com/mordred-r1/player-service/config"
	"github.com/mordred-r1/player-service/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}
	playerStorage := bootstrap.InitPGStorage(cfg)
	playerEventProducer := bootstrap.InitPlayerEventProducer(cfg)
	redisCache := bootstrap.InitRedis(cfg)
	playerService := bootstrap.InitPlayerService(playerStorage, cfg, playerEventProducer, redisCache)

	roomEventProcessor := bootstrap.InitRoomEventProcessor(playerService)
	roomEventConsumer := bootstrap.InitRoomEventConsumer(cfg, roomEventProcessor)

	playerCommandProcessor := bootstrap.InitPlayerCommandProcessor(playerService)
	playerCommandConsumer := bootstrap.InitPlayerCommandConsumer(cfg, playerCommandProcessor)

	bootstrap.AppRun(*playerCommandConsumer, *roomEventConsumer)
}
