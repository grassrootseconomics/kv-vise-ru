-- name: get-session-data
-- Get all keys and values for a session ID.
-- $1: complete prefix (including 20) as bytea
SELECT key, convert_from(value, 'UTF8') as value
FROM kv_vise
WHERE encode(key, 'hex') LIKE encode($1, 'hex') || '%';

-- name: get
-- Get the value for a specific key.
-- $1: session-id
SELECT convert_from(value, 'UTF8') as value
FROM kv_vise
WHERE key = $1;

-- name: get-profile-details-for-sms
-- Get profile details for a session ID.
-- $1: session-id
WITH base AS (SELECT $1::bytea as base_key)
SELECT 
    COALESCE(convert_from(v1.value, 'UTF8'), '') as public_key,
    COALESCE(convert_from(v2.value, 'UTF8'), '') as first_name,
    COALESCE(convert_from(v3.value, 'UTF8'), '') as family_name,
    COALESCE(convert_from(v4.value, 'UTF8'), '') as account_alias,
    COALESCE(convert_from(v5.value, 'UTF8'), '') as language_code
FROM 
    base
    LEFT JOIN kv_vise v1 ON v1.key = base.base_key || '\x0001'::bytea
    LEFT JOIN kv_vise v2 ON v2.key = base.base_key || '\x0003'::bytea
    LEFT JOIN kv_vise v3 ON v3.key = base.base_key || '\x0004'::bytea
    LEFT JOIN kv_vise v4 ON v4.key = base.base_key || '\x0013'::bytea
    LEFT JOIN kv_vise v5 ON v5.key = base.base_key || '\x0015'::bytea;

-- name: address-reverse-lookup
SELECT key, value
FROM kv_vise
WHERE value = convert_to($1, 'UTF8');