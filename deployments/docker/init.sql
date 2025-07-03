-- Create database if not exists
CREATE DATABASE golang_service;

-- Connect to the database
\c golang_service;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create initial tables (these will be managed by GORM migrations in the app)
-- This file is mainly for any custom setup or initial data