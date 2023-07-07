CREATE TABLE dividends (
    id SERIAL PRIMARY KEY NOT NULL,
    ticker VARCHAR NOT NULL,
    dollar_value MONEY,
    received_on TIMESTAMP NOT NULL
);

CREATE INDEX ix_dividends_ticker ON dividends(ticker);
CREATE INDEX ix_dividends_date ON dividends(date);