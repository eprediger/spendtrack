CREATE TABLE expense_item_discounts (
    expense_id UUID NOT NULL,
    item_position INTEGER NOT NULL,
    position INTEGER NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC NOT NULL CHECK (amount > 0),
    PRIMARY KEY (expense_id, item_position, position),
    FOREIGN KEY (expense_id, item_position) REFERENCES expense_items(expense_id, position) ON DELETE CASCADE
);
