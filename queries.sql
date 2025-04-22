-- name: get-session-data
-- Get all keys and values for a session ID.
-- $1: complete prefix (including 20) as bytea
SELECT key, convert_from(value, 'UTF8') as value
FROM kv_vise
WHERE encode(key, 'hex') LIKE encode($1, 'hex') || '%';

-- name: get
-- Get the value for a specific key.
-- $1: complete key as bytea
SELECT convert_from(value, 'UTF8') as value
FROM kv_vise
WHERE key = $1; 