-- name: CreateDividends :exec
INSERT INTO dividends (ticker, dollar_value, activity_type, received_on) VALUES ($1, $2, $3, $4) ON CONFLICT (received_on) DO NOTHING;