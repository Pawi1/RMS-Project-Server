PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: Users
CREATE TABLE IF NOT EXISTS Users (
id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
user_name TEXT UNIQUE,
email TEXT UNIQUE,
password_hash TEXT,
role TEXT NOT NULL CHECK (role IN ('admin','editor','operator','mechanic','viewer','client')) DEFAULT 'client',
created_at TEXT DEFAULT (datetime('now')),
updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_users_email ON Users(email);
CREATE INDEX IF NOT EXISTS idx_users_user_name ON Users(user_name);

CREATE TRIGGER IF NOT EXISTS trg_users_updated_at
AFTER UPDATE ON Users
FOR EACH ROW
BEGIN
UPDATE Users SET updated_at = datetime('now') WHERE id = OLD.id;
END;

-- Table: Cars
CREATE TABLE IF NOT EXISTS Cars (
id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
id_owner INTEGER REFERENCES Users (id) ON DELETE SET NULL,
vin TEXT UNIQUE,
plate_number TEXT,
make TEXT,
model TEXT,
year INTEGER,
last_mileage INTEGER CHECK (last_mileage >= 0),
fuel_type TEXT CHECK (fuel_type IN ('petrol','diesel','hybrid','electric','gas','other')),
engine_capacity NUMERIC,
engine_type TEXT,
default_hourly_rate NUMERIC,
notes TEXT,
created_at TEXT DEFAULT (datetime('now')),
updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_cars_owner ON Cars(id_owner);
CREATE INDEX IF NOT EXISTS idx_cars_plate ON Cars(plate_number);

CREATE TRIGGER IF NOT EXISTS trg_cars_updated_at
AFTER UPDATE ON Cars
FOR EACH ROW
BEGIN
UPDATE Cars SET updated_at = datetime('now') WHERE id = OLD.id;
END;

-- Table: Mechanic
CREATE TABLE IF NOT EXISTS Mechanic (
id INTEGER PRIMARY KEY REFERENCES Users (id) ON DELETE CASCADE,
hourly_rate NUMERIC,
created_at TEXT DEFAULT (datetime('now')),
updated_at TEXT DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_mechanic_updated_at
AFTER UPDATE ON Mechanic
FOR EACH ROW
BEGIN
UPDATE Mechanic SET updated_at = datetime('now') WHERE id = OLD.id;
END;

-- Table: Persons
CREATE TABLE IF NOT EXISTS Persons (
user_id INTEGER PRIMARY KEY REFERENCES Users (id) ON DELETE CASCADE NOT NULL,
first_name TEXT,
last_name TEXT,
phone TEXT CHECK (
phone GLOB '+48[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]'
OR phone GLOB '[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]'
OR phone GLOB '[0-9][0-9][0-9]-[0-9][0-9][0-9]-[0-9][0-9][0-9]'
),
created_at TEXT DEFAULT (datetime('now')),
updated_at TEXT DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_persons_updated_at
AFTER UPDATE ON Persons
FOR EACH ROW
BEGIN
UPDATE Persons SET updated_at = datetime('now') WHERE user_id = OLD.user_id;
END;

-- Table: RepairOrders
CREATE TABLE IF NOT EXISTS RepairOrders (
id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
id_car INTEGER REFERENCES Cars (id) ON DELETE CASCADE NOT NULL,
title TEXT,
description TEXT,
total_cost NUMERIC,
created_at TEXT DEFAULT (datetime('now')),
updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_repairorders_car ON RepairOrders(id_car);

CREATE TRIGGER IF NOT EXISTS trg_repairorders_updated_at
AFTER UPDATE ON RepairOrders
FOR EACH ROW
BEGIN
UPDATE RepairOrders SET updated_at = datetime('now') WHERE id = OLD.id;
END;

-- Table: RepairTasks
CREATE TABLE IF NOT EXISTS RepairTasks (
id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
order_id INTEGER REFERENCES RepairOrders (id) ON DELETE CASCADE NOT NULL,
mechanic_id INTEGER REFERENCES Mechanic (id) NOT NULL,
title TEXT NOT NULL,
description TEXT,
hours REAL,
image_path TEXT,
created_at TEXT DEFAULT (datetime('now')),
updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_repairtasks_order ON RepairTasks(order_id);
CREATE INDEX IF NOT EXISTS idx_repairtasks_mechanic ON RepairTasks(mechanic_id);

CREATE TRIGGER IF NOT EXISTS trg_repairtasks_updated_at
AFTER UPDATE ON RepairTasks
FOR EACH ROW
BEGIN
UPDATE RepairTasks SET updated_at = datetime('now') WHERE id = OLD.id;
END;

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;