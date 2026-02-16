package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/JacksonDahl2/drone-platform/cmd/shared/models"
	sqlc "github.com/JacksonDahl2/drone-platform/internal/platform/db/sqlc"
)

type Processor struct {
	queries  *sqlc.Queries
}

func NewProcessor(db *sql.DB) *Processor {
	queries := sqlc.New(db)
	return &Processor{
		queries: queries,
	}
}

func (p *Processor) ProcessGps(ctx context.Context, msg []byte) error {
	// unmarshal
	var input models.GpsInput
	err := json.Unmarshal(msg, &input)
	if err != nil {
		log.Printf("failed to unmarshal data")
		return err
	}

	time, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil {
		log.Printf("failed to parse timestamp")
		return err
	}

	params := sqlc.InsertGpsParams{
		DroneID: input.DroneId,
		Time: time,
		Latitude: input.Latitude,
		Longitude: input.Longitude,
		Altitude: input.Altitude,
		Heading: input.Heading,
		Pitch: input.Pitch,
		Roll: input.Roll,
		Speed: input.Speed,
		ClimbRate: input.ClimbRate,
		AngularRate: input.AngularRate,
	}

	if err := p.queries.InsertGps(ctx, params); err != nil {
		log.Printf("Failed to add to db")
		return err
	}
	return nil
}

func (p *Processor) ProcessState(ctx context.Context, msg []byte) error {
	// unmarshal
	var input models.StateInput
	err := json.Unmarshal(msg, &input)
	if err != nil {
		log.Printf("failed to unmarshal data")
		return err
	}

	time, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil {
		log.Printf("failed to parse timestamp")
		return err
	}

	params := sqlc.InsertStateParams {
		DroneID: input.DroneId,
		Time: time,
		Status: input.Status,
		BatteryPct: input.BatteryPct,
		Voltage: input.Voltage,
		Connected: input.Connected,
		FlightMode: input.FlightMode,
	}

	if err := p.queries.InsertState(ctx, params); err != nil {
		log.Printf("failed to add to db")
		return err
	}
	return nil
}

func (p *Processor) ProcessEvent(ctx context.Context, msg []byte) error{
	// unmarshal
	var input models.EventInput
	err := json.Unmarshal(msg, &input)
	if err != nil {
		log.Printf("failed to unmarshal data")
		return err
	}

	time, err := time.Parse(time.RFC3339, input.Timestamp)
	if err != nil {
		log.Printf("failed to parse timestamp")
		return err
	}

	payloadBytes, err := json.Marshal(input.Payload)
	if err != nil {
		return err
	}

	params := sqlc.InsertEventParams {
		DroneID: input.DroneId,
		Time: time,
		EventType: input.EventType,
		Payload: payloadBytes,
	}

	if err := p.queries.InsertEvent(ctx, params); err != nil {
		log.Printf("failed to add to db")
		return err
	}
	return nil
}