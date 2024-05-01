-- Table: booking
DROP TABLE IF EXISTS "bookings";

-- Table: event
DROP TABLE IF EXISTS "events";

-- Table: room
DROP TABLE IF EXISTS "rooms";
DROP TYPE IF EXISTS rooms_type_enum;

-- Table: user_cinema
DROP TABLE IF EXISTS "users_cinemas";

-- Table: cinema
DROP TABLE IF EXISTS "cinemas";

-- Table: user
DROP TABLE IF EXISTS "users";
DROP TYPE IF EXISTS users_account_type_enum;
DROP TYPE IF EXISTS users_status_enum;
DROP TYPE IF EXISTS users_role_enum;

-- Table: address
DROP TABLE IF EXISTS "addresses";
