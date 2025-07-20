--sqlite d1 (cloudflare.com)
CREATE TABLE investments (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    symbol TEXT,
    bond_index TEXT,
    bond_rate NUMERIC(12,4),
    quantity NUMERIC(12,6) NOT NULL DEFAULT 0,
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
    investment_id TEXT DEFAULT NULL,
    last_operation_date TEXT NOT NULL, 
    brokerage TEXT DEFAULT NULL,
    type TEXT NOT NULL,
    symbol TEXT,
    bond_index TEXT,
    bond_rate NUMERIC(12,4),
    quantity NUMERIC(12,4) NOT NULL DEFAULT 0,
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
CREATE INDEX idx_investments_summary_investment_id ON investments_summary(investment_id);
CREATE INDEX idx_investments_summary_last_operation_date ON investments_summary(last_operation_date);


insert into investments_summary(
    id,
    investment_id,
    last_operation_date,
    brokerage,
    type,
    symbol,
    bond_index,
    bond_rate,
    quantity,
    average_price,
    total_value,
    cost,
    redemption_policy_type,
    due_date
) select
    id,
    id as investment_id,
    operation_date,
    brokerage,
    type,
    symbol,
    bond_index,
    bond_rate,
    quantity,
    unit_price,
    total_value,
    cost,
    redemption_policy_type,
    due_date
  from investments;


ALTER TABLE investments ADD COLUMN sell_investment_id TEXT DEFAULT NULL;


CREATE TABLE investments_summary_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    investment_id TEXT DEFAULT NULL,
    last_operation_date TEXT NOT NULL,
    operation_month INTEGER NOT NULL,
    operation_year INTEGER NOT NULL, 
    brokerage TEXT DEFAULT NULL,
    type TEXT NOT NULL,
    symbol TEXT,
    bond_index TEXT,
    bond_rate NUMERIC(12,4),
    quantity NUMERIC(12,4) NOT NULL DEFAULT 0,
    average_price NUMERIC(12,4) NOT NULL DEFAULT 0,
    total_value NUMERIC(12,4) NOT NULL DEFAULT 0,
    market_value NUMERIC(12,4) NOT NULL DEFAULT 0,
    cost NUMERIC(12,4) NOT NULL,
    redemption_policy_type TEXT DEFAULT NULL,
    due_date TEXT DEFAULT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

ALTER TABLE investments_summary_history ADD COLUMN investment_summary_id TEXT NOT NULL;


INSERT INTO investments_summary_history(
    investment_id,
    last_operation_date,
    operation_month,
    operation_year, 
    brokerage,
    type,
    symbol,
    bond_index,
    bond_rate,
    quantity,
    average_price,
    total_value,
    market_value,
    cost,
    redemption_policy_type,
    due_date,
    investment_summary_id   
) SELECT 
    investment_id,
    last_operation_date,
    strftime('%m', last_operation_date) as operation_month,
    strftime('%Y', last_operation_date) as operation_year,
    brokerage,
    type,
    symbol,
    bond_index,
    bond_rate,
    quantity,
    average_price,
    total_value,
    market_value,
    cost,
    redemption_policy_type,
    due_date,
    id as investment_summary_id 
  FROM investments_summary
  WHERE id = ?;


ALTER TABLE investments ADD COLUMN pnl NUMERIC(12, 4) NOT NULL DEFAULT 0;
ALTER TABLE investments ADD COLUMN average_selling_price NUMERIC(12, 4) NOT NULL DEFAULT 0;


