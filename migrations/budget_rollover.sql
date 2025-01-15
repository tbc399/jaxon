CREATE table budget_periods (
    id VARCHAR(22) PRIMARY KEY,
    user_id VARCHAR(22) REFERENCES users(id),
    start TIMESTAMP,
    "end" TIMESTAMP,
    created_at TIMESTAMP
);

CREATE INDEX budget_periods_id_index ON budget_periods(id);
CREATE INDEX budget_periods_user_id_index ON budget_periods(user_id);
