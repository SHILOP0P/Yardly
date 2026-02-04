-- 014_fix_booking_confirmation_columns_to_timestamptz.sql
-- Цель: привести колонки подтверждений к timestamptz, чтобы работали COALESCE/CASE и нормальная семантика времени.

DO $$
DECLARE
  col text;
BEGIN
  FOREACH col IN ARRAY ARRAY[
    'handover_confirmed_by_owner_at',
    'handover_confirmed_by_requester_at',
    'return_confirmed_by_owner_at',
    'return_confirmed_by_requester_at'
  ]
  LOOP
    -- если колонки нет — пропускаем (на случай разного состояния схемы)
    IF NOT EXISTS (
      SELECT 1
      FROM information_schema.columns
      WHERE table_name = 'bookings'
        AND column_name = col
    ) THEN
      RAISE NOTICE 'Column %.% does not exist, skipping', 'bookings', col;
      CONTINUE;
    END IF;

    -- если уже timestamptz — пропускаем
    IF EXISTS (
      SELECT 1
      FROM information_schema.columns
      WHERE table_name = 'bookings'
        AND column_name = col
        AND udt_name = 'timestamptz'
    ) THEN
      RAISE NOTICE 'Column %.% already timestamptz, skipping', 'bookings', col;
      CONTINUE;
    END IF;

    -- приводим тип:
    -- 1) пустую строку превращаем в NULL
    -- 2) пытаемся кастануть к timestamptz
    --    если там мусор, то миграция упадёт — поэтому перед ALTER лучше почистить такие строки вручную (см. ниже).
    EXECUTE format(
      'ALTER TABLE bookings
         ALTER COLUMN %I TYPE timestamptz
         USING NULLIF(%I::text, '''')::timestamptz;',
      col, col
    );

    RAISE NOTICE 'Column %.% converted to timestamptz', 'bookings', col;
  END LOOP;
END $$;
