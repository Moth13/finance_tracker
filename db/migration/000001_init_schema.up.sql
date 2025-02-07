CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "currency" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "create_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "init_balance" numeric(19,4) NOT NULL DEFAULT 0,
  "balance" numeric(19,4) NOT NULL DEFAULT 0,
  "final_balance" numeric(19,4) NOT NULL DEFAULT 0
);

CREATE TABLE "lines" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "account_id" bigint NOT NULL,
  "month_id" bigint NOT NULL,
  "year_id" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "amount" numeric(19,4) NOT NULL,
  "checked" bool DEFAULT (false) NOT NULL,
  "description" varchar NOT NULL,
  "due_date" date NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "reclines" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "account_id" bigint NOT NULL,
  "amount" numeric(19,4) NOT NULL,
  "category_id" bigint NOT NULL,
  "description" varchar NOT NULL,
  "recurrency" varchar NOT NULL,
  "due_date" date NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "months" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "year_id" bigint NOT NULL,
  "balance" numeric(19,4) NOT NULL DEFAULT 0,
  "final_balance" numeric(19,4) NOT NULL DEFAULT 0,
  "start_date" date NOT NULL DEFAULT '0001-01-01',
  "end_date" date NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "years" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "balance" numeric(19,4) NOT NULL DEFAULT 0,
  "final_balance" numeric(19,4) NOT NULL DEFAULT 0,
  "start_date" date NOT NULL DEFAULT '0001-01-01',
  "end_date" date NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "owner" varchar NOT NULL
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "lines" ("account_id");

CREATE INDEX ON "lines" ("month_id");

CREATE INDEX ON "lines" ("owner");

CREATE INDEX ON "lines" ("category_id");

CREATE INDEX ON "lines" ("account_id", "month_id", "owner", "category_id");

CREATE INDEX ON "reclines" ("account_id");

CREATE INDEX ON "reclines" ("category_id");

CREATE INDEX ON "reclines" ("account_id", "owner", "category_id");

CREATE INDEX ON "months" ("year_id");

CREATE INDEX ON "months" ("owner");

CREATE INDEX ON "months" ("year_id", "owner");

CREATE INDEX ON "years" ("owner");

CREATE INDEX ON "categories" ("owner");

COMMENT ON COLUMN "lines"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "reclines"."amount" IS 'can be negative or position';

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "lines" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "lines" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "lines" ADD FOREIGN KEY ("month_id") REFERENCES "months" ("id");

ALTER TABLE "lines" ADD FOREIGN KEY ("year_id") REFERENCES "years" ("id");

ALTER TABLE "lines" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "reclines" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "reclines" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "reclines" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "months" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "months" ADD FOREIGN KEY ("year_id") REFERENCES "years" ("id");

ALTER TABLE "years" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "categories" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");