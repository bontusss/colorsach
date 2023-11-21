CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "bio" varchar,
  "email" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "collections" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "title" varchar NOT NULL,
  "user_id" integer,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "collections" ("name");

ALTER TABLE "collections" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
