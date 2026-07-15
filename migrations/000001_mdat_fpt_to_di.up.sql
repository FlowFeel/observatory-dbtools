-- Migration 001: Fix _MDAT routing mismatch (SMW 6.0.1 bug)
-- Copies _MDAT entries from smw_fpt_mdat to smw_di_time where missing.
-- Property ID 29 = _MDAT (Modification date)
INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey)
SELECT s_id, 29, o_serialized, o_sortkey FROM smw_fpt_mdat
WHERE s_id NOT IN (SELECT s_id FROM smw_di_time WHERE p_id = 29);
