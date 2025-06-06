create database elRodeo;
use elRodeo;

CREATE TABLE users (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,  
  name VARCHAR(45) NOT NULL,
  surname VARCHAR(70) DEFAULT NULL, 
  password VARCHAR(70) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  activo BOOL DEFAULT TRUE,
  is_admin BOOL DEFAULT FALSE,
  is_barber bool default false,
  phone_number VARCHAR(30),  
  PRIMARY KEY (id)
);


CREATE TABLE services (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT, 
  title VARCHAR(150) NOT NULL,
  description TEXT DEFAULT NULL,
  price DECIMAL(12, 0) NOT NULL DEFAULT 0,
  created_by_id INT UNSIGNED NOT NULL,
  service_duration INT DEFAULT NULL,
  preview_url TEXT DEFAULT NULL,
  category VARCHAR(40) NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (created_by_id) REFERENCES users(id)
);


CREATE TABLE orders (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    title VARCHAR(150) NOT NULL,
    description TEXT DEFAULT NULL,
    price FLOAT NOT NULL DEFAULT 0,
    service_duration INT DEFAULT 0,
    user_id INT UNSIGNED NOT NULL,
    barber_id INT UNSIGNED NOT NULL,
    service_id INT UNSIGNED NOT NULL,
    payer_name VARCHAR(50) NOT NULL,
    payer_surname VARCHAR(100) NOT NULL,
    payer_phone VARCHAR(30) DEFAULT NULL,
    email VARCHAR(255) NOT NULL,
    date_approved TIMESTAMP NULL DEFAULT NULL,
    mp_status VARCHAR(20) DEFAULT NULL,
    created_by_id INT UNSIGNED NOT NULL,
	shift_id INT UNSIGNED NOT NULL,
	schedule_day_date datetime not null,
    schedule_start_time VARCHAR(50) NOT NULL,
    transaction_type VARCHAR(60) NOT NULL, 
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (service_id) REFERENCES services(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id) 
);

CREATE TABLE reviews (
    order_id INT UNSIGNED NOT NULL,
    schedule_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    rating TINYINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review_status BOOLEAN DEFAULT TRUE,
    comment TEXT,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (order_id, schedule_id, user_id)
);

Create TABLE schedules (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	created_by_name varchar(50) not null,
    barber_id INT UNSIGNED NOT NULL,
	available BOOL DEFAULT true,
    schedule_day_date datetime not null,
    start_time varchar(40) not null,
	created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (barber_id) REFERENCES users(id)
);



CREATE TABLE coupons (
    code VARCHAR(12) PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    discount_percent DECIMAL(5,2),
    available BOOLEAN DEFAULT true,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    transaction_type VARCHAR(60) NOT NULL, 
    available_to_date DATETIME NOT NULL,
    used_at DATETIME default NULL,
    coupon_type ENUM('reembolso', 'promo'),
    FOREIGN KEY (user_id) REFERENCES users(id)
);



