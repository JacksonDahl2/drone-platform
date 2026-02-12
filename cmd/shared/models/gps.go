package models

type GpsInput struct {
	DroneID     string `json:"drone_id"`
	Timestamp   string `json:"timestamp"`
	Lattitude   string `json:"lattitude"`
	Longitude   string `json:"longitude"`
	Altitude    string `json:"altitude"`
	Heading     string `json:"heading"`
	Pitch       string `json:"pitch"`
	Roll        string `json:"roll"`
	Speed       string `json:"speed"`
	ClimbRate   string `json:"climb_rate"`
	AngularRate string `json:"angular_rate"`
}
