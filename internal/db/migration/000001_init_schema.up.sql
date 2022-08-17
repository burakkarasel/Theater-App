CREATE TABLE "movies" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "director_id" bigint NOT NULL,
  "rating" numeric NOT NULL,
  "poster" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "directors" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "oscars" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tickets" (
  "id" bigserial PRIMARY KEY,
  "movie_id" bigint NOT NULL,
  "ticket_owner" varchar NOT NULL,
  "child" smallint NOT NULL,
  "adult" smallint NOT NULL,
  "total" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "access_level" smallint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "movies" ("id");

CREATE INDEX ON "directors" ("id");

CREATE INDEX ON "tickets" ("id");

CREATE INDEX ON "users" ("username");

ALTER TABLE "movies" ADD FOREIGN KEY ("director_id") REFERENCES "directors" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("ticket_owner") REFERENCES "users" ("username");