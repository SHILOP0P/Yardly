BEGIN
ALTER TABLE bookings
  ADD COLUMN handover_confirmed_by_owner_at timestamptz NULL,
  ADD COLUMN handover_confirmed_by_requester_at timestamptz NULL,
  ADD COLUMN return_confirmed_by_owner_at timestamptz NULL,
  ADD COLUMN return_confirmed_by_requester_at timestamptz NULL;
COMMIT