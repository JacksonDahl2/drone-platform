CREATE INDEX idx_gps_drone_id_time ON gps (drone_id, time DESC);
CREATE INDEX idx_state_drone_id_time ON state (drone_id, time DESC);
CREATE INDEX idx_events_drone_id_time ON events (drone_id, time DESC);

