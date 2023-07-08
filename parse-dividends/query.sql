-- name: CreateDividends :exec
INSERT INTO dividends (ticker, dollar_value, received_on) VALUES ($1, $2, $3) ON CONFLICT (received_on) DO NOTHING;