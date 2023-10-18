-- +migrate Up
ALTER TABLE users add COLUMN password varchar(255) NOT NULL;


-- +migrate Down
ALTER TABLE users drop COLUMN password;