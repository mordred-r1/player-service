package roomeventprocessor

import (
	"context"
	"strings"

	"github.com/mordred-r1/player-service/internal/models"
)

func (p *RoomEventProcessor) Handle(ctx context.Context, roomEvent *models.RoomEvent) error {
	content := strings.TrimSpace(roomEvent.Content)

	switch {
	case content == "created":
		p.playerService.Create(ctx, &models.PlayerState{ID: roomEvent.Name, State: "created"})
	case content == "deleted":
		p.playerService.Delete(ctx, roomEvent.Name)
	default:
		// unknown command
	}
	return nil
}
