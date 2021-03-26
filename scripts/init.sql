CREATE TABLE users (
    user_id INT,
    PRIMARY KEY(user_id),
    username VARCHAR(39),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE gists (
    gist_id VARCHAR(100),
    PRIMARY KEY(gist_id),
    raw_url_link TEXT,
    user_id INT,
    CONSTRAINT fk_user_id
                   FOREIGN KEY(user_id)
                   REFERENCES users(user_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE session (
    session_id INT GENERATED ALWAYS AS IDENTITY,
    PRIMARY KEY(session_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE session_gist(
    session_id INT,
    CONSTRAINT fk_session_id
                         FOREIGN KEY (session_id)
                         REFERENCES session(session_id),
    gist_id VARCHAR(100),
    CONSTRAINT fk_gist_id
                         FOREIGN KEY (gist_id)
                         REFERENCES gists(gist_id)

);
