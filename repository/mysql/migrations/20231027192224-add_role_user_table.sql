
-- +migrate Up
-- Do not chage the order of enum values
-- TODO find a better solution to this problem
ALTER TABLE `users` ADD COLUMN `role` ENUM('user', 'admin') NOT NULL;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `role`; 