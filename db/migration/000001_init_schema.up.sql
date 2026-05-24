CREATE TABLE
    "users" (
        "id" bigserial PRIMARY KEY,
        "username" varchar UNIQUE NOT NULL,
        "hashed_password" varchar NOT NULL,
        "full_name" varchar NOT NULL,
        "email" varchar UNIQUE NOT NULL,
        "role" varchar NOT NULL DEFAULT 'depositor',
        "is_email_verified" boolean NOT NULL DEFAULT false,
        "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01',
        "created_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE TABLE
    "accounts" (
        "id" bigserial PRIMARY KEY,
        "owner" varchar NOT NULL REFERENCES users (username) ON DELETE CASCADE,
        "balance" bigint NOT NULL DEFAULT 0,
        "currency" varchar NOT NULL DEFAULT 'USD',
        "created_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

CREATE TABLE
    "entries" (
        "id" bigserial PRIMARY KEY,
        "account_id" bigint NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
        "amount" bigint NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE INDEX ON "entries" ("account_id");

CREATE TABLE
    "transfers" (
        "id" bigserial PRIMARY KEY,
        "from_account_id" bigint NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
        "to_account_id" bigint NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
        "amount" bigint NOT NULL CHECK (amount > 0),
        "created_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

CREATE TABLE
    "sessions" (
        "id" uuid PRIMARY KEY,
        "username" varchar NOT NULL REFERENCES users (username) ON DELETE CASCADE,
        "refresh_token" varchar NOT NULL,
        "user_agent" varchar NOT NULL,
        "client_ip" varchar NOT NULL,
        "is_blocked" boolean NOT NULL DEFAULT false,
        "expires_at" timestamptz NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE TABLE
    "verify_emails" (
        "id" bigserial PRIMARY KEY,
        "username" varchar NOT NULL REFERENCES users (username) ON DELETE CASCADE,
        "email" varchar NOT NULL,
        "secret_code" varchar NOT NULL,
        "is_used" boolean NOT NULL DEFAULT false,
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        "expired_at" timestamptz NOT NULL DEFAULT (now () + interval '15 minutes')
    );