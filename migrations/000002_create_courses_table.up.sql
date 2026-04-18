CREATE TABLE IF NOT EXISTS courses (
                                       id SERIAL PRIMARY KEY,
                                       title VARCHAR(255) NOT NULL,
                                       description TEXT,
                                       price NUMERIC(10,2) NOT NULL,
                                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);