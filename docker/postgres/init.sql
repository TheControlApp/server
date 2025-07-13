-- PostgreSQL initialization script for ControlMe
-- This script is run automatically when the container starts

-- Create the controlme database (if not already created by POSTGRES_DB)
-- CREATE DATABASE IF NOT EXISTS controlme;

-- Connect to the controlme database
\c controlme;

-- Enable UUID extension for UUID primary keys
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create a development user (optional)
-- CREATE USER controlme_dev WITH PASSWORD 'dev_password';
-- GRANT ALL PRIVILEGES ON DATABASE controlme TO controlme_dev;

-- Set timezone to UTC
SET timezone = 'UTC';

-- Create indexes that will be needed for performance
-- (These will be created by GORM migrations, but we can pre-create them)

-- Log the completion
\echo 'ControlMe PostgreSQL initialization complete';
