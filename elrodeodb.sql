create database elRodeo;
use elRodeo;

CREATE TABLE users (
  id SERIAL NOT NULL PRIMARY KEY,  
  name VARCHAR(45) NOT NULL,
  surname VARCHAR(70) DEFAULT NULL, 
  password VARCHAR(70) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  activo BOOL DEFAULT TRUE,
  is_admin BOOL DEFAULT FALSE,
  is_barber bool default false,
  phone_number VARCHAR(30), 
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NULL
);


-- v1
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description VARCHAR(500),
  price DECIMAL(12, 2) NOT NULL,
  category VARCHAR(40) NOT NULL,
  preview_url TEXT,
  rating_sum int default 0,
  number_of_reviews int default 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE slots (
    id SERIAL PRIMARY KEY,
    date datetime not null,
    time VARCHAR(30) NOT NULL,
    is_booked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE appointments (
    id SERIAL PRIMARY KEY,
    client_name VARCHAR(100) NOT NULL,
    client_surname VARCHAR(100) NOT NULL,
    slot_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED not null,
    status VARCHAR(100) not null,
    payment_percentage INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_slot FOREIGN KEY (slot_id) REFERENCES slots(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE reviews (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    appointment_id BIGINT UNSIGNED NOT NULL UNIQUE,
    rating TINYINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_review_appointment FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE
);


CREATE TABLE coupons (
    id SERIAL PRIMARY KEY,
    code VARCHAR(12) UNIQUE,
    user_id BIGINT UNSIGNED NOT NULL,
    discount_percentage DECIMAL(10,2),
    is_available BOOL default true,
    used_at DATETIME default NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expire_at DATETIME DEFAULT NULL,
    CONSTRAINT fk_user_coupon FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);



-- Relacion many to many entre products y appointment
CREATE TABLE appointment_products (
    appointment_id BIGINT UNSIGNED REFERENCES appointments(id) ON DELETE CASCADE,
    product_id BIGINT UNSIGNED REFERENCES products(id) ON DELETE CASCADE,
    PRIMARY KEY (appointment_id, product_id)
);



CREATE INDEX idx_slots_date ON slots(date, time);
-- OLD






