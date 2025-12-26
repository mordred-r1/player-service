package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	playercommandconsumer "github.com/mordred-r1/player-service/internal/consumer/player_command_consumer"
	roomeventconsumer "github.com/mordred-r1/player-service/internal/consumer/room_event_consumer"
)

// AppRun starts consumers concurrently and blocks until an interrupt signal is received.
func AppRun(playerCommandConsumer playercommandconsumer.PlayerCommandConsumer, roomEventConsumer roomeventconsumer.RoomEventConsumer) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		playerCommandConsumer.Consume(ctx)
		fmt.Println("PlayerCommandConsumer stopped")
	}()

	go func() {
		defer wg.Done()
		roomEventConsumer.Consume(ctx)
		fmt.Println("RoomEventConsumer stopped")
	}()

	fmt.Println("Consumers started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("Shutdown signal received, cancelling consumers...")
	cancel()
	wg.Wait()
	fmt.Println("All consumers stopped")
}
