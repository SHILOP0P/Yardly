ALTER TABLE bookings
ADD CONSTRAINT bookings_dates_by_type_chk
CHECK (
  (type = 'rent' AND start_at IS NOT NULL AND end_at IS NOT NULL)
  OR
  (type IN ('buy','give') AND start_at IS NULL AND end_at IS NULL)
);
