CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
    roles (name, description, level)
VALUES
    (
        'user',
        'A user can create posts and comments.',
        1
    );

INSERT INTO
    roles (name, description, level)
VALUES
    (
        'moderator',
        'A moderator can update other users posts and comments.',
        2
    );

INSERT INTO
    roles (name, description, level)
VALUES
    (
        'admin',
        'An admin can update and deletes all posts and comments.',
        1
    );
