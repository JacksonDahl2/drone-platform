package models

type StateInput struct {
	DroneId    string  `json:"drone_id"`
	Timestamp  string  `json:"timestamp"`
	Status     string  `json:"status"` // "idle", "armed", "flying", "landing", "error"
	BatteryPct float64 `json:"battery_pct"` // 0 - 100
	Voltage    float64 `json:"voltage"` // 
	Connected  bool    `json:"connected"` // link to ground station
	FlightMode string  `json:"flight_mode"` // "manual", "auto", "loiter"
}
