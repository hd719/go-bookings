// Functions that will be used to mutate the data in the database
package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hd719/go-bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
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

	query := `select rooms.id, rooms.room_name from rooms where rooms.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

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

// GetRoomByID gets room by ID
func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
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

// Returns a user by Id
func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// Updates user
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		fmt.Println("Failed to update User")
		return err
	}

	fmt.Println("User updated!")
	return nil
}

// Authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// This will be the id of the user if the credentials are correct
	var id int

	// Will hold the password we get from the DB
	var hashedPassword string

	// check to if the email exists in the Db
	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	// Check to see if the passwords match
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// Returns a slice of all ressys
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close() // prevents memory leak

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName, &i.Email, &i.Phone, &i.StartDate, &i.EndDate, &i.RoomID, &i.CreatedAt, &i.UpdatedAt, &i.Processed, &i.Room.ID, &i.Room.RoomName)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// Returns new reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) where processed = 0 order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close() // prevents memory leak

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(&i.ID, &i.FirstName, &i.LastName, &i.Email, &i.Phone, &i.StartDate, &i.EndDate, &i.RoomID, &i.CreatedAt, &i.UpdatedAt, &i.Room.ID, &i.Room.RoomName)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// Returns 1 reservation by id
func (m *postgresDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name from reservations r left join rooms rm on (r.room_id = rm.id) where r.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID, &res.FirstName, &res.LastName, &res.Email, &res.Phone, &res.StartDate, &res.EndDate, &res.RoomID, &res.CreatedAt, &res.UpdatedAt, &res.Processed, &res.Room.ID, &res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

// Update Reservation
func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5 where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)
	if err != nil {
		fmt.Println("Failed to update Ressy")
		return err
	}

	fmt.Println("Ressy updated!")
	return nil
}

// Deletes ressy
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`

	fmt.Printf("Delete Reservation")

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

// Updates processed for a reservation by id
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `udpate from reservations set processed = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		return err
	}

	return nil
}
