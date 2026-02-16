package models

type GpsInput struct {
	DroneId     string  `json:"drone_id"`
	Timestamp   string  `json:"timestamp"`
	Latitude   float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	Heading     float64 `json:"heading"`
	Pitch       float64 `json:"pitch"`
	Roll        float64 `json:"roll"`
	Speed       float64 `json:"speed"`
	ClimbRate   float64 `json:"climb_rate"`
	AngularRate float64 `json:"angular_rate"`
}
