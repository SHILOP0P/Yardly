-- У одного requester не может быть 2 активных заявок на один item для buy/give
CREATE UNIQUE INDEX IF NOT EXISTS uniq_active_transfer_per_item_requester
ON bookings (item_id, requester_id)
WHERE type IN ('buy', 'give')
  AND status IN ('requested', 'approved', 'handover_pending');
