-- Table: addresses

CREATE TABLE "addresses" (
  "id" SERIAL PRIMARY KEY,
  "country" VARCHAR(255) NOT NULL,
  "city" VARCHAR(255) NOT NULL,
  "zip_code" VARCHAR(20) NOT NULL,
  "street" VARCHAR(255) NOT NULL,
  "longitude" DECIMAL(10, 8) NOT NULL,
  "latitude" DECIMAL(11, 8) NOT NULL
);

-- Table: users

CREATE TYPE users_role_enum AS ENUM ('ADMIN', 'MANAGER', 'VIEWER');

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "first_name" VARCHAR(255) NOT NULL,
  "last_name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "role" users_role_enum NOT NULL,
  "deleted_at" TIMESTAMP
);

-- Table: cinemas

CREATE TABLE "cinemas" (
  "id" SERIAL PRIMARY KEY,
  "address_id" INTEGER NOT NULL REFERENCES "addresses"("id"),
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT DEFAULT '',
  "deleted_at" TIMESTAMP
  --"images" TEXT[] DEFAULT '{}'
);

-- Table: users_cinemas

CREATE TABLE "users_cinemas" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  "cinema_id" INTEGER NOT NULL REFERENCES "cinemas"("id"),
  UNIQUE ("user_id", "cinema_id")
);

-- Table: rooms

CREATE TYPE rooms_type_enum AS ENUM ('SMALL', 'MEDIUM', 'LARGE');

CREATE TABLE "rooms" (
  "id" SERIAL PRIMARY KEY,
  "cinema_id" INTEGER NOT NULL REFERENCES "cinemas"("id"),
  "number" VARCHAR(255) NOT NULL,
  "type" rooms_type_enum NOT NULL,
  UNIQUE ("cinema_id", "number")
);

-- Table: movies

CREATE TABLE "movies" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "description" TEXT DEFAULT '',
  "poster" TEXT NOT NULL,
  "backdrop" TEXT NOT NULL,
  "language" VARCHAR(255) NOT NULL,
  "released_at" VARCHAR(255) NOT NULL
);

-- Table: events

CREATE TABLE "events" (
  "id" SERIAL PRIMARY KEY,
  "room_id" INTEGER NOT NULL REFERENCES "rooms"("id"),
  "movie_id" INTEGER NOT NULL REFERENCES "movies"("id"),
  "price" INTEGER NOT NULL DEFAULT 0,
  "starts_at" TIMESTAMP NOT NULL,
  "ends_at" TIMESTAMP NOT NULL,
  "deleted_at" TIMESTAMP
);

-- Table: bookings

CREATE TABLE "bookings" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  "event_id" INTEGER NOT NULL REFERENCES "events"("id"),
  "place" VARCHAR(255) NOT NULL,
  "canceled_at" TIMESTAMP,
  UNIQUE ("user_id", "event_id", "place")
);
