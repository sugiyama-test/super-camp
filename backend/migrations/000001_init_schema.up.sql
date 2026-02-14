CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- gears
CREATE TABLE gears (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    category     VARCHAR(100) NOT NULL DEFAULT '',
    brand        VARCHAR(255) NOT NULL DEFAULT '',
    weight_grams DOUBLE PRECISION,
    notes        TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_gears_user_id ON gears(user_id);

-- campsites
CREATE TABLE campsites (
    id        BIGSERIAL PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    address   TEXT         NOT NULL DEFAULT '',
    latitude  DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    notes     TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- checklists
CREATE TABLE checklists (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_checklists_user_id ON checklists(user_id);

CREATE TABLE checklist_items (
    id           BIGSERIAL PRIMARY KEY,
    checklist_id BIGINT       NOT NULL REFERENCES checklists(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    is_checked   BOOLEAN      NOT NULL DEFAULT FALSE,
    quantity     INTEGER      NOT NULL DEFAULT 1,
    sort_order   INTEGER      NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_checklist_items_checklist_id ON checklist_items(checklist_id);

-- layouts
CREATE TABLE layouts (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    data        JSONB        NOT NULL DEFAULT '{}',
    campsite_id BIGINT       REFERENCES campsites(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_layouts_user_id ON layouts(user_id);

-- fire_logs
CREATE TABLE fire_logs (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date             DATE         NOT NULL,
    location         VARCHAR(255) NOT NULL DEFAULT '',
    wood_type        VARCHAR(100) NOT NULL DEFAULT '',
    duration_minutes INTEGER      NOT NULL DEFAULT 0,
    notes            TEXT         NOT NULL DEFAULT '',
    temperature      DOUBLE PRECISION,
    campsite_id      BIGINT       REFERENCES campsites(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_fire_logs_user_id ON fire_logs(user_id);
CREATE INDEX idx_fire_logs_date ON fire_logs(date);

-- meal_plans
CREATE TABLE meal_plans (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title     VARCHAR(255) NOT NULL,
    meal_type VARCHAR(50)  NOT NULL DEFAULT 'dinner'
              CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    servings  INTEGER      NOT NULL DEFAULT 2,
    notes     TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_meal_plans_user_id ON meal_plans(user_id);

-- Auto-update updated_at trigger
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at BEFORE UPDATE ON users           FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON gears           FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON campsites       FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON checklists      FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON checklist_items FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON layouts         FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON fire_logs       FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
CREATE TRIGGER set_updated_at BEFORE UPDATE ON meal_plans      FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
