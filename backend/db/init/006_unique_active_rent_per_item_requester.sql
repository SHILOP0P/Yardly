BEGIN;

CREATE UNIQUE INDEX IF NOT EXISTS uq_bookings_item_requester_active_rent
ON bookings (item_id, requester_id)
WHERE type = 'rent' AND status IN ('requested', 'approved', 'in_use');

COMMIT;
