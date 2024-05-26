-- Create enum type "gender"
CREATE TYPE "public"."gender" AS ENUM ('MALE', 'FEMALE');
-- Create "admins" table
CREATE TABLE "public"."admins" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "unique_id" bigint NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_admins_deleted_at" to table: "admins"
CREATE INDEX "idx_admins_deleted_at" ON "public"."admins" ("deleted_at");
-- Create "subscriber_lists" table
CREATE TABLE "public"."subscriber_lists" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "unique_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_subscriber_lists_deleted_at" to table: "subscriber_lists"
CREATE INDEX "idx_subscriber_lists_deleted_at" ON "public"."subscriber_lists" ("deleted_at");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "unique_id" bigint NULL,
  "name" text NULL,
  "email" text NULL,
  "role" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create "subscribers" table
CREATE TABLE "public"."subscribers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "unique_id" bigint NULL,
  "list_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_subscriber_lists_subscribers" FOREIGN KEY ("list_id") REFERENCES "public"."subscriber_lists" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_subscribers_deleted_at" to table: "subscribers"
CREATE INDEX "idx_subscribers_deleted_at" ON "public"."subscribers" ("deleted_at");
