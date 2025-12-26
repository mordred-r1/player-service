package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mordred-r1/player-service/internal/models"
	"github.com/pkg/errors"
)

type PGstorage struct {
	db *pgxpool.Pool
}

func NewPGStorage(connString string) (*PGstorage, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}
	storage := &PGstorage{
		db: db,
	}
	err = storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGstorage) initTables() error {
	sql := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %v (
        %v TEXT PRIMARY KEY,
        %v TEXT
    )`, tableName, IDColumnName, StateColumnName)
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "initition tables")
	}
	return nil
}

// Create inserts a new player record.
func (s *PGstorage) Create(ctx context.Context, player *models.PlayerState) error {
	sql := fmt.Sprintf("INSERT INTO %v (%v, %v) VALUES ($1, $2)", tableName, IDColumnName, StateColumnName)
	_, err := s.db.Exec(ctx, sql, player.ID, player.State)
	if err != nil {
		return errors.Wrap(err, "create player")
	}
	return nil
}

// Get retrieves a player by id.
func (s *PGstorage) Get(ctx context.Context, id string) (*models.PlayerState, error) {
	sql := fmt.Sprintf("SELECT %v, %v FROM %v WHERE %v = $1", IDColumnName, StateColumnName, tableName, IDColumnName)
	row := s.db.QueryRow(ctx, sql, id)
	var state string
	if err := row.Scan(&id, &state); err != nil {
		return nil, errors.Wrap(err, "get player")
	}
	return &models.PlayerState{ID: id, State: state}, nil
}

// Update updates the state for an existing player.
func (s *PGstorage) Update(ctx context.Context, player *models.PlayerState) error {
	sql := fmt.Sprintf("UPDATE %v SET %v = $1 WHERE %v = $2", tableName, StateColumnName, IDColumnName)
	_, err := s.db.Exec(ctx, sql, player.State, player.ID)
	if err != nil {
		return errors.Wrap(err, "update player")
	}
	return nil
}

// Delete removes a player by id.
func (s *PGstorage) Delete(ctx context.Context, id string) error {
	sql := fmt.Sprintf("DELETE FROM %v WHERE %v = $1", tableName, IDColumnName)
	_, err := s.db.Exec(ctx, sql, id)
	if err != nil {
		return errors.Wrap(err, "delete player")
	}
	return nil
}
