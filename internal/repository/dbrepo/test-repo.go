package dbrepo

import (
	"errors"
	"time"

	"github.com/eador/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 1, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestricition(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start time.Time, end time.Time, roomID int) (bool, error) {
	if roomID == 2 {
		return false, nil
	}
	if roomID > 2 {
		return false, errors.New("some error")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for a given start and end date
func (m *testDBRepo) SearchAvailablitiyForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	if start.Before(time.Now()) {
		if end.Before(time.Now()) {
			return rooms, errors.New("some error")
		} else {
			return rooms, nil
		}
	}
	var room models.Room

	rooms = append(rooms, room)
	return rooms, nil
}

// GetRoomById gets a room by id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}
