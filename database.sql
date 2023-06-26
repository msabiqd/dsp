CREATE TABLE users (
	id serial PRIMARY KEY,
  full_name VARCHAR (60) NOT NULL,
  phone_number VARCHAR (13) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  successful_login INT DEFAULT 0
);

CREATE INDEX index_users_on_phone_number ON users USING HASH (phone_number);

INSERT INTO users (full_name, phone_number, password) VALUES ('User 1', '082111111111', 'Password1!');
INSERT INTO users (full_name, phone_number, password) VALUES ('User 2', '082111111112', 'Password2!');
