-- +goose Up
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `users`;
 
CREATE TABLE `users` (
    `id`         varchar(36) NOT NULL,
    `name`       varchar(50) NOT NULL UNIQUE,
    `password`   binary(64) NOT NULL,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;
 
CREATE TABLE `tasks` (
    `id` varchar(36) NOT NULL,
	`user_id` varchar(36) NOT NULL,
    `title` varchar(50) NOT NULL,
    `is_done` boolean NOT NULL DEFAULT b'0',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
) DEFAULT CHARSET=utf8mb4;