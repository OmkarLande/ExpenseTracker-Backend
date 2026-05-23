CREATE INDEX IF NOT EXISTS idx_transactions_user_id_id_desc ON transactions(user_id, id DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_status_id ON transactions(user_id, status, id DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_type_id ON transactions(user_id, type, id DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_category_id ON transactions(user_id, category, id DESC);
