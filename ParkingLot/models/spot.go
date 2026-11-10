package models

type ParkingSpot struct {
	ID      string
	Floor   int
	Type    VehicleType
	Vehicle *Vehicle
}

func NewParkingSpot(id string, floor int, vType VehicleType) *ParkingSpot {
	return &ParkingSpot{
		ID:    id,
		Floor: floor,
		Type:  vType,
	}
}

func (s *ParkingSpot) IsFree() bool {
	return s.Vehicle == nil
}

func (s *ParkingSpot) CanFit(vType VehicleType) bool {
	return s.IsFree() && s.Type == vType
}

func (s *ParkingSpot) Park(v *Vehicle) {
	s.Vehicle = v
}

func (s *ParkingSpot) UnPark() {
	s.Vehicle = nil
}
