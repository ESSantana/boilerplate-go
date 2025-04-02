GRANT ALL PRIVILEGES ON *.* TO 'root' @'%' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'admin' @'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;

CREATE USER 'ellas-api'@'%' IDENTIFIED BY '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8';

GRANT SELECT, INSERT, UPDATE ON `ellas-db`.* TO 'ellas-api' @'%';

FLUSH PRIVILEGES;

USE `ellas-db`;

CREATE TABLE IF NOT EXISTS customer (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    birth_date DATETIME NOT NULL,
    phone_number VARCHAR(11) DEFAULT '',
    provider_origin VARCHAR(50) NOT NULL,
    external_id VARCHAR(100) NOT NULL,
    profile_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_customer_external_id ON customer (external_id);

CREATE TABLE IF NOT EXISTS staff (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    staff_role ENUM('admin', 'manager', 'service_provider') NOT NULL,
    birth_date DATETIME NOT NULL,
    phone_number VARCHAR(11) DEFAULT '',
    provider_origin VARCHAR(50) NOT NULL,
    external_id VARCHAR(100) NOT NULL,
    profile_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_staff_external_id ON staff (external_id);

CREATE TABLE IF NOT EXISTS product_category (
    id varchar(36) NOT NULL PRIMARY KEY,
    `name` varchar(40) NOT NULL,
    `description` TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS product (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    provider_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    `name` VARCHAR(50) NOT NULL,
    `description` TEXT,
    credit_cost INT NOT NULL,
    average_duration INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (provider_id) REFERENCES staff(id),
    FOREIGN KEY (category_id) REFERENCES product_category(id)
);

CREATE INDEX idx_product_name ON product (`name`);

CREATE TABLE IF NOT EXISTS credit_purchase_history (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    credit_amount INT NOT NULL,
    payment_value VARCHAR(10) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_status ENUM(
        'pending',
        'paid',
        'canceled',
        'failed',
        'refunded'
    ) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customer(id)
);

CREATE TABLE IF NOT EXISTS credit_balance (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    credit_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customer(id)
);

-- CREATE TABLE IF NOT EXISTS offered_product_schedule ();