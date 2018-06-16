package users

// Schema is the database schema for users
// it runs everytime the application starts
var Schema = []string{`
CREATE TABLE IF NOT EXISTS ` + Table + ` (
  id           SERIAL UNIQUE PRIMARY KEY,

  name         VARCHAR(25),
  username     VARCHAR(25) NOT NULL, -- 25 is more than enough -> 1234567890123456789012345 -> JetFuelCantMeltSteelBeams

  password     TEXT NOT NULL,
  email        VARCHAR(255) NOT NULL,

  deleted      BOOLEAN NOT NULL DEFAULT FALSE,
  activated    BOOLEAN NOT NULL DEFAULT FALSE,

  power        INTEGER NOT NULL DEFAULT 0,

  created      TIMESTAMP NOT NULL,
  seen         TIMESTAMP
);

`, `
CREATE TABLE IF NOT EXISTS ` + TableActivation + ` (
  id      SERIAL UNIQUE PRIMARY KEY,
  code    VARCHAR(255) NOT NULL,
  user_id INTEGER NOT NULL
);
`, `
CREATE TABLE IF NOT EXISTS ` + TableBan + ` (
  id        SERIAL UNIQUE PRIMARY KEY,
  user_id   INTEGER NOT NULL,
  state     BOOLEAN NOT NULL DEFAULT FALSE,
  temporary BOOLEAN NOT NULL DEFAULT FALSE,
  starts    TIMESTAMP,
  until     TIMESTAMP NOT NULL
);
`, `
CREATE TABLE IF NOT EXISTS ` + TableEvents + ` (
  id        SERIAL UNIQUE PRIMARY KEY,
  user_id   INTEGER NOT NULL,
  event     VARCHAR(255) NOT NULL,
  data      TEXT NULL,
  ip        INET NOT NULL,
  at        TIMESTAMP NOT NULL
);
`}

// SchemaTest is the database schema for testing the users table
// it runs before tests starts
var SchemaTest = []string{
	`TRUNCATE ` + Table + `, ` + TableActivation + `, ` + TableBan + `, ` + TableEvents + ` CASCADE;`,
}
