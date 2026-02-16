package models

type GpsInput struct {
	DroneID     string  `json:"drone_id"`
	Timestamp   string  `json:"timestamp"`
	Lattitude   float64 `json:"lattitude"`
	Longitude   float64 `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	Heading     float64 `json:"heading"`
	Pitch       float64 `json:"pitch"`
	Roll        float64 `json:"roll"`
	Speed       float64 `json:"speed"`
	ClimbRate   float64 `json:"climb_rate"`
	AngularRate float64 `json:"angular_rate"`
}
