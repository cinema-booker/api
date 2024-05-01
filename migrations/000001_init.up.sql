-- Table: addresses

CREATE TABLE "addresses" (
  "id" SERIAL PRIMARY KEY,
  "country" VARCHAR(255),
  "city" VARCHAR(255),
  "zip_code" VARCHAR(20),
  "street" VARCHAR(255),
  "longitude" DECIMAL(10, 8),
  "latitude" DECIMAL(11, 8)
);

-- Table: users

CREATE TYPE users_role_enum AS ENUM ('ADMIN', 'MANAGER', 'VIEWER');
CREATE TYPE users_status_enum AS ENUM ('ACTIVE', 'ARCHIVED');
CREATE TYPE users_account_type_enum AS ENUM ('CREDENTIALS', 'GOOGLE');

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "first_name" VARCHAR(255),
  "last_name" VARCHAR(255),
  "email" VARCHAR(255) UNIQUE,
  "password" VARCHAR(255),
  "role" users_role_enum,
  "status" users_status_enum,
  "account_type" users_account_type_enum
);

-- Table: cinemas

CREATE TABLE "cinemas" (
  "id" SERIAL PRIMARY KEY,
  "address_id" INTEGER REFERENCES "addresses"("id"),
  "name" VARCHAR(255),
  "description" TEXT,
  "images" TEXT[]
);

-- Table: users_cinemas

CREATE TABLE "users_cinemas" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER REFERENCES "users"("id"),
  "cinema_id" INTEGER REFERENCES "cinemas"("id")
);

-- Table: rooms

CREATE TYPE rooms_type_enum AS ENUM ('SMALL', 'MEDIUM', 'LARGE');

CREATE TABLE "rooms" (
  "id" SERIAL PRIMARY KEY,
  "cinema_id" INTEGER REFERENCES "cinemas"("id"),
  "number" VARCHAR(255),
  "type" rooms_type_enum
);

-- Table: event

CREATE TABLE "events" (
  "id" SERIAL PRIMARY KEY,
  "room_id" INTEGER REFERENCES "rooms"("id"),
  "movie_id" INTEGER,
  "price" DECIMAL(10, 2),
  "starts_at" TIMESTAMP,
  "ends_at" TIMESTAMP
);

-- Table: bookings

CREATE TABLE "bookings" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER REFERENCES "users"("id"),
  "event_id" INTEGER REFERENCES "events"("id"),
  "place" VARCHAR(255)
);
