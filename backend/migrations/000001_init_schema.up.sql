-- Users
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Amenities
CREATE TABLE IF NOT EXISTS amenities (
    amenity_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Room Types
CREATE TABLE IF NOT EXISTS roomtypes (
    room_type_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    size_sqm DECIMAL(10, 2),
    bed_type VARCHAR(50),
    capacity INT DEFAULT 2,
    picture_url TEXT
);

-- Room Type Amenities
CREATE TABLE IF NOT EXISTS roomtype_amenities (
    room_type_id INT REFERENCES roomtypes(room_type_id) ON DELETE CASCADE,
    amenity_id INT REFERENCES amenities(amenity_id) ON DELETE CASCADE,
    PRIMARY KEY (room_type_id, amenity_id)
);

-- Rooms
CREATE TABLE IF NOT EXISTS rooms (
    room_id SERIAL PRIMARY KEY,
    room_type_id INT REFERENCES roomtypes(room_type_id) ON DELETE CASCADE,
    room_number VARCHAR(20) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'available' -- available, occupied, maintenance
);

-- Room Blocks (Maintenance/Manual Block)
CREATE TABLE IF NOT EXISTS room_blocks (
    block_id SERIAL PRIMARY KEY,
    room_id INT REFERENCES rooms(room_id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT
);

-- Rate Plans
CREATE TABLE IF NOT EXISTS rate_plans (
    rate_plan_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_special_package BOOLEAN DEFAULT FALSE,
    allow_free_cancel BOOLEAN DEFAULT FALSE,
    allow_pay_later BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Room Type Rate Prices
CREATE TABLE IF NOT EXISTS room_type_rate_prices (
    room_type_id INT REFERENCES roomtypes(room_type_id) ON DELETE CASCADE,
    rate_plan_id INT REFERENCES rate_plans(rate_plan_id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (room_type_id, rate_plan_id)
);

-- Addon Categories
CREATE TABLE IF NOT EXISTS addon_categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Addons
CREATE TABLE IF NOT EXISTS addons (
    addon_id SERIAL PRIMARY KEY,
    category_id INT REFERENCES addon_categories(category_id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) DEFAULT 0,
    unit_name VARCHAR(50) -- e.g. "per person", "per night"
);

-- Bookings
CREATE TABLE IF NOT EXISTS bookings (
    booking_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE SET NULL,
    rate_plan_id INT REFERENCES rate_plans(rate_plan_id) ON DELETE SET NULL,
    room_id INT REFERENCES rooms(room_id) ON DELETE SET NULL,
    
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    num_adults INT DEFAULT 1,
    
    room_subtotal DECIMAL(10, 2) DEFAULT 0,
    addon_subtotal DECIMAL(10, 2) DEFAULT 0,
    taxes_amount DECIMAL(10, 2) DEFAULT 0,
    total_price DECIMAL(10, 2) DEFAULT 0,
    
    status VARCHAR(20) DEFAULT 'pending', -- pending, confirmed, cancelled, completed
    expired_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Booking Addons
CREATE TABLE IF NOT EXISTS booking_addons (
    booking_addon_id SERIAL PRIMARY KEY,
    booking_id INT REFERENCES bookings(booking_id) ON DELETE CASCADE,
    addon_id INT REFERENCES addons(addon_id) ON DELETE SET NULL,
    quantity INT DEFAULT 1,
    price_at_time_of_booking DECIMAL(10, 2) NOT NULL
);
