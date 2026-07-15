-- Rollback: remove DI entries that exist in FPT but were added by this migration.
-- WARNING: This is destructive only for entries that had no DI entry before migration.
-- We track by matching (s_id, o_serialized) pairs that exist in both tables.
DELETE di FROM smw_di_time di
INNER JOIN smw_fpt_mdat fpt ON di.s_id = fpt.s_id
    AND di.o_serialized = fpt.o_serialized
WHERE di.p_id = 29;
