GRANT ALL PRIVILEGES ON *.* TO 'root' @'%';

USE `ellas-db`;

CREATE TABLE IF NOT EXISTS user (
    id varchar(40) NOT NULL PRIMARY KEY,
    `name` varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    provider_origin varchar(255) NOT NULL,
    external_id varchar(255) NOT NULL,
    profile_image_url varchar(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_external_id ON user (external_id);

CREATE TABLE IF NOT EXISTS user_role (
    id varchar(40) NOT NULL PRIMARY KEY,
    user_id varchar(40) NOT NULL,
    role ENUM('customer', 'service_provider', 'admin') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE INDEX idx_user_role_user_id ON user_role (user_id);