BEGIN;
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id serial PRIMARY KEY,
    name varchar(50) UNIQUE NOT NULL,
    password VARCHAR(128) NOT NULL,
    salt varchar(50) NOT NULL,
    reauth boolean NOT NULL,
    created_on timestamp DEFAULT now(),
    last_login timestamp DEFAULT now(),
    archived_on timestamp DEFAULT NULL,
    is_admin boolean DEFAULT FALSE NOT NULL
);
COMMIT;

