
-- +migrate Up
CREATE TABLE `access_controls` (
    `id` INT PRIMARY KEY AUTOINCREMENT,
    `actor_id` VARCHAR(255) NOT NULL,
    `actor_type` ENUM('role','user') NOT NULL,
    `permission_id` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`permission_id`) REFERENCES `permission('id')`


);

-- +migrate Down
DROP TABLE IF EXISTS `access_controls`;