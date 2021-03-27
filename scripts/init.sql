CREATE TABLE users (
    user_id INT,
    PRIMARY KEY(user_id),
    username VARCHAR(39),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE gists (
    gist_unique_id INT GENERATED ALWAYS AS IDENTITY,
    gist_id VARCHAR(100),
    PRIMARY KEY(gist_unique_id),
    raw_url_link TEXT,
    username VARCHAR(39),
    gist_file_title TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE routine (
    routine_id INT GENERATED ALWAYS AS IDENTITY,
    PRIMARY KEY(routine_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE routine_gist_user(
    routine_id INT,
    CONSTRAINT fk_routine_id
                         FOREIGN KEY (routine_id)
                         REFERENCES routine(routine_id),
    gist_id VARCHAR(100),
    user_id INT,
    CONSTRAINT fk_user_id
                FOREIGN KEY(user_id)
                REFERENCES users(user_id)

);

CREATE TABLE session(
    session_id INT GENERATED ALWAYS AS IDENTITY,
    user_id INT,
    CONSTRAINT fk_user_id
               FOREIGN KEY(user_id)
               REFERENCES users(user_id),
    PRIMARY KEY(session_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
)
