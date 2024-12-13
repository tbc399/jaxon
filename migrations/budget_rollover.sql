CREATE table budget_rollovers (
    id VARCAHR(22) PRIMARY KEY,
    user_id VARCHAR(22),
    year INTEGER,
    month INTEGER,
    created_at TIMESTAMP
);

CREATE INDEX budget_rollover_id_index ON budget_rollovers(id);
CREATE INDEX budget_rollover_user_id_index ON budget_rollovers(user_id);
