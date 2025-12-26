package pgstorage

type PlayerState struct {
	ID    string `db:"id"`
	State string `db:"state"`
}

const (
	tableName = "playerStates"

	IDColumnName    = "id"
	StateColumnName = "state"
)
