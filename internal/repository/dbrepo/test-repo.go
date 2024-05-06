package dbrepo

import (
	"time"

	"github.com/hd719/go-bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction into the Database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDates returns true if availability exists for a specific room and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesForRoomId(start, end time.Time, roomId int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for a given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets room by ID
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	return room, nil
}
