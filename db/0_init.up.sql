CREATE TABLE prices (
    ticker VARCHAR,
    price_date VARCHAR,
    open VARCHAR, 
    high VARCHAR, 
    low VARCHAR, 
    close VARCHAR,
    PRIMARY KEY(ticker, price_date)
);