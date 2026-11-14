package services

import (
	"parkinglot/models"
	"time"
)

// EntryGate is one boom barrier where vehicles come in.
// Many gates share the SAME *ParkingLotService.
type EntryGate struct {
	ID  string
	lot *ParkingLotService
}

func NewEntryGate(id string, lot *ParkingLotService) *EntryGate {
	return &EntryGate{
		ID:  id,
		lot: lot,
	}
}

func (g *EntryGate) Enter(v *models.Vehicle, at time.Time) (*models.Ticket, error) {
	return g.lot.Park(v, at)
}

// ExitGate is one boom barrier where vehicles leave and pay.
type ExitGate struct {
	ID  string
	lot *ParkingLotService
}

func NewExitGate(id string, lot *ParkingLotService) *ExitGate {
	return &ExitGate{
		ID:  id,
		lot: lot,
	}
}

func (g *ExitGate) Exit(ticketID string, at time.Time) (*models.Ticket, error) {
	return g.lot.Unpark(ticketID, at)
}
