package strategies

import (
	"math"
	"parkinglot/models"
)

type FeeStrategy interface {
	CalculateFee(ticket *models.Ticket) float64
}

// HourlyFeeStrategy: rate per hour by vehicle type, rounded up, minimum hour.
type HourlyFeeStrategy struct {
	Rates map[models.VehicleType]float64
}

func NewHourlyStrategy() *HourlyFeeStrategy {
	return &HourlyFeeStrategy{
		Rates: map[models.VehicleType]float64{
			models.Bike:  10,
			models.Car:   20,
			models.Truck: 50,
		},
	}
}

func (h *HourlyFeeStrategy) CalculateFee(ticket *models.Ticket) float64 {
	hours := int(math.Ceil(ticket.Duration().Hours())) // 2.5h -> 3
	if hours < 1 {
		hours = 1
	}

	return float64(hours) * h.Rates[ticket.Vehicle.Type]
}

// FlatFeeStrategy: same amount no matter how long
type FlatFeeStrategy struct {
	Amount float64
}

func (f FlatFeeStrategy) CalculateFee(ticket *models.Ticket) float64 {
	return f.Amount
}
