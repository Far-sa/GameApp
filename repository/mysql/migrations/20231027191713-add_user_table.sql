
-- +migrate Up
CREATE TABLE `users` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `phone_number` VARCHAR(255) NOT NULL UNIQUE,
    `name` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS `users`;