

CREATE TABLE users (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    phone_number varchar(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE users;
ALTER TABLE users add COLUMN password varchar(255) NOT NULL;
ALTER TABLE users drop COLUMN password;