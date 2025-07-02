USE `boilerplate-db`;

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
    preferred_communication_channel ENUM(
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
