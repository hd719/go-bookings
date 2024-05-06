// Functions that will be used to mutate the data in the database
package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/hd719/go-bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// If the the insert operation is taking more than 3 seconds cancel it
	// Context: Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	// OLD: Inserting into DB
	// _, err := m.DB.ExecContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())

	// New: Query that returns an id and sets the value to the memory address of var newId
	err := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&newId)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction into the Database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt, r.StartDate, r.EndDate, r.RoomID, r.ReservationID, time.Now(), time.Now(), r.RestrictionID)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDates returns true if availability exists for a specific room and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityByDatesForRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Iterate through all the rows for a given room and see if there are any overlapping dates
	query := `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	var numRows int
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for a given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select rooms.id, rooms.room_name from rooms where rooms.id not in (select rr.room_id from room_restrictions rr where '2021-02-19' < rr.end_date and '2021-02-21' > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows", err)
		return rooms, err
	}

	return rooms, nil
}

// Logic Ex. for SearchAvailabilityByDates
// --existing reservation
// start_date: 2021-02-01
// end_date: 2021-02-04

// -- start date is exactly the same as existing reservation
// select count(id) from room_restrictions where '2021-02-01' < end_date and "2021-02-04" > start_date

// -- start date is before existing reservation and end date is the same
// select count(id) from room_restrictions where '2021-01-31' < end_date and "2021-02-04" > start_date

// -- end date is after the existing reservation and start date is the same
// select count(id) from room_restrictions where '2021-02-01' < end_date and "2021-02-05" > start_date

// -- both start and end dates are outside of all existing reservations, but cover the reservations
// select count(id) from room_restrictions where '2021-31-01' < end_date and "2021-02-05" > start_date

// What does the query do? Lists all available rooms within the dates given
// Get me all of the room ids and room names from the rooms table where the id from the rooms table (rooms.id) is not in this query (select rr.room_id from room_restrictions rr where '2021-02-19' < rr.end_date and '2021-02-21' > rr.start_date -> returns a row of room ids that are booked within the given dates)
// select rooms.id, rooms.room_name from rooms where rooms.id not in (select rr.room_id from room_restrictions rr where '2021-02-19' < rr.end_date and '2021-02-21' > rr.start_date)
