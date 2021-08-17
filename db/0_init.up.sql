CREATE TABLE prices (
    ticker VARCHAR,
    price_date VARCHAR,
    open VARCHAR, 
    high VARCHAR, 
    low VARCHAR, 
    close VARCHAR,
    PRIMARY KEY(ticker, price_date)
);

CREATE TABLE users (
    id uuid NOT NULL,
    name varchar NOT NULL,
    email varchar,
    username varchar NOT NULL,
    password varchar NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE tokens (
    id uuid NOT NULL,
    token varchar NOT NULL,
    user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(id)
);

INSERT INTO users (id, name, username, password) VALUES ('0622dea2-ee79-4aa9-8560-b3ba5a09fa26', 'Super User', 'admin', '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918');
