package repository

import (
	"time"

	"github.com/eador/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestricition(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start time.Time, end time.Time, roomID int) (bool, error)
	SearchAvailablitiyForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
}
