CREATE TABLE expense_items (
    expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    description TEXT NOT NULL,
    quantity NUMERIC NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC NOT NULL,
    line_total NUMERIC NOT NULL,
    PRIMARY KEY (expense_id, position)
);
