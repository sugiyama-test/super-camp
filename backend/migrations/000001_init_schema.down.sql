DROP TRIGGER IF EXISTS set_updated_at ON meal_plans;
DROP TRIGGER IF EXISTS set_updated_at ON fire_logs;
DROP TRIGGER IF EXISTS set_updated_at ON layouts;
DROP TRIGGER IF EXISTS set_updated_at ON checklist_items;
DROP TRIGGER IF EXISTS set_updated_at ON checklists;
DROP TRIGGER IF EXISTS set_updated_at ON campsites;
DROP TRIGGER IF EXISTS set_updated_at ON gears;
DROP TRIGGER IF EXISTS set_updated_at ON users;
DROP FUNCTION IF EXISTS trigger_set_updated_at;

DROP TABLE IF EXISTS meal_plans;
DROP TABLE IF EXISTS fire_logs;
DROP TABLE IF EXISTS layouts;
DROP TABLE IF EXISTS checklist_items;
DROP TABLE IF EXISTS checklists;
DROP TABLE IF EXISTS campsites;
DROP TABLE IF EXISTS gears;
DROP TABLE IF EXISTS users;
