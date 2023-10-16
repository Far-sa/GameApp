CREATE TABLE users (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    phone_number varchar(255) NOT NULL UNIQUE,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);