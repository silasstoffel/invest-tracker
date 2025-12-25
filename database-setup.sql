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



CREATE TABLE symbol_details (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL, 
    segment TEXT NOT NULL,
    sub_segment TEXT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_symbol_details_type ON symbol_details(type);
CREATE INDEX idx_symbol_details_segment ON symbol_details(segment);

insert into symbol_details (id, "type", segment, sub_segment) values
('B3SA3', '#', '#', null),
('BARI11','#', '#', null),
('BBDC3','#', '#', null),
('BPAN4','#', '#', null),
('BRSR6','#', '#', null),
('BTHF11','#', '#', null),
('BTLG11','#', '#', null),
('CLIN11','#', '#', null),
('COGN3','#', '#', null),
('CVCB3','#', '#', null),
('FIGS11','#', '#', null),
('FIIB11','#', '#', null),
('HGBS11','#', '#', null),
('HGFF11','#', '#', null),
('HGLG11','#', '#', null),
('HGRU11','#', '#', null),
('HLOG11','#', '#', null),
('HSML11','#', '#', null),
('HYPE3','#', '#', null),
('IRDM11','#', '#', null),
('IRIM11','#', '#', null),
('ITRI11','#', '#', null),
('JSRE11','#', '#', null),
('KCRE11','#', '#', null),
('KNRI11','#', '#', null),
('LOGG3','#', '#', null),
('LREN3','#', '#', null),
('LVBI11','#', '#', null),
('MGLU3','#', '#', null),
('MILS3','#', '#', null),
('NEOE3','#', '#', null),
('PCIP11','#', '#', null),
('PLCR11','#', '#', null),
('PMLL11','#', '#', null),
('PORD11','#', '#', null),
('PSEC11','#', '#', null),
('RBRF11','#', '#', null),
('RBRL11','#', '#', null),
('RBRX11','#', '#', null),
('RECR11','#', '#', null),
('RECV3','#', '#', null),
('RFOF11','#', '#', null),
('RVBI11','#', '#', null),
('SOJA3','#', '#', null),
('TRXF11','#', '#', null),
('VALE3','#', '#', null),
('VAMO3','#', '#', null),
('VILG11','#', '#', null),
('VINO11','#', '#', null),
('WHGR11','#', '#', null),
('WIZC3','#', '#', null),
('XPIN11','#', '#', null);

