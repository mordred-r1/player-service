package models

type PlayerState struct {
	ID    string
	State string
}

type RoomEvent struct {
	Name    string
	Content string
}

type PlayerCommand struct {
	Name    string
	Content string
}

type PlayerEvent struct {
	ID    string
	State string
}
