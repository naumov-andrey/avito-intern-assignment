CREATE TABLE accounts (
    user_id INTEGER PRIMARY KEY,
    balance NUMERIC NOT NULL DEFAULT 0.00 CHECK (balance >= 0.00)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    account_id INTEGER REFERENCES accounts ON DELETE CASCADE,
    amount NUMERIC NOT NULL,
    description VARCHAR(255)
);
