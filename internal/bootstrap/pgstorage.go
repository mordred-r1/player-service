package bootstrap

import (
	"fmt"
	"log"

	"github.com/mordred-r1/player-service/config"
	"github.com/mordred-r1/player-service/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGstorage {

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	storage, err := pgstorage.NewPGStorage(connectionString)
	if err != nil {
		log.Panicf("ошибка инициализации БД, %v", err)
		panic(err)
	}
	return storage
}
