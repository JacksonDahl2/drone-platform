package models

type EventInput struct {
	DroneId   string         `json:"drone_id"`
	Timestamp string         `json:"timestamp"`
	EventType string         `json:"event_type"`
	Payload   map[string]any `json:"payload"`
}
