CREATE TABLE accounts (
    user_id INTEGER,
    balance NUMERIC NOT NULL DEFAULT 0.00,
    PRIMARY KEY (user_id),
    CONSTRAINT not_negative_balance CHECK (balance >= 0.00)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    debit_user_id INTEGER,
    credit_user_id INTEGER,
    amount NUMERIC NOT NULL,
    description VARCHAR(255),
    CONSTRAINT fk_debit_user
        FOREIGN KEY (debit_user_id)
            REFERENCES accounts(user_id)
            ON DELETE SET NULL,
    CONSTRAINT fk_credit_user
        FOREIGN KEY (credit_user_id)
            REFERENCES accounts(user_id)
            ON DELETE SET NULL,
    CONSTRAINT debit_or_credit
        CHECK (debit_user_id IS NOT NULL OR credit_user_id IS NOT NULL),
    CONSTRAINT positive_amount CHECK (amount > 0.00)
);
