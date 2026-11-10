package strategies

import "parkinglot/models"

// This SpotAllocationStrategy decided which free spot to allocate the vehicle.
type SpotAllocationStrategy interface {
	FindSpot(floors []*models.Floor, vType models.VehicleType) *models.ParkingSpot
}

// NearestSpotStrategy : lowest floor first, first matching spot on it.
type NearestSpotStrategy struct{}

func (NearestSpotStrategy) FindSpot(floors []*models.Floor, vType models.VehicleType) *models.ParkingSpot {
	for _, floor := range floors {
		for _, spot := range floor.Spots {
			if spot.CanFit(vType) {
				return spot
			}
		}
	}

	return nil // no spot found
}

// TopFloorFirst: fill top floors first keep ground floors for VIP's
type TopFloorFirstStrategy struct{}

func (TopFloorFirstStrategy) FindSpot(floors []*models.Floor, vType models.VehicleType) *models.ParkingSpot {
	for i := len(floors); i >= 0; i-- {
		for _, spot := range floors[i].Spots {
			if spot.CanFit(vType) {
				return spot
			}
		}
	}

	return nil // no spot found
}
