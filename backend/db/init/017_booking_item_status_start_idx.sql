CREATE INDEX IF NOT EXISTS bookings_item_type_status_start_idx
ON bookings (item_id, type, status, start_at);
