package bootstrap

import (
	playerservice "github.com/mordred-r1/player-service/internal/services/playerService"
	playercommandprocessor "github.com/mordred-r1/player-service/internal/services/processors/player_command_processor"
	roomeventprocessor "github.com/mordred-r1/player-service/internal/services/processors/room_event_processor"
)

func InitRoomEventProcessor(playerService *playerservice.PlayerService) *roomeventprocessor.RoomEventProcessor {
	return roomeventprocessor.NewRoomEventProcessor(playerService)
}

func InitPlayerCommandProcessor(playerService *playerservice.PlayerService) *playercommandprocessor.PlayerCommandProcessor {
	return playercommandprocessor.NewPlayerCommandProcessor(playerService)
}
