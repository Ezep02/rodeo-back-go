create database elrodeodb;
use elrodeodb;

CREATE TABLE users (
  id SERIAL NOT NULL PRIMARY KEY,  
  name VARCHAR(45) NOT NULL,
  surname VARCHAR(70) DEFAULT NULL, 
  password VARCHAR(70) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  is_admin BOOL DEFAULT FALSE,
  is_barber bool default false,
  phone_number VARCHAR(30), 
  last_name_change TIMESTAMP NULL DEFAULT NULL,
  username VARCHAR(45) NOT NULL UNIQUE,
  avatar TEXT DEFAULT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT NULL
);



-- v1

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


-- NEW FEATURES
CREATE TABLE slots (
    id SERIAL PRIMARY KEY,
    barber_id BIGINT UNSIGNED NOT NULL,
    start DATETIME NOT NULL,
    end DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_barber_id FOREIGN KEY (barber_id) REFERENCES users(id) ON DELETE CASCADE
);

-- GOOGLE CALENDAR START
CREATE TABLE google_calendar_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expiry TIMESTAMP NOT NULL,
    token_type VARCHAR(50) NOT NULL DEFAULT 'Bearer',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- GOOGLE CALENDAR END

CREATE TABLE barbers (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    calendar_id VARCHAR(255),          -- ID del calendario de Google Calendar
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- SERVICES START
CREATE TABLE services (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    barber_id BIGINT UNSIGNED NOT NULL,
    preview_url TEXT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (barber_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE medias (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    service_id BIGINT UNSIGNED NOT NULL,
    url TEXT NOT NULL,
    type ENUM('image', 'video') DEFAULT 'image',
    position INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

-- Categorías
CREATE TABLE categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7),
    preview_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Relación muchos a muchos
CREATE TABLE service_categories (
    service_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (service_id, category_id),
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);



-- Promociones
CREATE TABLE promotions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    service_id BIGINT UNSIGNED NOT NULL,
    discount DECIMAL(5,2) NOT NULL, -- puede ser porcentaje
    start_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP NULL,
    type ENUM('percentage', 'fixed') DEFAULT 'percentage', 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

-- SERVICES END


-- PAYMENT AND BOOKING START
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    slot_id BIGINT UNSIGNED NOT NULL,
    client_id BIGINT UNSIGNED NOT NULL,
    
    -- Estado de la reserva (no financiero)
    status ENUM('pendiente_pago', 'confirmado', 'cancelado', 'rechazado', 'completado') NOT NULL DEFAULT 'pendiente_pago',
    
    total_amount DECIMAL(10,2) DEFAULT 0,
    google_event_id VARCHAR(255),
    
    coupon_code VARCHAR(12) DEFAULT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    
    expires_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_booking_slot FOREIGN KEY (slot_id) REFERENCES slots(id) ON DELETE CASCADE,
    CONSTRAINT fk_booking_client FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_booking_client (client_id),
    INDEX idx_booking_slot (slot_id),
    INDEX idx_booking_status (status)
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    booking_id BIGINT UNSIGNED NOT NULL,
    
    amount DECIMAL(10,2) NOT NULL,  -- monto del pago
    type ENUM('total','parcial','seña','restante') NOT NULL DEFAULT 'total',
    method ENUM('mercadopago','efectivo','tarjeta','transferencia') NOT NULL,
    
    status ENUM('pendiente','aprobado','rechazado','reembolsado') NOT NULL DEFAULT 'pendiente',

    mercado_pago_id VARCHAR(255) DEFAULT NULL,   -- ID en Mercado Pago
    payment_url TEXT DEFAULT NULL,               -- URL de preferencia / checkout
    paid_at DATETIME DEFAULT NULL,               -- fecha de confirmación de pago
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_payment_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    
    INDEX idx_payment_booking (booking_id),
    INDEX idx_payment_status (status),
    INDEX idx_mercado_pago_id (mercado_pago_id)
);


CREATE TABLE booking_services (
    id SERIAL PRIMARY KEY,
    booking_id BIGINT UNSIGNED NOT NULL,
    service_id BIGINT UNSIGNED NOT NULL,
    
    price DECIMAL(10,2) NOT NULL,   -- precio al momento de la reserva
    quantity INT DEFAULT 1,         -- unidades o cantidad de servicios
    notes VARCHAR(255) DEFAULT NULL, -- opcional: observaciones específicas del servicio

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_booking_service_booking FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
    CONSTRAINT fk_booking_service_service FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE,
    
    INDEX idx_booking_service_booking (booking_id),
    INDEX idx_booking_service_service (service_id)
);



-- PAYMENT AND BOOKING END






CREATE INDEX idx_slots_date ON slots(date, time);
-- OLD






