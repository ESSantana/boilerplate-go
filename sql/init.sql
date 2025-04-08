GRANT ALL PRIVILEGES ON *.* TO 'admin' @'%' WITH GRANT OPTION;

FLUSH PRIVILEGES;

CREATE USER 'ellas-api' @'%' IDENTIFIED BY '5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8';

GRANT
SELECT
,
INSERT
,
UPDATE
    ON `ellas-db`.* TO 'ellas-api' @'%';

FLUSH PRIVILEGES;

USE `ellas-db`;

CREATE TABLE IF NOT EXISTS customer (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    nickname VARCHAR(50),
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(64),
    birth_date DATETIME NOT NULL,
    phone_number VARCHAR(11),
    cpf VARCHAR(11) NOT NULL,
    gender VARCHAR(30) NOT NULL,
    address VARCHAR(100) NOT NULL,
    address_number VARCHAR(10) NOT NULL,
    address_complement VARCHAR(50),
    address_neighborhood VARCHAR(50) NOT NULL,
    address_city VARCHAR(50) NOT NULL,
    address_state VARCHAR(50) NOT NULL,
    address_zip_code VARCHAR(10) NOT NULL,
    provider_origin VARCHAR(50) NOT NULL,
    external_id VARCHAR(100),
    profile_image_url TEXT,
    interests TEXT,
    how_heard_about_us TEXT,
    preferred_comunication_channel ENUM(
        'email',
        'sms',
        'whatsapp',
        'no-preference'
    ) DEFAULT 'no-preference',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_customer_external_id ON customer (external_id);
CREATE INDEX idx_customer_email ON customer (email);

CREATE TABLE IF NOT EXISTS staff (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    nickname VARCHAR(50),
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(64),
    staff_role ENUM('admin', 'manager', 'service_provider') NOT NULL,
    birth_date DATETIME NOT NULL,
    phone_number VARCHAR(11),
    document VARCHAR(14) NOT NULL,
    document_type ENUM('cpf', 'cnpj') NOT NULL,
    address VARCHAR(100) NOT NULL,
    address_number VARCHAR(10) NOT NULL,
    address_complement VARCHAR(50),
    address_neighborhood VARCHAR(50) NOT NULL,
    address_city VARCHAR(50) NOT NULL,
    address_state VARCHAR(50) NOT NULL,
    address_zip_code VARCHAR(10) NOT NULL,
    provider_origin VARCHAR(50) NOT NULL,
    external_id VARCHAR(100),
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
    payment_vendor VARCHAR(50) NOT NULL,
    payment_transaction_id VARCHAR(100) NOT NULL,
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

CREATE TABLE IF NOT EXISTS customer_event(
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(36) NOT NULL,
    event_data TEXT,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customer(id)
);

-- CREATE TABLE IF NOT EXISTS offered_product_schedule ();