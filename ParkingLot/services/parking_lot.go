package services

import (
	"errors"
	"fmt"
	"parkinglot/models"
	"parkinglot/strategies"
	"sync"
	"time"
)

var (
	ErrNoSpotAvailable = errors.New("no spot available for this vehicle type")
	ErrAlreadyParked   = errors.New("vehicle is already parked")
	ErrTicketNotFound  = errors.New("ticket not found")
)

type ParkingLotService struct {
	mu             sync.Mutex // many gates could call Park/Unpark at the same time
	name           string
	floors         []*models.Floor
	spotStrategy   strategies.SpotAllocationStrategy // we should declare the interfaces in service struct.
	feeStrategy    strategies.FeeStrategy
	activeTickets  map[string]*models.Ticket // ticketID -> ticket
	parkedVehicles map[string]string         // vehicle number -> ticketID
	ticketCounter  int                       // It is used in generating ticketId like example T1, T2..
}

// For interfaces you shouldn't keep star.
func NewParkingLotService(name string, floors []*models.Floor, spotStrategy strategies.SpotAllocationStrategy, feeStrategy strategies.FeeStrategy) *ParkingLotService {
	return &ParkingLotService{
		name:           name,
		floors:         floors,
		spotStrategy:   spotStrategy,
		feeStrategy:    feeStrategy,
		activeTickets:  make(map[string]*models.Ticket),
		parkedVehicles: make(map[string]string),
	}
}

func (p *ParkingLotService) Park(v *models.Vehicle, entryTime time.Time) (*models.Ticket, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 1. Validate
	if _, ok := p.parkedVehicles[v.Number]; ok {
		return nil, ErrAlreadyParked
	}

	// 2. Ask the strategy for a spot
	spot := p.spotStrategy.FindSpot(p.floors, v.Type)
	if spot == nil {
		return nil, ErrNoSpotAvailable
	}

	// 3. Update state
	spot.Park(v)
	p.ticketCounter++
	ticketID := fmt.Sprint("T%d", p.ticketCounter)
	ticket := models.NewTicket(ticketID, v, spot, entryTime)

	p.activeTickets[ticketID] = ticket
	p.parkedVehicles[v.Number] = ticketID

	// 4. Return
	return ticket, nil

}

// Unpark: validate ticket -> set exit time -> fee (strategy) -> free spot -> return
func (p *ParkingLotService) Unpark(ticketID string, exitTime time.Time) (*models.Ticket, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ticket, ok := p.activeTickets[ticketID]
	if !ok {
		return nil, ErrTicketNotFound
	}

	ticket.ExitTime = exitTime
	ticket.Fee = p.feeStrategy.CalculateFee(ticket)
	ticket.Spot.UnPark()

	delete(p.activeTickets, ticketID)
	delete(p.parkedVehicles, ticket.Vehicle.Number)

	return ticket, nil
}

func (p *ParkingLotService) ShowAvailability() {
	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Printf("--- %s availability ---\n", p.name)
	types := []models.VehicleType{models.Bike, models.Car, models.Truck}
	for _, f := range p.floors {
		fmt.Printf("Floor %d:", f.Number)
		for _, t := range types {
			fmt.Printf(" %s=%d", t, f.FreeCount(t))
		}
		fmt.Println()
	}
}
