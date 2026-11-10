package models

import "fmt"

type Floor struct {
	Number int
	Spots  []*ParkingSpot
}

func NewFloor(number, bikeSpots, carSpots, truckSpots int) *Floor {
	f := &Floor{Number: number}
	f.addSpots(Bike, bikeSpots)
	f.addSpots(Car, carSpots)
	f.addSpots(Truck, truckSpots)
	return f
}

func (f *Floor) addSpots(vType VehicleType, count int) {
	prefix := vType.String()[:1] // "B", "C" or "T"
	for i := 1; i <= count; i++ {
		id := fmt.Sprintf("F%d-%s%d", f.Number, prefix, i)
		f.Spots = append(f.Spots, NewParkingSpot(id, f.Number, vType))
	}
}

func (f *Floor) FreeCount(vType VehicleType) int {
	count := 0
	for _, s := range f.Spots {
		if s.CanFit(vType) {
			count++
		}
	}
	return count
}
