CREATE TABLE dividends (
    id SERIAL PRIMARY KEY NOT NULL,
    ticker VARCHAR NOT NULL,
    activity_type VARCHAR NOT NULL,
    dollar_value MONEY,
    received_on TIMESTAMP NOT NULL
);

CREATE INDEX ix_dividends_ticker ON dividends(ticker);
CREATE INDEX ix_dividends_activity_type ON dividends(activity_type);
CREATE INDEX ix_dividends_date ON dividends(received_on);
CREATE UNIQUE INDEX ix_dividends_ticker_activity_type_received_on ON dividends(ticker, activity_type, received_on);