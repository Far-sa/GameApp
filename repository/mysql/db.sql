
CREATE TABLE users (
    id int PRIMARY KEY AUTO_INCREMENT,
    phone_number varchar(255) NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS users;
ALTER TABLE users add COLUMN password varchar(255) NOT NULL;
ALTER TABLE users drop COLUMN password;