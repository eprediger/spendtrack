CREATE TABLE expenses (
    id                      UUID PRIMARY KEY,
    occurred_at             TIMESTAMPTZ NOT NULL,
    occurred_offset_minutes SMALLINT NOT NULL,
    amount                  NUMERIC NOT NULL CHECK (amount > 0),
    currency                CHAR(3) NOT NULL,
    description             TEXT NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
