-- USERS TABLE
CREATE TABLE IF NOT EXISTS users (
    uid SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    work_location TEXT NOT NULL,
    balance DECIMAL(5,2) DEFAULT 5.0 CHECK(balance >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- AUTH SESSIONS TABLE
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(uid),
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    UNIQUE(user_id)
);

CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    assigned_to INT NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    description TEXT NOT NULL,
    task stat DEFAULT 'todo',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    priority INT NOT NULL CHECK(priority>=0 AND priority<=10)
);

CREATE TABLE IF NOT EXISTS notifications(
    id SERIAL PRIMARY KEY,
    assigned_to INT NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    created_at TIMESTAMPTZ WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE 
);

CREATE TYPE IF NOT EXISTS stat as ENUM('todo','inProgress','completed'); 