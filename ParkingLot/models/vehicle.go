package models

type VehicleType int

const (
	Bike VehicleType = iota
	Car
	Truck
)

func (v VehicleType) String() string {
	switch v {
	case Bike:
		return "Bike"
	case Car:
		return "Car"
	case Truck:
		return "Truck"
	default:
		return "Unknown"
	}
}

type Vehicle struct {
	Number string
	Type   VehicleType
}

func NewVehicle(Number string, vType VehicleType) *Vehicle {
	return &Vehicle{
		Number: Number,
		Type:   vType,
	}
}
