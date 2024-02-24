CREATE TYPE user_role as ENUM ('USER', 'SUPERUSER', 'ADMIN');

CREATE TYPE score_level AS ENUM ('EASY', 'MEDIUM', 'HARD');

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  username varchar(20) UNIQUE NOT NULL,
  email text UNIQUE NOT NULL,
  passhash varchar NOT NULL,
  role user_role NOT NULL,
  created_at varchar(30) NOT NULL DEFAULT TO_CHAR (
    current_timestamp,
    'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'
  ),
  updated_at varchar(30) NOT NULL DEFAULT TO_CHAR (
    current_timestamp,
    'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'
  )
);

CREATE TABLE IF NOT EXISTS scores (
  id bigserial PRIMARY KEY,
  level score_level NOT NULL,
  moves integer NOT NULL,
  time integer NOT NULL,
  created_at varchar(30) NOT NULL DEFAULT TO_CHAR (
    current_timestamp,
    'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'
  ),
  updated_at varchar(30) NOT NULL DEFAULT TO_CHAR (
    current_timestamp,
    'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'
  ),
  user_id bigserial REFERENCES users (id)
);
