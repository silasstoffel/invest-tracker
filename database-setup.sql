--sqlite d1 (cloudflare.com)
CREATE TABLE investments (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    symbol TEXT,
    bond_index TEXT,
    bond_rate NUMERIC(12,4),
    quantity INTEGER NOT NULL DEFAULT 0,
    unit_price NUMERIC(12,4) NOT NULL DEFAULT 0,
    total_value NUMERIC(12,4) NOT NULL,
    cost NUMERIC(12,4) NOT NULL,
    operation_type TEXT NOT NULL,
    operation_date TEXT NOT NULL,
    operation_year INTEGER NOT NULL,
    operation_month INTEGER NOT NULL,
    due_date TEXT DEFAULT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);


CREATE INDEX idx_investments_type ON investments(type);
CREATE INDEX idx_investments_symbol ON investments(symbol);
CREATE INDEX idx_investments_bond_index ON investments(bond_index);
CREATE INDEX operation_year ON investments(operation_year);
CREATE INDEX operation_month ON investments(operation_month);

ALTER TABLE investments ADD COLUMN brokerage TEXT DEFAULT NULL;