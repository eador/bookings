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

// GetUserByID returns a user by id
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

// UpdateUser updates a user in the database
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticate authentices a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "me@here.com" {
		return 1, "", nil
	}
	return 0, "", errors.New("some error")
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// NewReservations returns a slice of all reservations
func (m *testDBRepo) NewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

// GetReservationByID gets a reservation from the database using the ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation
	return reservation, nil
}

// UpdateReservation updates a user in the database
func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdatedProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRestrictionsForRoomByDate resturns restrictions for a room by a date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	return nil
}

// DeleteBlockByID deletes a room restriction
func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
