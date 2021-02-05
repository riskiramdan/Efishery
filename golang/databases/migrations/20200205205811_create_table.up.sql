-- Table Definition ----------------------------------------------
CREATE TABLE "user" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "roleId" int NOT NULL,
  "name" varchar(80) NOT NULL,
  "phone" varchar(80) NOT NULL,
  "password" varchar NOT NULL,
  "token" varchar,
  "tokenExpiredAt" timestamp,
  "createdAt" timestamp NOT NULL DEFAULT (now()),
  "createdBy" varchar(20) DEFAULT 'admin',
  "updatedAt" timestamp NOT NULL DEFAULT (now()),
  "updatedBy" varchar(20) DEFAULT 'admin',
  "deletedAt" timestamp,
  "deletedBy" varchar(20)
);

CREATE TABLE "role" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" varchar(80) NOT NULL,
  "createdAt" timestamp NOT NULL DEFAULT (now()),
  "createdBy" varchar(20) DEFAULT 'admin',
  "updatedAt" timestamp NOT NULL DEFAULT (now()),
  "updatedBy" varchar(20) DEFAULT 'admin',
  "deletedAt" timestamp,
  "deletedBy" varchar(20)
);

ALTER TABLE "user" ADD FOREIGN KEY ("roleId") REFERENCES "role" ("id");
