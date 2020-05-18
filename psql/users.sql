DROP TABLE users;

CREATE TABLE users
(
    id serial PRIMARY KEY,
    name VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (128) NOT NULL,
    salt VARCHAR (50) NOT NULL,
    reauth BOOLEAN NOT NULL,
    created_on TIMESTAMP DEFAULT now(),
    last_login TIMESTAMP DEFAULT now(),
    archived_on TIMESTAMP DEFAULT NULL,
    is_admin boolean DEFAULT FALSE NOT NULL
);
