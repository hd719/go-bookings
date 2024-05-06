package dbrepo

import (
	"errors"
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

	// At the moment we only have 2 rooms with id 1 and 2
	if id > 2 {
		return room, errors.New("room doesnt exist")
	}

	return room, nil
}

func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Authenticates a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil
}
