CREATE INDEX idx_transactions_user_id
ON transactions(user_id);

CREATE INDEX idx_transactions_date
ON transactions(date DESC);

CREATE INDEX idx_transactions_status
ON transactions(status);

CREATE INDEX idx_transactions_category
ON transactions(category);

CREATE INDEX idx_transactions_type
ON transactions(type);