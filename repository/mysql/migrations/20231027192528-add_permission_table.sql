
-- +migrate Up
CREATE TABLE `permissions` (
    `id` INT PRIMARY KEY AUTOINCREMENT,
    `title` VARCHAR(255) NOT NULL UNIQUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

-- +migrate Down
DROP TABLE IF EXISTS `permissions`;

