// Functions that will be used to mutate the data in the database
package dbrepo

import (
	"context"
	"time"

	"github.com/hd719/go-bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	// If the the insert operation is taking more than 3 seconds cancel it
	// Context: Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations (first_name, last_name, email, phone start_date,
		end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := m.DB.ExecContext(ctx, stmt, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}
