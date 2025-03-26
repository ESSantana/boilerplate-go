GRANT ALL PRIVILEGES ON *.* TO 'root' @'%';

USE `boilerplate-go`;

CREATE TABLE IF NOT EXISTS user (
    id varchar(40) NOT NULL PRIMARY KEY,
    `name` varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    provider_origin varchar(255) NOT NULL,
    external_id varchar(255) NOT NULL,
    profile_image_url varchar(255),
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    deleted_at datetime DEFAULT NULL
);

CREATE INDEX idx_external_id ON user (external_id);