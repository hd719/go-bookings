package repository

import (
	"time"

	"github.com/hd719/go-bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesForRoom(start, end time.Time, roomId int) (bool, error)
}
