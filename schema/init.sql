-- Bastion Host Database Schema
-- TiDB/MySQL compatible

CREATE DATABASE IF NOT EXISTS bastion;
USE bastion;

-- User accounts
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    role VARCHAR(16) NOT NULL DEFAULT 'operator',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Managed server assets
CREATE TABLE IF NOT EXISTS assets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    host VARCHAR(256) NOT NULL,
    port INT NOT NULL DEFAULT 22,
    username VARCHAR(64) NOT NULL,
    auth_type VARCHAR(16) NOT NULL DEFAULT 'password',
    credential TEXT NOT NULL,
    status VARCHAR(16) DEFAULT 'offline',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- SSH sessions
CREATE TABLE IF NOT EXISTS sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    asset_id BIGINT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    client_ip VARCHAR(64) DEFAULT '',
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME NULL
);

-- Command audit logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    session_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    asset_id BIGINT NOT NULL,
    command TEXT,
    executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_audit_session (session_id),
    INDEX idx_audit_user (user_id),
    INDEX idx_audit_asset (asset_id),
    INDEX idx_audit_time (executed_at)
);

-- Seed admin user: admin / admin123 (bcrypt hash)
INSERT INTO users (username, password_hash, role) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin')
ON DUPLICATE KEY UPDATE username=username;
