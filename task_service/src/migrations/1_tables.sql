CREATE TYPE stat as ENUM('todo','inProgress','completed'); 


CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    assigned_by INT NOT NULL,
    assigned_to INT NOT NULL,
    description TEXT NOT NULL,
    task_status stat DEFAULT 'todo',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deadline TIMESTAMPTZ NOT NULL,
    priority INT NOT NULL CHECK(priority>=0 AND priority<=10),
    CHECK (deadline> CURRENT_TIMESTAMP)
);

