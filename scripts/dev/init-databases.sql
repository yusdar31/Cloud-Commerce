-- Initialize databases for each microservice
-- This runs automatically when PostgreSQL container starts for the first time

CREATE DATABASE identity_db;
CREATE DATABASE tenant_db;
CREATE DATABASE catalog_db;
CREATE DATABASE inventory_db;
CREATE DATABASE order_db;
CREATE DATABASE payment_db;
CREATE DATABASE notification_db;
CREATE DATABASE review_db;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE identity_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE tenant_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE catalog_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE inventory_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE order_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE payment_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE notification_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE review_db TO postgres;
