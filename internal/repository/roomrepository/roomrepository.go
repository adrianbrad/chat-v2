package roomrepository

import (
	"database/sql"

	"github.com/adrianbrad/chat-v2/internal/repository"
	"github.com/adrianbrad/chat-v2/internal/room"
)

type RoomRepositoryDB struct {
	repository.DBRepository
}

func NewRoomRepositoryDB(db *sql.DB) *RoomRepositoryDB {
	return &RoomRepositoryDB{db}
}

func (r *RoomRepositoryDB) GetAll() (rooms []*room.Room, err error) {
	rows, err := r.Query(`
	SELECT room_id
	FROM rooms
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	var roomID string
	for rows.Next() {
		err = rows.Scan(&roomID)
		if err != nil {
			rooms = []*room.Room{}
			return
		}
		rooms = append(rooms, room.New(roomID))
	}

	err = rows.Err()

	return
}
