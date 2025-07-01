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
ALTER TABLE investments ADD COLUMN note TEXT DEFAULT NULL;
ALTER TABLE investments ADD COLUMN redemption_policy_type TEXT DEFAULT NULL;

CREATE TABLE investments_summary (
    id TEXT PRIMARY KEY,
    brokerage TEXT DEFAULT NULL,
    type TEXT NOT NULL,
    symbol TEXT,
    bond_index TEXT,
    bond_rate NUMERIC(12,4),
    quantity INTEGER NOT NULL DEFAULT 0,
    average_price NUMERIC(12,4) NOT NULL DEFAULT 0,
    total_value NUMERIC(12,4) NOT NULL DEFAULT 0,
    market_value NUMERIC(12,4) NOT NULL DEFAULT 0,
    cost NUMERIC(12,4) NOT NULL,
    redemption_policy_type TEXT DEFAULT NULL,
    due_date TEXT DEFAULT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_investments_summary_type ON investments_summary(type);
CREATE INDEX idx_investments_summary_brokerage ON investments_summary(brokerage);
CREATE INDEX idx_investments_summary_symbol ON investments_summary(symbol);
CREATE INDEX idx_investments_summary_bond_index ON investments_summary(bond_index);
CREATE INDEX idx_investments_summary_redemption_policy_type ON investments_summary(redemption_policy_type);