package models

import "time"

type Ticket struct {
	ID        string
	Vehicle   *Vehicle
	Spot      *ParkingSpot
	EntryTime time.Time
	ExitTime  time.Time
	Fee       float64
}

func NewTicket(id string, v *Vehicle, spot *ParkingSpot, entry time.Time) *Ticket {
	return &Ticket{
		ID:        id,
		Vehicle:   v,
		Spot:      spot,
		EntryTime: entry,
	}
}

func (t *Ticket) Duration() time.Duration {
	return t.ExitTime.Sub(t.EntryTime)
}
