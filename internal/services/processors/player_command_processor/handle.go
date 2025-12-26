package playercommandprocessor

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mordred-r1/player-service/internal/models"
)

func (p *PlayerCommandProcessor) Handle(ctx context.Context, playerCommand *models.PlayerCommand) error {
	content := strings.TrimSpace(playerCommand.Content)

	switch {
	case content == "Play":
		p.playerService.Update(ctx, &models.PlayerState{ID: playerCommand.Name, State: "Play"})
	case content == "Pause":
		p.playerService.Update(ctx, &models.PlayerState{ID: playerCommand.Name, State: "Pause"})
	case content == "Stop":
		p.playerService.Update(ctx, &models.PlayerState{ID: playerCommand.Name, State: "Stop"})
	case strings.HasPrefix(content, "Seek"):
		// content is like "Seek 120" (seconds) or "Seek:120"
		// extract numeric part robustly
		rest := strings.TrimSpace(strings.TrimPrefix(content, "Seek"))
		rest = strings.TrimPrefix(rest, ":")
		if rest == "" {
			return fmt.Errorf("seek command missing position")
		}
		secs, err := strconv.Atoi(strings.Fields(rest)[0])
		if err != nil {
			return fmt.Errorf("invalid seek position: %w", err)
		}
		p.playerService.Update(ctx, &models.PlayerState{ID: playerCommand.Name, State: fmt.Sprintf("Seek %d", secs)})
	default:
		// unknown command
	}
	return nil
}
