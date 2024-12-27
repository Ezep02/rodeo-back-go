use goMeli;

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
  created_by_id int unsigned not null,
  service_duration INT DEFAULT NULL,
  preview_url text default null,
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
    price DECIMAL(12, 0) NOT NULL DEFAULT 0,
    service_duration INT DEFAULT 0,
    user_id INT UNSIGNED NOT NULL,
    barber_id INT UNSIGNED NOT NULL,
    service_id INT UNSIGNED NOT NULL,
    payment_id VARCHAR(150) NOT NULL,
    payer_name VARCHAR(50) NOT NULL,
    payer_surname VARCHAR(100) NOT NULL,
    payer_phone VARCHAR(30) DEFAULT NULL,
    email VARCHAR(255) NOT NULL,
    mp_order_id BIGINT DEFAULT NULL,
    mp_payment_id BIGINT DEFAULT NULL,
    date_approved TIMESTAMP NULL DEFAULT NULL,
    transaction_amount DECIMAL(12, 2) DEFAULT 0.00,
    fee_amount DECIMAL(12, 2) DEFAULT 0.00,
    mp_status VARCHAR(20) DEFAULT NULL,
    mp_status_detail VARCHAR(50) DEFAULT NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    date DATE NOT NULL, -- Nuevo campo para almacenar la fecha de la orden
    created_by_id INT UNSIGNED NOT NULL,
	shift_id INT UNSIGNED NOT NULL,
    weak_day VARCHAR(20) NOT NULL, 
    schedule VARCHAR(50) NOT NULL, 
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (service_id) REFERENCES services(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id) -- Relación con usuarios para "created_by_id"
);

alter table orders add 	shift_id INT UNSIGNED NOT NULL;

CREATE TABLE schedules (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    barber_id INT UNSIGNED NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NULL, 
    schedule_day ENUM('Lunes', 'Martes', 'Miercoles', 'Jueves', 'Viernes', 'Sabado', 'Domingo') NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (barber_id) REFERENCES users(id),
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE shifts (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    schedule_id INT UNSIGNED NOT NULL,
    created_by_name varchar(50) not null,
    day ENUM('Lunes', 'Martes', 'Miercoles', 'Jueves', 'Viernes', 'Sabado', 'Domingo') NOT NULL,
    start_time varchar(40) NOT NULL,
    available BOOL DEFAULT true,
    PRIMARY KEY (id),
    FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE,
    Shift_status varchar(50),
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

select * from schedules;

